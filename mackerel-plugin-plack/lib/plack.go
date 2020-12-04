package mpplack

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// PlackPlugin mackerel plugin for Plack
type PlackPlugin struct {
	URI         string
	Prefix      string
	LabelPrefix string
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
	Uptime        interface{}    `json:"Uptime"` // Plack::Middleware::ServerStatus::Lite 0.35 outputs Uptime as a JSON number, though pre-0.35 outputs it as a JSON string.
	TotalAccesses interface{}    `json:"TotalAccesses"`
	TotalKbytes   interface{}    `json:"TotalKbytes"`
	BusyWorkers   interface{}    `json:"BusyWorkers"`
	IdleWorkers   interface{}    `json:"IdleWorkers"`
	Stats         []PlackRequest `json:"stats"`
}

// FetchMetrics interface for mackerelplugin
func (p PlackPlugin) FetchMetrics() (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, p.URI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "mackerel-plugin-plack")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return p.parseStats(resp.Body)
}

func parseNumber(s interface{}) (float64, error) {
	switch v := s.(type) {
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, nil
		}
	case float64:
		return v, nil
	}
	return 0, fmt.Errorf("failed to parse %v as Number", s)
}

func (p PlackPlugin) parseStats(body io.Reader) (map[string]interface{}, error) {
	stat := make(map[string]interface{})
	decoder := json.NewDecoder(body)

	var s PlackServerStatus
	err := decoder.Decode(&s)
	if err != nil {
		return nil, err
	}

	if s, err := parseNumber(s.BusyWorkers); err == nil {
		stat["busy_workers"] = s
	}

	if s, err := parseNumber(s.IdleWorkers); err == nil {
		stat["idle_workers"] = s
	}

	if s, err := parseNumber(s.TotalAccesses); err == nil {
		stat["requests"] = uint64(s)
	}

	if s, err := parseNumber(s.TotalKbytes); err == nil {
		stat["bytes_sent"] = uint64(s)
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (p PlackPlugin) GraphDefinition() map[string]mp.Graphs {
	var graphdef = map[string]mp.Graphs{
		(p.Prefix + ".workers"): {
			Label: p.LabelPrefix + " Workers",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "busy_workers", Label: "Busy Workers", Diff: false, Stacked: true},
				{Name: "idle_workers", Label: "Idle Workers", Diff: false, Stacked: true},
			},
		},
		(p.Prefix + ".req"): {
			Label: p.LabelPrefix + " Requests",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "requests", Label: "Requests", Diff: true, Type: "uint64"},
			},
		},
		(p.Prefix + ".bytes"): {
			Label: p.LabelPrefix + " Bytes",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "bytes_sent", Label: "Bytes Sent", Diff: true, Type: "uint64"},
			},
		},
	}

	return graphdef
}

// Do the plugin
func Do() {
	optURI := flag.String("uri", "", "URI")
	optScheme := flag.String("scheme", "http", "Scheme")
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "5000", "Port")
	optPath := flag.String("path", "/server-status?json", "Path")
	optPrefix := flag.String("metric-key-prefix", "plack", "Prefix")
	optLabelPrefix := flag.String("metric-label-prefix", "", "Label Prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	plack := PlackPlugin{URI: *optURI, Prefix: *optPrefix, LabelPrefix: *optLabelPrefix}
	if plack.URI == "" {
		plack.URI = fmt.Sprintf("%s://%s:%s%s", *optScheme, *optHost, *optPort, *optPath)
	}
	if plack.LabelPrefix == "" {
		plack.LabelPrefix = strings.Title(plack.Prefix)
	}

	helper := mp.NewMackerelPlugin(plack)
	helper.Tempfile = *optTempfile

	helper.Run()
}
