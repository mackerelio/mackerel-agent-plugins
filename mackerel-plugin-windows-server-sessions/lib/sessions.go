package mpwindowsserversessions

import (
	"flag"
	"os"
	"strings"

	"github.com/yusufpapurcu/wmi"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// WindowsServerSessionsPlugin store the name of servers
type WindowsServerSessionsPlugin struct {
}

// cf.) https://learn.microsoft.com/en-us/previous-versions/aa394265(v=vs.85)
type Win32_PerfFormattedData_PerfNet_Server struct {
	ServerSessions uint32
}

type ServerSessionsCount struct {
	Node           string
	ServerSessions uint32
}

func getCounts() ([]ServerSessionsCount, error) {
	var dst []Win32_PerfFormattedData_PerfNet_Server
	q := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(q, &dst); err != nil {
		return nil, err
	}

	if len(dst) == 0 {
		return []ServerSessionsCount{}, nil
	}

	node, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	counts := []ServerSessionsCount{
		{Node: node, ServerSessions: dst[0].ServerSessions},
	}

	return counts, nil
}

// FetchMetrics interface for mackerelplugin
func (m WindowsServerSessionsPlugin) FetchMetrics() (map[string]interface{}, error) {
	counts, err := getCounts()
	if err != nil {
		return nil, err
	}
	stat := make(map[string]interface{}, len(counts))
	for _, v := range counts {
		// node name of Windows can contain ".", which is the metric name delimiter on Mackerel.
		nodeMetricKey := strings.ReplaceAll(v.Node, ".", "_")
		stat["windows.server.sessions."+nodeMetricKey+".count"] = v.ServerSessions
	}
	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (m WindowsServerSessionsPlugin) GraphDefinition() map[string](mp.Graphs) {
	return map[string](mp.Graphs){
		"windows.server.sessions.#": mp.Graphs{
			Label: "Windows Server Sessions",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "count", Label: "count", Diff: false, Stacked: false},
			},
		},
	}
}

// Do the plugin
func Do() {
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var plugin WindowsServerSessionsPlugin

	helper := mp.NewMackerelPlugin(plugin)

	helper.Tempfile = *optTempfile
	helper.Run()
}
