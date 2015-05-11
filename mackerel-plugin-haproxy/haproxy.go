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
	"time"
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
	Uri      string
	Username string
	Password string
}

func (p HAProxyPlugin) FetchMetrics() (map[string]float64, error) {
	client := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}

	request_uri := p.Uri + ";csv;norefresh"
	req, err := http.NewRequest("GET", request_uri, nil)
	if err != nil {
		return nil, err
	}
	if p.Username != "" {
		req.SetBasicAuth(p.Username, p.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Request failed. Status: %s, URI: %s", resp.Status, request_uri)
	}

	return p.ParseStats(resp.Body)
}

func (p HAProxyPlugin) ParseStats(statsBody io.Reader) (map[string]float64, error) {
	stat := make(map[string]float64)
	reader := csv.NewReader(statsBody)

	for {
		columns, err := reader.Read()
		if err == io.EOF {
			break
		}

		if len(columns) < 60 {
			return nil, errors.New("Length of stats csv is too short. Specifed uri may be wrong.")
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
	optUsername := flag.String("username", "", "Username for Basic Auth")
	optPassword := flag.String("password", "", "Password for Basic Auth")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var haproxy HAProxyPlugin
	if *optUri != "" {
		haproxy.Uri = *optUri
	} else {
		haproxy.Uri = fmt.Sprintf("%s://%s:%s%s", *optScheme, *optHost, *optPort, *optPath)
	}

	if *optUsername != "" {
		haproxy.Username = *optUsername
	}

	if *optPassword != "" {
		haproxy.Password = *optPassword
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
