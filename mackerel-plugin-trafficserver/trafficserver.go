package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string](mp.Graphs){
	"trafficserver.cache": mp.Graphs{
		Label: "Trafficserver Cache Hits/Misses",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "cache_hits", Label: "Hits", Diff: true, Stacked: true, Type: "uint64"},
			mp.Metrics{Name: "cache_misses", Label: "Misses", Diff: true, Stacked: true, Type: "uint64"},
		},
	},
	"trafficserver.http_response_codes": mp.Graphs{
		Label: "Trafficserver HTTP Response Codes",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "http_2xx", Label: "2xx", Diff: true, Stacked: true, Type: "uint64"},
			mp.Metrics{Name: "http_3xx", Label: "3xx", Diff: true, Stacked: true, Type: "uint64"},
			mp.Metrics{Name: "http_4xx", Label: "4xx", Diff: true, Stacked: true, Type: "uint64"},
			mp.Metrics{Name: "http_5xx", Label: "5xx", Diff: true, Stacked: true, Type: "uint64"},
		},
	},
	"trafficserver.connections": mp.Graphs{
		Label: "Trafficserver Current Connections",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "conn_server", Label: "Server"},
			mp.Metrics{Name: "conn_client", Label: "Client"},
		},
	},
}

var metricVarDef = map[string]string{
	"cache_hits":   "proxy.node.cache_total_hits",
	"cache_misses": "proxy.node.cache_total_misses",
	"http_2xx":     "proxy.process.http.2xx_responses",
	"http_3xx":     "proxy.process.http.3xx_responses",
	"http_4xx":     "proxy.process.http.4xx_responses",
	"http_5xx":     "proxy.process.http.5xx_responses",
	"conn_server":  "proxy.node.current_server_connections",
	"conn_client":  "proxy.node.current_client_connections",
}

// TrafficserverPlugin mackerel plugin for apache trafficserver
type TrafficserverPlugin struct {
	Tempfile string
}

// FetchMetrics interface for mackerelplugin
func (m TrafficserverPlugin) FetchMetrics() (map[string]interface{}, error) {
	var err error
	strp, err := getDataWithCommand()
	if err != nil {
		return nil, err
	}

	stat := make(map[string]interface{})
	parseVars(strp, &stat)

	return stat, nil
}

func parseVars(text *string, statp *map[string]interface{}) error {
	stat := *statp

	varMetric := make(map[string]string)
	for metric, varkey := range metricVarDef {
		varMetric[varkey] = metric
	}

	lines := strings.Split(*text, "\n")
	for _, line := range lines {
		factors := strings.Split(line, " ")
		varkey := factors[0]

		if metric, present := varMetric[varkey]; present {
			stat[metric], _ = strconv.ParseUint(factors[1], 10, 64)
		}
	}

	return nil
}

func getDataWithCommand() (*string, error) {
	cmd := exec.Command("traffic_line", "-m", "^proxy")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	str := out.String()
	return &str, nil
}

// GraphDefinition interface for mackerelplugin
func (m TrafficserverPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

var stderrLogger *log.Logger

func getStderrLogger() *log.Logger {
	if stderrLogger == nil {
		stderrLogger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return stderrLogger
}

func main() {
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var trafficserver TrafficserverPlugin

	helper := mp.NewMackerelPlugin(trafficserver)
	helper.Tempfile = *optTempfile

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
