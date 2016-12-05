package mpsquid

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string]mp.Graphs{
	"squid.requests": {
		Label: "Squid Client Requests",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "requests", Label: "Requests", Diff: true},
		},
	},
	"squid.cache_hit_ratio.5min": {
		Label: "Squid Client Cache Hit Ratio (5min)",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "request_ratio", Label: "Request Ratio", Diff: false},
			{Name: "byte_ratio", Label: "Byte Ratio", Diff: false},
		},
	},
}

// SquidPlugin mackerel plugin for squid
type SquidPlugin struct {
	Target   string
	Tempfile string
}

// FetchMetrics interface for mackerelplugin
func (m SquidPlugin) FetchMetrics() (map[string]interface{}, error) {
	conn, err := net.Dial("tcp", m.Target)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(conn, "GET cache_object://"+m.Target+"/info HTTP/1.0\n\n")
	scanner := bufio.NewScanner(conn)

	stat := make(map[string]interface{})
	//regexpmap := make(map[string]*regexp.Regexp)
	regexpmap := map[*regexp.Regexp]string{
		regexp.MustCompile("Number of HTTP requests received:\t([0-9]+)"): "requests",
		// version 2
		regexp.MustCompile("Request Hit Ratios:\t5min: ([0-9\\.]+)%"): "request_ratio",
		regexp.MustCompile("Byte Hit Ratios:\t5min: ([0-9\\.]+)%"):    "byte_ratio",
		// version 3
		regexp.MustCompile("Hits as % of all requests:\t5min: ([0-9\\.]+)%"): "request_ratio",
		regexp.MustCompile("Hits as % of bytes sent:\t5min: ([0-9\\.]+)%"):   "byte_ratio",
	}

	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)

		for rexp, key := range regexpmap {
			match := rexp.FindStringSubmatch(s)
			if match == nil {
				continue
			}

			stat[key], err = strconv.ParseFloat(match[1], 64)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	return stat, err
}

// GraphDefinition interface for mackerelplugin
func (m SquidPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "3128", "Port")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var squid SquidPlugin
	squid.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	helper := mp.NewMackerelPlugin(squid)
	helper.Tempfile = *optTempfile

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
