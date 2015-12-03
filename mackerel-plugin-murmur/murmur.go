package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	mumble "github.com/layeh/gumble/gumble"
)

var graphdef = map[string](mp.Graphs){
	"murmur.connections": mp.Graphs{
		Label: "Murmur Connections",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "con_cur", Label: "Current users", Diff: false, Type: "uint32"},
			mp.Metrics{Name: "con_max", Label: "Maximum users", Diff: false, Type: "uint32"},
		},
	},
}

// MurmurPlugin mackerel plugin for Murmur
type MurmurPlugin struct {
	Host    string
	Timeout uint64
}

// FetchMetrics interface for mackerelplugin
func (m MurmurPlugin) FetchMetrics() (map[string]interface{}, error) {
	resp, err := mumble.Ping(m.Host, time.Millisecond * time.Duration(m.Timeout))

	if err != nil {
		return nil, err
	}

	metrics := map[string]interface{}{
		"con_cur":uint32(resp.ConnectedUsers),
		"con_max":uint32(resp.MaximumUsers),
	}

	return metrics, nil
}

// GraphDefinition interface for mackerelplugin
func (m MurmurPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "64738", "Port")
	optTimeout := flag.Uint64("timeout", 1000, "Timeout (ms)")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var murmur MurmurPlugin

	murmur.Host = fmt.Sprintf("%s:%s", *optHost, *optPort)
	murmur.Timeout = *optTimeout
	helper := mp.NewMackerelPlugin(murmur)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-murmur-%s-%s", *optHost, *optPort)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}