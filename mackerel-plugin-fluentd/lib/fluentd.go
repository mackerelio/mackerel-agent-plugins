package mpfluentd

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/golib/logging"
)

var logger = logging.GetLogger("metrics.plugin.fluentd")

func metricName(names ...string) string {
	return strings.Join(names, ".")
}

// FluentdPlugin mackerel plugin for Fluentd
type FluentdPlugin struct {
	Host            string
	Port            string
	Prefix          string
	Tempfile        string
	pluginType      string
	pluginIDPattern *regexp.Regexp
	extendedMetrics []string
	Workers         uint

	plugins []FluentdPluginMetrics
}

// FluentdMetrics is alias for backward compatibility.
type FluentdMetrics = FluentdPlugin

// MetricKeyPrefix interface for PluginWithPrefix
func (f FluentdPlugin) MetricKeyPrefix() string {
	if f.Prefix == "" {
		f.Prefix = "fluentd"
	}
	return f.Prefix
}

// FluentdPluginMetrics metrics
type FluentdPluginMetrics struct {
	RetryCount            uint64 `json:"retry_count"`
	BufferQueueLength     uint64 `json:"buffer_queue_length"`
	BufferTotalQueuedSize uint64 `json:"buffer_total_queued_size"`
	OutputPlugin          bool   `json:"output_plugin"`
	Type                  string `json:"type"`
	PluginCategory        string `json:"plugin_category"`
	PluginID              string `json:"plugin_id"`
	normalizedPluginID    string

	// extended metrics fluentd >= 1.6
	// https://www.fluentd.org/blog/fluentd-v1.6.0-has-been-released
	EmitRecords                      uint64  `json:"emit_records"`
	EmitCount                        uint64  `json:"emit_count"`
	WriteCount                       uint64  `json:"write_count"`
	RollbackCount                    uint64  `json:"rollback_count"`
	SlowFlushCount                   uint64  `json:"slow_flush_count"`
	FlushTimeCount                   uint64  `json:"flush_time_count"`
	BufferStageLength                uint64  `json:"buffer_stage_length"`
	BufferStageByteSize              uint64  `json:"buffer_stage_byte_size"`
	BufferQueueByteSize              uint64  `json:"buffer_queue_byte_size"`
	BufferAvailableBufferSpaceRatios float64 `json:"buffer_available_buffer_space_ratios"`
}

func (fpm FluentdPluginMetrics) getExtended(name string) float64 {
	switch name {
	case "emit_records":
		return float64(fpm.EmitRecords)
	case "emit_count":
		return float64(fpm.EmitCount)
	case "write_count":
		return float64(fpm.WriteCount)
	case "rollback_count":
		return float64(fpm.RollbackCount)
	case "slow_flush_count":
		return float64(fpm.SlowFlushCount)
	case "flush_time_count":
		return float64(fpm.FlushTimeCount)
	case "buffer_stage_length":
		return float64(fpm.BufferStageLength)
	case "buffer_stage_byte_size":
		return float64(fpm.BufferStageByteSize)
	case "buffer_queue_byte_size":
		return float64(fpm.BufferQueueByteSize)
	case "buffer_available_buffer_space_ratios":
		return fpm.BufferAvailableBufferSpaceRatios
	}
	return 0
}

// FluentMonitorJSON monitor json
type FluentMonitorJSON struct {
	Plugins []FluentdPluginMetrics `json:"plugins"`
}

var normalizePluginIDRe = regexp.MustCompile(`[^-a-zA-Z0-9_]`)

func normalizePluginID(in string) string {
	return normalizePluginIDRe.ReplaceAllString(in, "_")
}

func (fpm FluentdPluginMetrics) getNormalizedPluginID() string {
	if fpm.normalizedPluginID == "" {
		fpm.normalizedPluginID = normalizePluginID(fpm.PluginID)
	}
	return fpm.normalizedPluginID
}

func (f *FluentdPlugin) parseStats(body []byte) (map[string]interface{}, error) {
	var j FluentMonitorJSON
	err := json.Unmarshal(body, &j)
	f.plugins = j.Plugins

	metrics := make(map[string]interface{})
	for _, p := range f.plugins {
		if f.nonTargetPlugin(p) {
			continue
		}
		pid := p.getNormalizedPluginID()
		metrics[metricName("retry_count", pid)] = float64(p.RetryCount)
		metrics[metricName("buffer_queue_length", pid)] = float64(p.BufferQueueLength)
		metrics[metricName("buffer_total_queued_size", pid)] = float64(p.BufferTotalQueuedSize)
		for _, name := range f.extendedMetrics {
			metrics[metricName(name, pid)] = p.getExtended(name)
		}
	}
	return metrics, err
}

func (f *FluentdPlugin) nonTargetPlugin(plugin FluentdPluginMetrics) bool {
	if plugin.PluginCategory != "output" {
		return true
	}
	if f.pluginType != "" && f.pluginType != plugin.Type {
		return true
	}
	if f.pluginIDPattern != nil && !f.pluginIDPattern.MatchString(plugin.PluginID) {
		return true
	}
	return false
}

func (f *FluentdPlugin) fetchFluentdMetrics(host string, port int) (map[string]interface{}, error) {
	target := fmt.Sprintf("http://%s:%d/api/plugins.json", host, port)
	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "mackerel-plugin-fluentd")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return f.parseStats(body)
}

// FetchMetrics interface for mackerelplugin
func (f FluentdPlugin) FetchMetrics() (map[string]interface{}, error) {
	port, _ := strconv.Atoi(f.Port)
	if f.Workers > 1 {
		metrics := make(map[string]interface{})
		for workerNumber := 0; workerNumber < int(f.Workers); workerNumber++ {
			m, e := f.fetchFluentdMetrics(f.Host, port+workerNumber)
			if e != nil {
				continue
			}

			workerName := fmt.Sprintf("worker%d", workerNumber)
			for k, v := range m {
				ks := strings.Split(k, ".")
				ks, last := ks[:len(ks)-1], ks[len(ks)-1]
				ks = append(ks, workerName)
				ks = append(ks, last)
				metrics[strings.Join(ks, ".")] = v
			}
		}
		if len(metrics) == 0 {
			err := fmt.Errorf("failed to connect to fluentd's monitor_agent")
			return metrics, err
		}
		return metrics, nil
	}
	return f.fetchFluentdMetrics(f.Host, port)
}

var defaultGraphs = map[string]mp.Graphs{
	"retry_count": {
		Label: "retry count",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
	"buffer_queue_length": {
		Label: "queue length",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
	"buffer_total_queued_size": {
		Label: "buffer total queued size",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
}

var extendedGraphs = map[string]mp.Graphs{
	"emit_records": {
		Label: "emitted records",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: true},
		},
	},
	"emit_count": {
		Label: "emit calls",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: true},
		},
	},
	"write_count": {
		Label: "write/try_write calls",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: true},
		},
	},
	"rollback_count": {
		Label: "rollbacks",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: true},
		},
	},
	"slow_flush_count": {
		Label: "slow flushes",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: true},
		},
	},
	"flush_time_count": {
		Label: "buffer flush time in msec",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: true},
		},
	},
	"buffer_stage_length": {
		Label: "length of staged buffer chunks",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
	"buffer_stage_byte_size": {
		Label: "bytesize of staged buffer chunks",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
	"buffer_queue_byte_size": {
		Label: "bytesize of queued buffer chunks",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
	"buffer_available_buffer_space_ratios": {
		Label: "available space for buffer",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
}

// GraphDefinition interface for mackerelplugin
func (f FluentdPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(f.Prefix)
	graphs := make(map[string]mp.Graphs, len(defaultGraphs))
	for key, g := range defaultGraphs {
		graphs[key] = mp.Graphs{
			Label:   (labelPrefix + " " + g.Label),
			Unit:    g.Unit,
			Metrics: g.Metrics,
		}
	}
	for _, name := range f.extendedMetrics {
		fullName := metricName(name)
		if g, ok := extendedGraphs[fullName]; ok {
			graphs[fullName] = mp.Graphs{
				Label:   (labelPrefix + " " + g.Label),
				Unit:    g.Unit,
				Metrics: g.Metrics,
			}
		}
	}
	return graphs
}

// Do the plugin
func Do() {
	host := flag.String("host", "localhost", "fluentd monitor_agent host")
	port := flag.String("port", "24220", "fluentd monitor_agent port")
	pluginType := flag.String("plugin-type", "", "Gets the metric that matches this plugin type")
	pluginIDPatternString := flag.String("plugin-id-pattern", "", "Gets the metric that matches this plugin id pattern")
	prefix := flag.String("metric-key-prefix", "fluentd", "Metric key prefix")
	tempFile := flag.String("tempfile", "", "Temp file name")
	extendedMetricNames := flag.String("extended_metrics", "", "extended metric names joind with ',' or 'all' (fluentd >= v1.6.0)")
	workers := flag.Uint("workers", 1, "specifying the number of Fluentd's multi-process workers")
	flag.Parse()

	var pluginIDPattern *regexp.Regexp
	var err error
	if *pluginIDPatternString != "" {
		pluginIDPattern, err = regexp.Compile(*pluginIDPatternString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to exec mackerel-plugin-fluentd: invalid plugin-id-pattern: %s\n", err)
			os.Exit(1)
		}
	}

	var extendedMetrics []string
	switch *extendedMetricNames {
	case "all":
		for key := range extendedGraphs {
			extendedMetrics = append(extendedMetrics, key)
		}
	case "":
	default:
		for _, name := range strings.Split(*extendedMetricNames, ",") {
			fullName := metricName(name)
			if _, exists := extendedGraphs[fullName]; !exists {
				fmt.Fprintf(os.Stderr, "extended_metrics %s is not supported. See also https://www.fluentd.org/blog/fluentd-v1.6.0-has-been-released\n", name)
				os.Exit(1)
			}
			extendedMetrics = append(extendedMetrics, name)
		}
	}
	f := FluentdPlugin{
		Host:            *host,
		Port:            *port,
		Prefix:          *prefix,
		Tempfile:        *tempFile,
		pluginType:      *pluginType,
		pluginIDPattern: pluginIDPattern,
		extendedMetrics: extendedMetrics,
		Workers:         *workers,
	}

	helper := mp.NewMackerelPlugin(f)

	helper.Tempfile = *tempFile
	if *tempFile == "" {
		tempFileSuffix := []string{*host, *port}
		if *pluginType != "" {
			tempFileSuffix = append(tempFileSuffix, *pluginType)
		}
		if *pluginIDPatternString != "" {
			tempFileSuffix = append(tempFileSuffix, fmt.Sprintf("%x", md5.Sum([]byte(*pluginIDPatternString))))
		}
		helper.SetTempfileByBasename(fmt.Sprintf("mackerel-plugin-fluentd-%s", strings.Join(tempFileSuffix, "-")))
	}

	helper.Run()
}
