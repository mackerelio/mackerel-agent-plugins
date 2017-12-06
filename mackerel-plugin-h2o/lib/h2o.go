package mph2o

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

const prefix = "h2o"

var graphdef = map[string]mp.Graphs{
	"uptime": {
		Label: "H2O Uptime",
		Unit:  mp.UnitFloat,
		Metrics: []mp.Metrics{
			{Name: "uptime", Label: "Seconds"},
		},
	},
	"connections": {
		Label: "H2O Connections",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "connections", Label: "Active connections"},
			{Name: "max_connections", Label: "Max connections"},
		},
	},
	"listeners": {
		Label: "H2O Listeners",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "listeners", Label: "Listeners"},
		},
	},
	"worker_threads": {
		Label: "H2O Worker Threads",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "worker_threads", Label: "Worker Threads"},
		},
	},
	"num_sessions": {
		Label: "H2O Sessions",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "num_sessions", Label: "Sessions"},
		},
	},
	"requests": {
		Label: "H2O Requests",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "requests", Label: "In-flight Requests"},
		},
	},
	"status_errors": {
		Label: "H2O Status Errors",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "status_errors_503", Label: "Error 503", Diff: true},
			{Name: "status_errors_502", Label: "Error 502", Diff: true},
			{Name: "status_errors_500", Label: "Error 500", Diff: true},
			{Name: "status_errors_417", Label: "Error 417", Diff: true},
			{Name: "status_errors_416", Label: "Error 416", Diff: true},
			{Name: "status_errors_405", Label: "Error 405", Diff: true},
			{Name: "status_errors_404", Label: "Error 404", Diff: true},
			{Name: "status_errors_403", Label: "Error 403", Diff: true},
			{Name: "status_errors_400", Label: "Error 400", Diff: true},
		},
	},
	"http2_errors": {
		Label: "H2O HTTP2 Errors",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "http2_errors_protocol", Label: "Error protocol", Diff: true},
			{Name: "http2_errors_internal", Label: "Error internal", Diff: true},
			{Name: "http2_errors_flow_control", Label: "Error Flow Control", Diff: true},
			{Name: "http2_errors_settings_timeout", Label: "Error Setting Timeout", Diff: true},
			{Name: "http2_errors_frame_size", Label: "Error Frame Size", Diff: true},
			{Name: "http2_errors_refused_stream", Label: "Error Refused Stream", Diff: true},
			{Name: "http2_errors_cancel", Label: "Error Cancel", Diff: true},
			{Name: "http2_errors_compression", Label: "Error Compression", Diff: true},
			{Name: "http2_errors_connect", Label: "Error Connect", Diff: true},
			{Name: "http2_errors_enhance_your_calm", Label: "Error Enhance Your Calm", Diff: true},
			{Name: "http2_errors_inadequate_security", Label: "Error Inadequate Security", Diff: true},
		},
	},
	"read_closed": {
		Label: "H2O Read Closed",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "http2_read_closed", Label: "Read Closed", Diff: true},
		},
	},
	"write_closed": {
		Label: "H2O Write Closed",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "http2_write_closed", Label: "Write Closed", Diff: true},
		},
	},
	"connect_time": {
		Label: "H2O Connect Time",
		Unit:  mp.UnitFloat,
		Metrics: []mp.Metrics{
			{Name: "connect_time_99", Label: "99 Percentile"},
			{Name: "connect_time_75", Label: "75 Percentile"},
			{Name: "connect_time_50", Label: "50 Percentile"},
			{Name: "connect_time_25", Label: "25 Percentile"},
			{Name: "connect_time_0", Label: "0 Percentile"},
		},
	},
	"header_time": {
		Label: "H2O Header Time",
		Unit:  mp.UnitFloat,
		Metrics: []mp.Metrics{
			{Name: "header_time_99", Label: "99 Percentile"},
			{Name: "header_time_75", Label: "75 Percentile"},
			{Name: "header_time_50", Label: "50 Percentile"},
			{Name: "header_time_25", Label: "25 Percentile"},
			{Name: "header_time_0", Label: "0 Percentile"},
		},
	},
	"body_time": {
		Label: "H2O Body Time",
		Unit:  mp.UnitFloat,
		Metrics: []mp.Metrics{
			{Name: "body_time_99", Label: "99 Percentile"},
			{Name: "body_time_75", Label: "75 Percentile"},
			{Name: "body_time_50", Label: "50 Percentile"},
			{Name: "body_time_25", Label: "25 Percentile"},
			{Name: "body_time_0", Label: "0 Percentile"},
		},
	},
	"request_total_time": {
		Label: "H2O Request Total Time",
		Unit:  mp.UnitFloat,
		Metrics: []mp.Metrics{
			{Name: "request_total_time_99", Label: "99 Percentile"},
			{Name: "request_total_time_75", Label: "75 Percentile"},
			{Name: "request_total_time_50", Label: "50 Percentile"},
			{Name: "request_total_time_25", Label: "25 Percentile"},
			{Name: "request_total_time_0", Label: "0 Percentile"},
		},
	},
	"process_time": {
		Label: "H2O Process Time",
		Unit:  mp.UnitFloat,
		Metrics: []mp.Metrics{
			{Name: "process_time_99", Label: "99 Percentile"},
			{Name: "process_time_75", Label: "75 Percentile"},
			{Name: "process_time_50", Label: "50 Percentile"},
			{Name: "process_time_25", Label: "25 Percentile"},
			{Name: "process_time_0", Label: "0 Percentile"},
		},
	},
	"response_time": {
		Label: "H2O Response Time",
		Unit:  mp.UnitFloat,
		Metrics: []mp.Metrics{
			{Name: "response_time_99", Label: "99 Percentile"},
			{Name: "response_time_75", Label: "75 Percentile"},
			{Name: "response_time_50", Label: "50 Percentile"},
			{Name: "response_time_25", Label: "25 Percentile"},
			{Name: "response_time_0", Label: "0 Percentile"},
		},
	},
	"duration": {
		Label: "H2O Duration",
		Unit:  mp.UnitFloat,
		Metrics: []mp.Metrics{
			{Name: "duration_99", Label: "99 Percentile"},
			{Name: "duration_75", Label: "75 Percentile"},
			{Name: "duration_50", Label: "50 Percentile"},
			{Name: "duration_25", Label: "25 Percentile"},
			{Name: "duration_0", Label: "0 Percentile"},
		},
	},
}

type stringSlice []string

func (s *stringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

// H2OPlugin mackerel plugin for H2O
type H2OPlugin struct {
	Prefix string
	URI    string
	Header stringSlice
}

// MetricKeyPrefix interface for mackerelplugin
func (h2o H2OPlugin) MetricKeyPrefix() string {
	if h2o.Prefix == "" {
		h2o.Prefix = prefix
	}
	return h2o.Prefix
}

// GraphDefinition interface for mackerelplugin
func (h2o H2OPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// FetchMetrics interface for mackerelplugin
func (h2o H2OPlugin) FetchMetrics() (map[string]float64, error) {
	req, err := http.NewRequest("GET", h2o.URI, nil)
	if err != nil {
		return nil, err
	}
	for _, h := range h2o.Header {
		kv := strings.SplitN(h, ":", 2)
		var k, v string
		k = strings.TrimSpace(kv[0])
		if len(kv) == 2 {
			v = strings.TrimSpace(kv[1])
		}
		if http.CanonicalHeaderKey(k) == "Host" {
			req.Host = v
		} else {
			req.Header.Set(k, v)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return h2o.parseStats(resp.Body)
}

func (h2o H2OPlugin) parseStats(body io.Reader) (map[string]float64, error) {
	stat := make(map[string]interface{})
	metrics := make(map[string]float64)

	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &stat)
	if err != nil {
		return nil, err
	}

	r := strings.NewReplacer(".", "_", "-", "_")

	for k, v := range stat {
		switch k {
		case "server-version", "openssl-version", "current-time", "restart-time", "generation":
		case "requests":
			requests, ok := stat["requests"].([]interface{})
			if !ok {
				return nil, errors.New("cannot get \"requests\" value")
			}
			metrics["requests"] = float64(len(requests))
		default:
			f, ok := v.(float64)
			if !ok {
				return nil, fmt.Errorf("cannot get %q value", k)
			}
			metrics[r.Replace(k)] = f
		}
	}

	return metrics, nil
}

// Do the plugin
func Do() {
	optURI := flag.String("uri", "", "URI")
	optScheme := flag.String("scheme", "http", "Scheme")
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "8080", "Port")
	optPath := flag.String("path", "/server-status/json", "Path")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optPrefix := flag.String("metric-key-prefix", prefix, "Metric key prefix")
	optHeader := &stringSlice{}
	flag.Var(optHeader, "header", "Set http header (e.g. \"Host: servername\")")
	flag.Parse()

	var h2o = H2OPlugin{
		Prefix: *optPrefix,
		Header: *optHeader,
	}
	if *optURI != "" {
		h2o.URI = *optURI
	} else {
		h2o.URI = fmt.Sprintf("%s://%s:%s%s", *optScheme, *optHost, *optPort, *optPath)
	}

	helper := mp.NewMackerelPlugin(h2o)
	helper.Tempfile = *optTempfile
	helper.Run()
}
