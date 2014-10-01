package main

import (
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"errors"
	"net/http"
	"os"
	"strconv"
	"encoding/json"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"plack.workers": mp.Graphs{
		Label: "Plack Workers",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
            mp.Metrics{ Name: "busy_workers", Label: "Busy Workers", Diff: false, Stacked: true },
            mp.Metrics{ Name: "idle_workers", Label: "Idle Workers", Diff: false, Stacked: true },
		},
	},
    "plack.req": mp.Graphs{
        Label: "Plack Requests",
        Unit: "integer",
        Metrics: [](mp.Metrics){
            mp.Metrics{ Name: "requests", Label: "Requests", Diff: true },
        },
    },
    "plack.bytes": mp.Graphs{
        Label: "Plack Bytes",
        Unit: "integer",
        Metrics: [](mp.Metrics){
            mp.Metrics{ Name: "bytes_sent", Label: "Bytes Sent", Diff: true },
        },
    },
}

type PlackPlugin struct {
	Uri string
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
type PlackRequest struct {}
type PlackServerStatus struct {
	Uptime string          `json:"Uptime"`
	TotalAccesses string   `json:"TotalAccesses"`
	TotalKbytes string     `json:"TotalKbytes"`
	BusyWorkers string     `json:"BusyWorkers"`
	IdleWorkers string     `json:"IdleWorkers"`
	Stats []PlackRequest   `json:"stats"`
}

func (p PlackPlugin) FetchMetrics() (map[string]float64, error) {
	resp, err := http.Get(p.Uri)
    if err != nil {
        return nil, err
    }
	defer resp.Body.Close()

	stat := make(map[string]float64)
	decoder := json.NewDecoder(resp.Body)

	var s PlackServerStatus
	err = decoder.Decode(&s)
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

	stat["requests"], err = strconv.ParseFloat(s.TotalAccesses, 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}

	stat["bytes_sent"], err = strconv.ParseFloat(s.TotalKbytes, 64)
	if err != nil {
		return nil, errors.New("cannot get values")
	}
	return stat, nil
}


func (n PlackPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optUri := flag.String("uri", "", "URI")
	optScheme := flag.String("scheme", "http", "Scheme")
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "5000", "Port")
	optPath := flag.String("path", "/server-status?json", "Path")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var plack PlackPlugin
	if *optUri != "" {
		plack.Uri = *optUri
	} else {
		plack.Uri = fmt.Sprintf("%s://%s:%s%s", *optScheme, *optHost, *optPort, *optPath)
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
