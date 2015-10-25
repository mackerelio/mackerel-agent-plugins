package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

var graphdef = map[string](mp.Graphs){
	"varnish.requests": mp.Graphs{
		Label: "Varnish Client Requests",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "requests", Label: "Requests", Diff: true},
			mp.Metrics{Name: "cache_hits", Label: "Hits", Diff: true},
		},
	},
}

// VarnishPlugin mackerel plugin for varnish
type VarnishPlugin struct {
	VarnishStatPath string
	VarnishName     string
	Tempfile        string
}

// FetchMetrics interface for mackerelplugin
func (m VarnishPlugin) FetchMetrics() (map[string]float64, error) {
	var out []byte
	var err error

	if m.VarnishName == "" {
		out, err = exec.Command(m.VarnishStatPath, "-1").CombinedOutput()
	} else {
		out, err = exec.Command(m.VarnishStatPath, "-1", "-n", m.VarnishName).CombinedOutput()
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, out)
	}

	lineexp, err := regexp.Compile("^([^ ]+) +(\\d+)")

	var cacheHits float64
	var cacheMisses float64
	var cacheHitsForPass float64
	for _, line := range strings.Split(string(out), "\n") {
		match := lineexp.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		switch match[1] {
		case "cache_hit", "MAIN.cache_hit":
			cacheHits, err = strconv.ParseFloat(match[2], 64)
		case "cache_miss", "MAIN.cache_miss":
			cacheMisses, err = strconv.ParseFloat(match[2], 64)
		case "cache_hitpass", "MAIN.cache_hitpass":
			cacheHitsForPass, err = strconv.ParseFloat(match[2], 64)
		}
	}

	stat := map[string]float64{
		"requests":   cacheHits + cacheMisses + cacheHitsForPass,
		"cache_hits": cacheHits,
	}
	return stat, err
}

// GraphDefinition interface for mackerelplugin
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
