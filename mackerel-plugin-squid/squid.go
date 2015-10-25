package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

var graphdef = map[string](mp.Graphs){
	"squid.requests": mp.Graphs{
		Label: "Squid Client Requests",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "requests", Label: "Requests", Diff: true},
		},
	},
	"squid.cache_hit_ratio.5min": mp.Graphs{
		Label: "Squid Client Cache Hit Ratio (5min)",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "request_ratio", Label: "Request Ratio", Diff: false},
			mp.Metrics{Name: "byte_ratio", Label: "Byte Ratio", Diff: false},
		},
	},
}

// SquidPlugin mackerel plugin for squid
type SquidPlugin struct {
	Target   string
	Tempfile string
}

// FetchMetrics interface for mackerelplugin
func (m SquidPlugin) FetchMetrics() (map[string]float64, error) {
	conn, err := net.Dial("tcp", m.Target)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(conn, "GET cache_object://"+m.Target+"/info HTTP/1.0\n\n")
	scanner := bufio.NewScanner(conn)
	stat := make(map[string]float64)
	regexpmap := make(map[string]*regexp.Regexp)

	regexpmap["requests"], _ = regexp.Compile("Number of HTTP requests received:\t([0-9]+)")
	regexpmap["request_ratio"], _ = regexp.Compile("Request Hit Ratios:\t5min: ([0-9\\.]+)%")
	regexpmap["byte_ratio"], _ = regexp.Compile("Byte Hit Ratios:\t5min: ([0-9\\.]+)%")

	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)

		for key, rexp := range regexpmap {
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
func (m SquidPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "3128", "Port")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var squid SquidPlugin
	squid.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	helper := mp.NewMackerelPlugin(squid)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-squid-%s-%s", *optHost, *optPort)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
