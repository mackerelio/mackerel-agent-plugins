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
)

// FluentdPlugin mackerel plugin
type FluentdPlugin struct {
	Target          string
	Prefix          string
	Tempfile        string
	pluginType      string
	pluginIDPattern *regexp.Regexp

	plugins []FluentdPluginMetrics
}

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
		metrics["retry_count."+pid] = float64(p.RetryCount)
		metrics["buffer_queue_length."+pid] = float64(p.BufferQueueLength)
		metrics["buffer_total_queued_size."+pid] = float64(p.BufferTotalQueuedSize)
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

// FetchMetrics interface for mackerelplugin
func (f FluentdPlugin) FetchMetrics() (map[string]interface{}, error) {
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

// GraphDefinition interface for mackerelplugin
func (f FluentdPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(f.Prefix)

	return map[string]mp.Graphs{
		"retry_count": {
			Label: (labelPrefix + " retry count"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "*", Label: "%1", Diff: false},
			},
		},
		"buffer_queue_length": {
			Label: (labelPrefix + " queue length"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "*", Label: "%1", Diff: false},
			},
		},
		"buffer_total_queued_size": {
			Label: (labelPrefix + " buffer total queued size"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "*", Label: "%1", Diff: false},
			},
		},
	}
}

// Do the plugin
func Do() {
	host := flag.String("host", "localhost", "fluentd monitor_agent host")
	port := flag.String("port", "24220", "fluentd monitor_agent port")
	pluginType := flag.String("plugin-type", "", "Gets the metric that matches this plugin type")
	pluginIDPatternString := flag.String("plugin-id-pattern", "", "Gets the metric that matches this plugin id pattern")
	prefix := flag.String("metric-key-prefix", "fluentd", "Metric key prefix")
	tempFile := flag.String("tempfile", "", "Temp file name")
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

	f := FluentdPlugin{
		Target:          fmt.Sprintf("http://%s:%s/api/plugins.json", *host, *port),
		Prefix:          *prefix,
		Tempfile:        *tempFile,
		pluginType:      *pluginType,
		pluginIDPattern: pluginIDPattern,
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
