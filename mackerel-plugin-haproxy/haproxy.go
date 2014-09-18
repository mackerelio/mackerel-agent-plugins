package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"io"
	"net/http"
	"os"
	"strconv"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"haproxy.total.sessions": mp.Graphs{
		Label: "HAProxy Total Sessions",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "sessions", Label: "Sessions", Diff: true},
		},
	},
	"haproxy.total.bytes": mp.Graphs{
		Label: "HAProxy Total Bytes",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "bytes_in", Label: "Bytes In", Diff: true},
			mp.Metrics{Name: "bytes_out", Label: "Bytes Out", Diff: true},
		},
	},
	"haproxy.total.connection_errors": mp.Graphs{
		Label: "HAProxy Total Connection Errors",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "connection_errors", Label: "Connection Errors", Diff: true},
		},
	},
}

type HAProxyPlugin struct {
	Uri string
}

func (p HAProxyPlugin) FetchMetrics() (map[string]float64, error) {
	resp, err := http.Get(p.Uri + ";csv;norefresh")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	stat := make(map[string]float64)
	reader := csv.NewReader(resp.Body)

	for {
		columns, err := reader.Read()

		if err == io.EOF {
			break
		}

		if columns[1] != "BACKEND" {
			continue
		}

		var data float64

		data, err = strconv.ParseFloat(columns[7], 64)
		if err != nil {
			return nil, errors.New("cannot get values")
		}
		stat["sessions"] += data

		data, err = strconv.ParseFloat(columns[8], 64)
		if err != nil {
			return nil, errors.New("cannot get values")
		}
		stat["bytes_in"] += data

		data, err = strconv.ParseFloat(columns[9], 64)
		if err != nil {
			return nil, errors.New("cannot get values")
		}
		stat["bytes_out"] += data

		data, err = strconv.ParseFloat(columns[13], 64)
		if err != nil {
			return nil, errors.New("cannot get values")
		}
		stat["connection_errors"] += data
	}

	return stat, nil
}

func (n HAProxyPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optUri := flag.String("uri", "", "URI")
	optScheme := flag.String("scheme", "http", "Scheme")
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "80", "Port")
	optPath := flag.String("path", "/", "Path")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var haproxy HAProxyPlugin
	if *optUri != "" {
		haproxy.Uri = *optUri
	} else {
		haproxy.Uri = fmt.Sprintf("%s://%s:%s%s", *optScheme, *optHost, *optPort, *optPath)
	}

	helper := mp.NewMackerelPlugin(haproxy)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-haproxy")
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
