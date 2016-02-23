package main

import (
	"flag"
	"fmt"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// UptimePlugin mackerel plugin
type UptimePlugin struct {
	Prefix string
}

// GraphDefinition interface for mackerelplugin
func (u UptimePlugin) GraphDefinition() map[string](mp.Graphs) {
	labelPrefix := strings.Title(u.Prefix)
	return map[string](mp.Graphs){
		u.Prefix: mp.Graphs{
			Label: labelPrefix,
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "seconds", Label: "Seconds"},
			},
		},
	}
}

// FetchMetrics interface for mackerelplugin
func (u UptimePlugin) FetchMetrics() (map[string]interface{}, error) {
	return fetchMetrics()
}

func main() {
	optPrefix := flag.String("metric-key-prefix", "uptime", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	u := UptimePlugin{
		Prefix: *optPrefix,
	}
	helper := mp.NewMackerelPlugin(u)
	helper.Tempfile = *optTempfile
	if helper.Tempfile == "" {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-%s", *optPrefix)
	}
	helper.Run()
}
