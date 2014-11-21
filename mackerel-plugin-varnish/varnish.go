package main

import (
	"errors"
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"varnish.requests": mp.Graphs{
		Label: "Varnish Client Requests",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "requests", Label: "Requests", Diff: true},
			mp.Metrics{Name: "cache_hits", Label: "Hits", Diff: true},
		},
	},
}

type VarnishPlugin struct {
	VarnishStatPath string
	VarnishName     string
	Tempfile        string
}

func (m VarnishPlugin) FetchMetrics() (map[string]float64, error) {
	var out []byte
	var err error

	if m.VarnishName == "" {
		out, err = exec.Command(m.VarnishStatPath, "-1").CombinedOutput()
	} else {
		out, err = exec.Command(m.VarnishStatPath, "-1", "-n", m.VarnishName).CombinedOutput()
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s", err, out))
	}

	lineexp, err := regexp.Compile("^([^ ]+) +(\\d+)")

	stat := make(map[string]float64)
	for _, line := range strings.Split(string(out), "\n") {
		match := lineexp.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		switch match[1] {
		case "client_req", "MAIN.client_req":
			stat["requests"], err = strconv.ParseFloat(match[2], 64)
		case "cache_hit", "MAIN.cache_hit":
			stat["cache_hits"], err = strconv.ParseFloat(match[2], 64)
		}
	}

	return stat, err
}

func (m VarnishPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optVarnishStatPath := flag.String("varnishstat", "/usr/bin/varnishstat", "Path of varnishstat")
	optVarnishName := flag.String("varnish-name", "", "Varnish name")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var varnish VarnishPlugin
	varnish.VarnishStatPath = *optVarnishStatPath
	varnish.VarnishName = *optVarnishName
	helper := mp.NewMackerelPlugin(varnish)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = "/tmp/mackerel-plugin-varnish"
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
