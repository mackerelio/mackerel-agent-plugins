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

type FluentdMetrics struct {
	Target   string
	Tempfile string

	plugins []FluentdPluginMetrics
}

type FluentdPluginMetrics struct {
	RetryCount            uint64 `json:"retry_count"`
	BufferQueueLength     uint64 `json:"buffer_queue_length"`
	BufferTotalQueuedSize uint64 `json:"buffer_total_queued_size"`
	OutputPlugin          bool   `json:"output_plugin"`
	Type                  string `json:"type"`
	PluginCategory        string `json:"plugin_category"`
	PluginID              string `json:"plugin_id"`
	PluginIDModified      string
}

type FluentMonitorJSON struct {
	Plugins []FluentdPluginMetrics `json:"plugins"`
}

var normalizePluginIDRe = regexp.MustCompile(`[^-a-zA-Z0-9_]`)

func normalizePluginID(str string) string {
	return normalizePluginIDRe.ReplaceAllString(str, "_")
}

func (f *FluentdMetrics) ParseStats(body []byte) (map[string]interface{}, error) {
	var j FluentMonitorJSON
	err := json.Unmarshal(body, &j)
	f.plugins = j.Plugins
	for i, _ := range f.plugins {
		f.plugins[i].PluginIDModified = normalizePluginID(f.plugins[i].PluginID)
	}

	metrics := make(map[string]interface{})
	for _, plugin := range f.plugins {
		metrics["fluentd.retry_count."+plugin.PluginIDModified] = float64(plugin.RetryCount)
		metrics["fluentd.buffer_queue_length."+plugin.PluginIDModified] = float64(plugin.BufferQueueLength)
		metrics["fluentd.buffer_total_queued_size."+plugin.PluginIDModified] = float64(plugin.BufferTotalQueuedSize)
	}
	return metrics, err
}

func (f FluentdMetrics) FetchMetrics() (map[string]interface{}, error) {
	resp, err := http.Get(f.Target)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return f.ParseStats(body)
}

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
