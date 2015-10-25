package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string](mp.Graphs){
	"plack.workers": mp.Graphs{
		Label: "Plack Workers",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "busy_workers", Label: "Busy Workers", Diff: false, Stacked: true},
			mp.Metrics{Name: "idle_workers", Label: "Idle Workers", Diff: false, Stacked: true},
		},
	},
	"plack.req": mp.Graphs{
		Label: "Plack Requests",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "requests", Label: "Requests", Diff: true, Type: "uint64"},
		},
	},
	"plack.bytes": mp.Graphs{
		Label: "Plack Bytes",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "bytes_sent", Label: "Bytes Sent", Diff: true, Type: "uint64"},
		},
	},
}

// PlackPlugin mackerel plugin for Plack
type PlackPlugin struct {
	URI string
}

// {
//   "Uptime": "1410520211",
//   "TotalAccesses": "2",
//   "IdleWorkers": "2",
//   "TotalKbytes": "0",
//   "BusyWorkers": "1",
//   "stats": [
//     {
//       "pid": 11062,
//       "method": "GET",
//       "ss": 51,
//       "remote_addr": "127.0.0.1",
//       "host": "localhost:8000",
//       "protocol": "HTTP/1.1",
//       "status": "_",
//       "uri": "/server-status?json"
//     },
//     {
//       "ss": 41,
//       "remote_addr": "127.0.0.1",
//       "host": "localhost:8000",
//       "protocol": "HTTP/1.1",
//       "pid": 11063,
//       "method": "GET",
//       "status": "_",
//       "uri": "/server-status?json"
//     },
//     {
//       "ss": 0,
//       "remote_addr": "127.0.0.1",
//       "host": "localhost:8000",
//       "protocol": "HTTP/1.1",
//       "pid": 11064,
//       "method": "GET",
//       "status": "A",
//       "uri": "/server-status?json"
//     }
//   ]
// }

// field types vary between versions

// PlackRequest request
type PlackRequest struct{}

// PlackServerStatus sturct for server-status's json
type PlackServerStatus struct {
	Uptime        string         `json:"Uptime"`
	TotalAccesses string         `json:"TotalAccesses"`
	TotalKbytes   string         `json:"TotalKbytes"`
	BusyWorkers   string         `json:"BusyWorkers"`
	IdleWorkers   string         `json:"IdleWorkers"`
	Stats         []PlackRequest `json:"stats"`
}

// FetchMetrics interface for mackerelplugin
func (p PlackPlugin) FetchMetrics() (map[string]interface{}, error) {
	resp, err := http.Get(p.URI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return p.parseStats(resp.Body)
}

func (p PlackPlugin) parseStats(body io.Reader) (map[string]interface{}, error) {
	stat := make(map[string]interface{})
	decoder := json.NewDecoder(body)

	var s PlackServerStatus
	err := decoder.Decode(&s)
	if err != nil {
		return nil, err
	}

	stat["busy_workers"], err = strconv.ParseFloat(s.BusyWorkers, 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}

	stat["idle_workers"], err = strconv.ParseFloat(s.IdleWorkers, 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}

	stat["requests"], err = strconv.ParseUint(s.TotalAccesses, 10, 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}

	stat["bytes_sent"], err = strconv.ParseUint(s.TotalKbytes, 10, 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}
	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (p PlackPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optURI := flag.String("uri", "", "URI")
	optScheme := flag.String("scheme", "http", "Scheme")
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "5000", "Port")
	optPath := flag.String("path", "/server-status?json", "Path")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var plack PlackPlugin
	if *optURI != "" {
		plack.URI = *optURI
	} else {
		plack.URI = fmt.Sprintf("%s://%s:%s%s", *optScheme, *optHost, *optPort, *optPath)
	}

	helper := mp.NewMackerelPlugin(plack)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-plack")
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
