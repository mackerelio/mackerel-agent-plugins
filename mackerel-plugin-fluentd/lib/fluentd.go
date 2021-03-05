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
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/golib/logging"
)

var logger = logging.GetLogger("metrics.plugin.fluentd")

// FluentdMetrics plugin for fluentd
type FluentdMetrics struct {
	Target          string
	Tempfile        string
	pluginType      string
	pluginIDPattern *regexp.Regexp
	extendedMetrics []string

	plugins []FluentdPluginMetrics
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
		return float64(f.EmitRecords)
	case "emit_count":
		return float64(f.EmitCount)
	case "write_count":
		return float64(f.WriteCount)
	case "rollback_count":
		return float64(f.RollbackCount)
	case "slow_flush_count":
		return float64(f.SlowFlushCount)
	case "flush_time_count":
		return float64(f.FlushTimeCount)
	case "buffer_stage_length":
		return float64(f.BufferStageLength)
	case "buffer_stage_byte_size":
		return float64(f.BufferStageByteSize)
	case "buffer_queue_byte_size":
		return float64(f.BufferQueueByteSize)
	case "buffer_available_buffer_space_ratios":
		return f.BufferAvailableBufferSpaceRatios
	default:
		logger.Warningf("extended-metrics %s not defined", name)
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

func (f *FluentdMetrics) parseStats(body []byte) (map[string]interface{}, error) {
	var j FluentMonitorJSON
	err := json.Unmarshal(body, &j)
	f.plugins = j.Plugins

	metrics := make(map[string]interface{})
	for _, p := range f.plugins {
		if f.nonTargetPlugin(p) {
			continue
		}
		pid := p.getNormalizedPluginID()
		metrics["fluentd.retry_count."+pid] = float64(p.RetryCount)
		metrics["fluentd.buffer_queue_length."+pid] = float64(p.BufferQueueLength)
		metrics["fluentd.buffer_total_queued_size."+pid] = float64(p.BufferTotalQueuedSize)
		for _, name := range f.extendedMetrics {
			key := strings.Join([]string{"fluentd", name, pid}, ".")
			metrics[key] = p.getExtended(name)
		}
	}
	return metrics, err
}

func (f *FluentdMetrics) nonTargetPlugin(plugin FluentdPluginMetrics) bool {
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

// FetchMetrics interface for mackerelplugin
func (f FluentdMetrics) FetchMetrics() (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, f.Target, nil)
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

var defaultGraphs = map[string]mp.Graphs{
	"fluentd.retry_count": {
		Label: "Fluentd retry count",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
	"fluentd.buffer_queue_length": {
		Label: "Fluentd queue length",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
	"fluentd.buffer_total_queued_size": {
		Label: "Fluentd buffer total queued size",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
}

var extendedGraphs = map[string]mp.Graphs{
	"fluentd.emit_records": {
		Label: "Fluentd emitted records",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: true},
		},
	},
	"fluentd.emit_count": {
		Label: "Fluentd emit calls",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: true},
		},
	},
	"fluentd.write_count": {
		Label: "Fluentd write/try_write calls",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: true},
		},
	},
	"fluentd.rollback_count": {
		Label: "Fluentd rollbacks",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: true},
		},
	},
	"fluentd.slow_flush_count": {
		Label: "Fluentd slow flushes",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: true},
		},
	},
	"fluentd.flush_time_count": {
		Label: "Fluentd buffer flush time in msec",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: true},
		},
	},
	"fluentd.buffer_stage_length": {
		Label: "Fluentd length of staged buffer chunks",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
	"fluentd.buffer_stage_byte_size": {
		Label: "Fluentd bytesize of staged buffer chunks",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
	"fluentd.buffer_queue_byte_size": {
		Label: "Fluentd bytesize of queued buffer chunks",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
	"fluentd.buffer_available_buffer_space_ratios": {
		Label: "Fluentd available space for buffer",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "*", Label: "%1", Diff: false},
		},
	},
}

// GraphDefinition interface for mackerelplugin
func (f FluentdMetrics) GraphDefinition() map[string]mp.Graphs {
	graphs := make(map[string]mp.Graphs, len(defaultGraphs))
	for key, g := range defaultGraphs {
		g := g
		graphs[key] = g
	}
	for _, name := range f.extendedMetrics {
		if g, ok := extendedGraphs["fluentd."+name]; ok {
			graphs["fluentd."+name] = g
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
	tempFile := flag.String("tempfile", "", "Temp file name")
	extendedMetricNames := flag.String("extended_metrics", "", "extended metric names joind with ',' or 'all' (fluentd >= v1.6.0)")
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
			extendedMetrics = append(extendedMetrics, strings.TrimPrefix(key, "fluentd."))
		}
	case "":
	default:
		extendedMetrics = strings.Split(*extendedMetricNames, ",")
	}
	f := FluentdMetrics{
		Target:          fmt.Sprintf("http://%s:%s/api/plugins.json", *host, *port),
		Tempfile:        *tempFile,
		pluginType:      *pluginType,
		pluginIDPattern: pluginIDPattern,
		extendedMetrics: extendedMetrics,
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
