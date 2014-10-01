package main

import (
	"bufio"
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"net"
	"os"
	"strconv"
	"regexp"
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
	Target   string
	Tempfile string
}

func (m VarnishPlugin) FetchMetrics() (map[string]float64, error) {
	conn, err := net.Dial("tcp", m.Target)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(conn, "stats")
	scanner := bufio.NewScanner(conn)
	stat := make(map[string]float64)

	metricexp, err := regexp.Compile("^ *([0-9]+)  (.*)$")
	fetched := false
	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)

		match := metricexp.FindStringSubmatch(s)
		if match == nil {
		    // after metric rows
		    if fetched {
			break
		    }
		    continue
		}

		fetched = true

		switch match[2] {
		case "Client requests received":
		    stat["requests"], err = strconv.ParseFloat(match[1], 64)
		case "Cache hits":
		    stat["cache_hits"], err = strconv.ParseFloat(match[1], 64)
		}
	}

	return stat, err
}

func (m VarnishPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "6082", "Port")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var varnish VarnishPlugin
	varnish.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	helper := mp.NewMackerelPlugin(varnish)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-varnish-%s-%s", *optHost, *optPort)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
