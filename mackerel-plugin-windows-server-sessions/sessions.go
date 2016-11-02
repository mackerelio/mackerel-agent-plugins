package main

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

func getNames() ([]string, error) {
	// WMIC OS GET CSName /FORMAT:CSV
	output, err := exec.Command("WMIC", "OS", "GET", "CSName", "/FORMAT:CSV").Output()
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(strings.NewReader(string(output[1:])))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	var names []string
	dup := make(map[string]bool)
	for _, record := range records[1:] {
		name := strings.TrimSpace(record[1])
		if _, ok := dup[name]; !ok {
			names = append(names, name)
		}
		dup[name] = true
	}
	return names, nil
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
		stat[k+".count"] = uint64(v)
	}
	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (m WindowsServerSessionsPlugin) GraphDefinition() map[string](mp.Graphs) {
	m.names, _ = getNames()

	var metrics []mp.Metrics
	for _, v := range m.names {
		metrics = append(metrics, mp.Metrics{
			Name:    v + ".count",
			Label:   "Windows Server Sessions on " + v,
			Diff:    false,
			Stacked: true,
		})
	}

	return map[string](mp.Graphs){
		"windows-server-sessions": mp.Graphs{
			Label:   "Windows Server Sessions",
			Unit:    "uint64",
			Metrics: metrics,
		},
	}
}

func main() {
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var plugin WindowsServerSessionsPlugin

	helper := mp.NewMackerelPlugin(plugin)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = "/tmp/mackerel-plugin-windows-server-sessions"
	}

	helper.Run()
}
