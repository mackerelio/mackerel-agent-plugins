package mpwindowsserversessions

import (
	"encoding/csv"
	"flag"
	"os/exec"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metrics.plugin.windows-server-sessions")

// WindowsServerSessionsPlugin store the name of servers
type WindowsServerSessionsPlugin struct {
	names []string
}

func getCounts() (map[string]int, error) {
	// WMIC PATH Win32_PerfFormattedData_PerfNet_Server GET ServerSessions /FORMAT:CSV
	output, err := exec.Command("WMIC", "PATH", "Win32_PerfFormattedData_PerfNet_Server", "GET", "ServerSessions", "/FORMAT:CSV").Output()
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(strings.NewReader(string(output[1:])))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, record := range records[1:] {
		name := strings.TrimSpace(record[0])
		n, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			continue
		}
		counts[name] = n
	}
	return counts, nil
}

// FetchMetrics interface for mackerelplugin
func (m WindowsServerSessionsPlugin) FetchMetrics() (map[string]interface{}, error) {
	counts, err := getCounts()
	if err != nil {
		return nil, err
	}
	stat := make(map[string]interface{})
	for k, v := range counts {
		stat["windows.server.sessions."+k+".count"] = uint64(v)
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
