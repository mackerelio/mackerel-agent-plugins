package mpsidekiq

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	r "github.com/redis/go-redis/v9"
)

// SidekiqPlugin for fetching metrics
type SidekiqPlugin struct {
	Client    *r.Client
	Namespace string
	Prefix    string
}

var graphdef = map[string]mp.Graphs{
	"ProcessedANDFailed": {
		Label: "Sidekiq processed and failed count",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "processed", Label: "Processed", Type: "uint64", Diff: true},
			{Name: "failed", Label: "Failed", Type: "uint64", Diff: true},
		},
	},
	"Stats": {
		Label: "Sidekiq stats",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "busy", Label: "Busy", Type: "uint64"},
			{Name: "enqueued", Label: "Enqueued", Type: "uint64"},
			{Name: "schedule", Label: "Schedule", Type: "uint64"},
			{Name: "retry", Label: "Retry", Type: "uint64"},
			{Name: "dead", Label: "Dead", Type: "uint64"},
		},
	},
	"QueueLatency": {
		Label: "Sidekiq queue latency",
		Unit:  "float",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1"},
		},
	},
}

// GraphDefinition Graph definition
func (sp SidekiqPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

func (sp SidekiqPlugin) get(ctx context.Context, key string) uint64 {
	val, err := sp.Client.Get(ctx, key).Result()
	if err == r.Nil {
		return 0
	}

	r, _ := strconv.ParseUint(val, 10, 64)
	return r
}

func (sp SidekiqPlugin) zCard(ctx context.Context, key string) uint64 {
	val, err := sp.Client.ZCard(ctx, key).Result()
	if err == r.Nil {
		return 0
	}

	return uint64(val)
}

func (sp SidekiqPlugin) sMembers(ctx context.Context, key string) []string {
	val, err := sp.Client.SMembers(ctx, key).Result()
	if err == r.Nil {
		return make([]string, 0)
	}

	return val
}

func (sp SidekiqPlugin) hGet(ctx context.Context, key string, field string) uint64 {
	val, err := sp.Client.HGet(ctx, key, field).Result()
	if err == r.Nil {
		return 0
	}

	r, _ := strconv.ParseUint(val, 10, 64)
	return r
}

func (sp SidekiqPlugin) lLen(ctx context.Context, key string) uint64 {
	val, err := sp.Client.LLen(ctx, key).Result()
	if err == r.Nil {
		return 0
	}

	return uint64(val)
}

func addNamespace(namespace string, key string) string {
	if namespace == "" {
		return key
	}
	return namespace + ":" + key
}

func (sp SidekiqPlugin) getProcessed(ctx context.Context) uint64 {
	key := addNamespace(sp.Namespace, "stat:processed")
	return sp.get(ctx, key)
}

func (sp SidekiqPlugin) getFailed(ctx context.Context) uint64 {
	key := addNamespace(sp.Namespace, "stat:failed")
	return sp.get(ctx, key)
}

func inject(slice []uint64, base uint64) uint64 {
	for _, e := range slice {
		base += uint64(e)
	}

	return base
}

func (sp SidekiqPlugin) getBusy(ctx context.Context) uint64 {
	key := addNamespace(sp.Namespace, "processes")
	processes := sp.sMembers(ctx, key)
	busies := make([]uint64, 10)
	for _, e := range processes {
		e := addNamespace(sp.Namespace, e)
		busies = append(busies, sp.hGet(ctx, e, "busy"))
	}

	return inject(busies, 0)
}

func (sp SidekiqPlugin) getEnqueued(ctx context.Context) uint64 {
	key := addNamespace(sp.Namespace, "queues")
	queues := sp.sMembers(ctx, key)
	queuesLlens := make([]uint64, 10)

	prefix := addNamespace(sp.Namespace, "queue:")
	for _, e := range queues {
		queuesLlens = append(queuesLlens, sp.lLen(ctx, prefix+e))
	}

	return inject(queuesLlens, 0)
}

func (sp SidekiqPlugin) getSchedule(ctx context.Context) uint64 {
	key := addNamespace(sp.Namespace, "schedule")
	return sp.zCard(ctx, key)
}

func (sp SidekiqPlugin) getRetry(ctx context.Context) uint64 {
	key := addNamespace(sp.Namespace, "retry")
	return sp.zCard(ctx, key)
}

func (sp SidekiqPlugin) getDead(ctx context.Context) uint64 {
	key := addNamespace(sp.Namespace, "dead")
	return sp.zCard(ctx, key)
}

func (sp SidekiqPlugin) getProcessedFailed(ctx context.Context) map[string]interface{} {
	data := make(map[string]interface{}, 20)

	data["processed"] = sp.getProcessed(ctx)
	data["failed"] = sp.getFailed(ctx)

	return data
}

func (sp SidekiqPlugin) getStats(ctx context.Context, field []string) map[string]interface{} {
	stats := make(map[string]interface{}, 20)
	for _, e := range field {
		switch e {
		case "busy":
			stats[e] = sp.getBusy(ctx)
		case "enqueued":
			stats[e] = sp.getEnqueued(ctx)
		case "schedule":
			stats[e] = sp.getSchedule(ctx)
		case "retry":
			stats[e] = sp.getRetry(ctx)
		case "dead":
			stats[e] = sp.getDead(ctx)
		}
	}

	return stats
}

func metricName(names ...string) string {
	return strings.Join(names, ".")
}

func (sp SidekiqPlugin) getQueueLatency(ctx context.Context) map[string]interface{} {
	latency := make(map[string]interface{}, 10)

	key := addNamespace(sp.Namespace, "queues")
	queues := sp.sMembers(ctx, key)

	prefix := addNamespace(sp.Namespace, "queue:")
	for _, q := range queues {
		queuesLRange, err := sp.Client.LRange(ctx, prefix+q, -1, -1).Result()
		if err != nil {
			fmt.Fprintf(os.Stderr, "get last queue error")
		}

		if len(queuesLRange) == 0 {
			latency[metricName("QueueLatency", q)] = 0.0
			continue
		}
		var job map[string]interface{}
		var thence float64

		err = json.Unmarshal([]byte(queuesLRange[0]), &job)
		if err != nil {
			fmt.Fprintf(os.Stderr, "json parse error")
			continue
		}
		now := float64(time.Now().Unix())
		if enqueuedAt, ok := job["enqueued_at"]; ok {
			enqueuedAt := enqueuedAt.(float64)
			thence = enqueuedAt
		} else {
			thence = now
		}
		latency[metricName("QueueLatency", q)] = (now - thence)
	}

	return latency
}

// FetchMetrics fetch the metrics
func (sp SidekiqPlugin) FetchMetrics() (map[string]interface{}, error) {
	field := []string{"busy", "enqueued", "schedule", "retry", "dead", "latency"}
	ctx := context.Background()
	stats := sp.getStats(ctx, field)
	pf := sp.getProcessedFailed(ctx)
	latency := sp.getQueueLatency(ctx)

	// merge maps
	m := func(m ...map[string]interface{}) map[string]interface{} {
		r := make(map[string]interface{}, 20)
		for _, c := range m {
			for k, v := range c {
				r[k] = v
			}
		}

		return r
	}(stats, pf, latency)

	return m, nil
}

// MetricKeyPrefix interface for PluginWithPrefix
func (sp SidekiqPlugin) MetricKeyPrefix() string {
	if sp.Prefix == "" {
		sp.Prefix = "sidekiq"
	}
	return sp.Prefix
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "6379", "Port")
	optPassword := flag.String("password", os.Getenv("SIDEKIQ_PASSWORD"), "Password")
	optDB := flag.Int("db", 0, "DB")
	optNamespace := flag.String("redis-namespace", "", "Redis namespace")
	optPrefix := flag.String("metric-key-prefix", "sidekiq", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	client := r.NewClient(&r.Options{
		Addr:     fmt.Sprintf("%s:%s", *optHost, *optPort),
		Password: *optPassword,
		DB:       *optDB,
	})

	sp := SidekiqPlugin{
		Client:    client,
		Namespace: *optNamespace,
		Prefix:    *optPrefix,
	}
	helper := mp.NewMackerelPlugin(sp)
	helper.Tempfile = *optTempfile

	helper.Run()
}
