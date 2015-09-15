package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
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

type TrafficserverPlugin struct {
	Tempfile string
}

func (m TrafficserverPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	var err error
	stat["cache_hits"], err = m.CollectValue("proxy.node.cache_total_hits")
	if err != nil {
		getStderrLogger().Println(err)
	}
	stat["cache_misses"], err = m.CollectValue("proxy.node.cache_total_misses")
	if err != nil {
		getStderrLogger().Println(err)
	}
	stat["http_2xx"], err = m.CollectValue("proxy.process.http.2xx_responses")
	if err != nil {
		getStderrLogger().Println(err)
	}
	stat["http_3xx"], err = m.CollectValue("proxy.process.http.3xx_responses")
	if err != nil {
		getStderrLogger().Println(err)
	}
	stat["http_4xx"], err = m.CollectValue("proxy.process.http.4xx_responses")
	if err != nil {
		getStderrLogger().Println(err)
	}
	stat["http_5xx"], err = m.CollectValue("proxy.process.http.5xx_responses")
	if err != nil {
		getStderrLogger().Println(err)
	}
	stat["conn_server"], err = m.CollectValue("proxy.node.current_server_connections")
	if err != nil {
		getStderrLogger().Println(err)
	}
	stat["conn_client"], err = m.CollectValue("proxy.node.current_client_connections")
	if err != nil {
		getStderrLogger().Println(err)
	}

	return stat, nil
}

func (m TrafficserverPlugin) CollectValue(key string) (uint64, error) {
	cmd := exec.Command("traffic_line", "-r", key)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	ret, err := strconv.ParseUint(strings.TrimRight(out.String(), "\n"), 10, 64)
	if err != nil {
		return 0, err
	}

	return ret, nil
}

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

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-trafficserver")
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
