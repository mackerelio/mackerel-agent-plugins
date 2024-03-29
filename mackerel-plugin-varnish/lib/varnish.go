package mpvarnish

import (
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string]mp.Graphs{
	"varnish.requests": {
		Label: "Varnish Client Requests",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "requests", Label: "Requests", Diff: true},
			{Name: "cache_hits", Label: "Hits", Diff: true},
		},
	},
	"varnish.backend": {
		Label: "Varnish Backend",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "backend_req", Label: "Requests", Diff: true},
			{Name: "backend_conn", Label: "Conn success", Diff: true},
			{Name: "backend_fail", Label: "Conn fail", Diff: true},
			{Name: "backend_reuse", Label: "Conn reuse", Diff: true},
			{Name: "backend_recycle", Label: "Conn recycle", Diff: true},
		},
	},
	"varnish.objects": {
		Label: "Varnish Objects",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "n_object", Label: "object", Diff: false},
			{Name: "n_objectcore", Label: "objectcore", Diff: false},
			{Name: "n_objecthead", Label: "objecthead", Diff: false},
		},
	},
	"varnish.objects_expire": {
		Label: "Varnish Objects Expire",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "n_expired", Label: "expire", Diff: true},
		},
	},
	"varnish.busy_requests": {
		Label: "Varnish Busy Requests",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "busy_sleep", Label: "sleep", Diff: true},
			{Name: "busy_wakeup", Label: "wakeup", Diff: true},
		},
	},
	"varnish.sma.g_alloc.#": {
		Label: "Varnish SMA Allocations",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "g_alloc", Label: "num", Diff: false},
		},
	},
	"varnish.sma.memory.#": {
		Label: "Varnish SMA Memory",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "allocated", Label: "Allocated", Diff: false},
			{Name: "available", Label: "Available", Diff: false},
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
func (m VarnishPlugin) FetchMetrics() (map[string]interface{}, error) {
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

	lineexp := regexp.MustCompile(`^([^ ]+) +(\d+)`)
	smaexp := regexp.MustCompile(`^SMA\.([^\.]+)\.(.+)$`)

	stat := map[string]interface{}{
		"requests": float64(0),
	}

	var tmpv float64
	for _, line := range strings.Split(string(out), "\n") {
		match := lineexp.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		tmpv, err = strconv.ParseFloat(match[2], 64)
		if err != nil {
			continue
		}

		switch match[1] {
		case "cache_hit", "MAIN.cache_hit":
			stat["cache_hits"] = tmpv
			stat["requests"] = stat["requests"].(float64) + tmpv
		case "cache_miss", "MAIN.cache_miss":
			stat["requests"] = stat["requests"].(float64) + tmpv
		case "cache_hitpass", "MAIN.cache_hitpass":
			stat["requests"] = stat["requests"].(float64) + tmpv
		case "MAIN.backend_req":
			stat["backend_req"] = tmpv
		case "MAIN.backend_conn":
			stat["backend_conn"] = tmpv
		case "MAIN.backend_fail":
			stat["backend_fail"] = tmpv
		case "MAIN.backend_reuse":
			stat["backend_reuse"] = tmpv
		case "MAIN.backend_recycle":
			stat["backend_recycle"] = tmpv
		case "MAIN.n_object":
			stat["n_object"] = tmpv
		case "MAIN.n_objectcore":
			stat["n_objectcore"] = tmpv
		case "MAIN.n_expired":
			stat["n_expired"] = tmpv
		case "MAIN.n_objecthead":
			stat["n_objecthead"] = tmpv
		case "MAIN.busy_sleep":
			stat["busy_sleep"] = tmpv
		case "MAIN.busy_wakeup":
			stat["busy_wakeup"] = tmpv
		default:
			smamatch := smaexp.FindStringSubmatch(match[1])
			if smamatch == nil {
				continue
			}
			if smamatch[2] == "g_alloc" {
				stat["varnish.sma.g_alloc."+smamatch[1]+".g_alloc"] = tmpv
			} else if smamatch[2] == "g_bytes" {
				stat["varnish.sma.memory."+smamatch[1]+".allocated"] = tmpv
			} else if smamatch[2] == "g_space" {
				stat["varnish.sma.memory."+smamatch[1]+".available"] = tmpv
			}
		}
	}

	return stat, err
}

// GraphDefinition interface for mackerelplugin
func (m VarnishPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
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
	}

	helper.Run()
}
