package mpsidekiq

import (
	"flag"
	"fmt"
	"strconv"

	r "github.com/go-redis/redis"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// SidekiqPlugin for fetching metrics
type SidekiqPlugin struct {
	Tempfile string
	client   *r.Client
}

var graphdef = map[string]mp.Graphs{
	"Sidekiq.ProcessedANDFailed": mp.Graphs{
		Label: "Sidekiq processed and failed count",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "processed", Label: "Processed", Type: "uint64", Diff: true},
			{Name: "failed", Label: "Failed", Type: "uint64", Diff: true},
		},
	},
	"Sidekiq.Stats": mp.Graphs{
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
}

// GraphDefinition Graph definition
func (sp SidekiqPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

func (sp SidekiqPlugin) get(key string) uint64 {
	val, err := sp.client.Get(key).Result()
	if err == r.Nil {
		return 0
	}

	r, _ := strconv.ParseUint(val, 10, 64)
	return r
}

func (sp SidekiqPlugin) zCard(key string) uint64 {
	val, err := sp.client.ZCard(key).Result()
	if err == r.Nil {
		return 0
	}

	return uint64(val)
}

func (sp SidekiqPlugin) sMembers(key string) []string {
	val, err := sp.client.SMembers(key).Result()
	if err == r.Nil {
		return make([]string, 0)
	}

	return val
}

func (sp SidekiqPlugin) hGet(key string, field string) uint64 {
	val, err := sp.client.HGet(key, field).Result()
	if err == r.Nil {
		return 0
	}

	r, _ := strconv.ParseUint(val, 10, 64)
	return r
}

func (sp SidekiqPlugin) lLen(key string) uint64 {
	val, err := sp.client.LLen(key).Result()
	if err == r.Nil {
		return 0
	}

	return uint64(val)
}

func (sp SidekiqPlugin) getProcessed() uint64 {
	return sp.get("stat:processed")
}

func (sp SidekiqPlugin) getFailed() uint64 {
	return sp.get("stat:failed")
}

func inject(slice []uint64, base uint64) uint64 {
	for _, e := range slice {
		base += uint64(e)
	}

	return base
}

func (sp SidekiqPlugin) getBusy() uint64 {
	processes := sp.sMembers("processes")
	busies := make([]uint64, 10)
	for _, e := range processes {
		busies = append(busies, sp.hGet(e, "busy"))
	}

	return inject(busies, 0)
}

func (sp SidekiqPlugin) getEnqueued() uint64 {
	queues := sp.sMembers("queues")
	queuesLlens := make([]uint64, 10)

	for _, e := range queues {
		queuesLlens = append(queuesLlens, sp.lLen("queue:"+e))
	}

	return inject(queuesLlens, 0)
}

func (sp SidekiqPlugin) getSchedule() uint64 {
	return sp.zCard("schedule")
}

func (sp SidekiqPlugin) getRetry() uint64 {
	return sp.zCard("retry")
}

func (sp SidekiqPlugin) getDead() uint64 {
	return sp.zCard("dead")
}

func (sp SidekiqPlugin) getProcessedFailed() map[string]interface{} {
	data := make(map[string]interface{}, 20)

	data["processed"] = sp.getProcessed()
	data["failed"] = sp.getFailed()

	return data
}

func (sp SidekiqPlugin) getStats(field []string) map[string]interface{} {
	stats := make(map[string]interface{}, 20)
	for _, e := range field {
		switch e {
		case "busy":
			stats[e] = sp.getBusy()
		case "enqueued":
			stats[e] = sp.getEnqueued()
		case "schedule":
			stats[e] = sp.getSchedule()
		case "retry":
			stats[e] = sp.getRetry()
		case "dead":
			stats[e] = sp.getDead()
		}
	}

	return stats
}

// FetchMetrics fetch the metrics
func (sp SidekiqPlugin) FetchMetrics() (map[string]interface{}, error) {
	field := []string{"busy", "enqueued", "schedule", "retry", "dead"}
	stats := sp.getStats(field)
	pf := sp.getProcessedFailed()

	// merge maps
	m := func(map1 map[string]interface{}, map2 map[string]interface{}) map[string]interface{} {
		for k, v := range map2 {
			map1[k] = v
		}

		return map1
	}(stats, pf)

	return m, nil
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "6379", "Port")
	optPassword := flag.String("password", "", "Password")
	optDB := flag.Int("db", 0, "DB")
	optTempfile := flag.String("tempfile", "/tmp/mackerel-plugin-sidekiq", "Temp file name")
	flag.Parse()

	client := r.NewClient(&r.Options{
		Addr:     fmt.Sprintf("%s:%s", *optHost, *optPort),
		Password: *optPassword, // no password set
		DB:       *optDB,       // use default DB
	})

	sp := SidekiqPlugin{client: client}
	helper := mp.NewMackerelPlugin(sp)
	helper.Tempfile = *optTempfile

	helper.Run()
}
