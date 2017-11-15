package mpuptime

import (
	"flag"
	"fmt"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/mackerelio/golib/uptime"
)

// UptimePlugin mackerel plugin
type UptimePlugin struct {
	Prefix string
}

// MetricKeyPrefix interface for PluginWithPrefix
func (u UptimePlugin) MetricKeyPrefix() string {
	if u.Prefix == "" {
		u.Prefix = "uptime"
	}
	return u.Prefix
}

// GraphDefinition interface for mackerelplugin
func (u UptimePlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(u.Prefix)
	return map[string]mp.Graphs{
		"": {
			Label: labelPrefix,
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "seconds", Label: "Seconds"},
			},
		},
	}
}

// FetchMetrics interface for mackerelplugin
func (u UptimePlugin) FetchMetrics() (map[string]float64, error) {
	ut, err := uptime.Get()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch uptime metrics: %s", err)
	}
	return map[string]float64{"seconds": ut}, nil
}

// Do the plugin
func Do() {
	optPrefix := flag.String("metric-key-prefix", "uptime", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	u := UptimePlugin{
		Prefix: *optPrefix,
	}
	helper := mp.NewMackerelPlugin(u)
	helper.Tempfile = *optTempfile
	helper.Run()
}
