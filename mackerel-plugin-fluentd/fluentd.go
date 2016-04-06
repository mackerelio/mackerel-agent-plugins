package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// FluentdMetrics plugin for fluentd
type FluentdMetrics struct {
	Target   string
	Tempfile string

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
		if p.PluginCategory != "output" {
			continue
		}
		pid := p.getNormalizedPluginID()
		metrics["fluentd.retry_count."+pid] = float64(p.RetryCount)
		metrics["fluentd.buffer_queue_length."+pid] = float64(p.BufferQueueLength)
		metrics["fluentd.buffer_total_queued_size."+pid] = float64(p.BufferTotalQueuedSize)
	}
	return metrics, err
}

// FetchMetrics interface for mackerelplugin
func (f FluentdMetrics) FetchMetrics() (map[string]interface{}, error) {
	resp, err := http.Get(f.Target)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return f.parseStats(body)
}

// GraphDefinition interface for mackerelplugin
func (f FluentdMetrics) GraphDefinition() map[string](mp.Graphs) {
	return map[string](mp.Graphs){
		"fluentd.retry_count": mp.Graphs{
			Label: "Fluentd retry count",
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "*", Label: "%1", Diff: false},
			},
		},
		"fluentd.buffer_queue_length": mp.Graphs{
			Label: "Fluentd queue length",
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "*", Label: "%1", Diff: false},
			},
		},
		"fluentd.buffer_total_queued_size": mp.Graphs{
			Label: "Fluentd buffer total queued size",
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "*", Label: "%1", Diff: false},
			},
		},
	}
}

func main() {
	host := flag.String("host", "localhost", "fluentd monitor_agent port")
	port := flag.String("port", "24220", "fluentd monitor_agent port")
	tempFile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	f := FluentdMetrics{
		Target:   fmt.Sprintf("http://%s:%s/api/plugins.json", *host, *port),
		Tempfile: *tempFile,
	}
	helper := mp.NewMackerelPlugin(f)

	if *tempFile != "" {
		helper.Tempfile = *tempFile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-fluentd-%s-%s", *host, *port)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
