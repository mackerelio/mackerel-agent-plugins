package mph2o

import (
	"encoding/json"
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
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "uptime", Label: "Seconds"},
		},
	},
	"connections": {
		Label: "H2O Connections",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "max-connections", Label: "Max connections"},
			{Name: "connections", Label: "Active connections"},
		},
	},
	"listeners": {
		Label: "H2O Listeners",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "listeners", Label: "Listeners"},
		},
	},
	"worker-threads": {
		Label: "H2O Worker Threads",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "worker-threads", Label: "Worker Threads"},
		},
	},
	"num-sessions": {
		Label: "H2O Sessions",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "num-sessions", Label: "Sessions"},
		},
	},
	"requests": {
		Label: "H2O Requests",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "requests", Label: "In-flight Requests"},
		},
	},
	"status-errors": {
		Label: "H2O Status Errors",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "status-errors_400", Label: "Error 400", Diff: true},
			{Name: "status-errors_403", Label: "Error 403", Diff: true},
			{Name: "status-errors_404", Label: "Error 404", Diff: true},
			{Name: "status-errors_405", Label: "Error 405", Diff: true},
			{Name: "status-errors_416", Label: "Error 416", Diff: true},
			{Name: "status-errors_417", Label: "Error 417", Diff: true},
			{Name: "status-errors_500", Label: "Error 500", Diff: true},
			{Name: "status-errors_502", Label: "Error 502", Diff: true},
			{Name: "status-errors_503", Label: "Error 503", Diff: true},
		},
	},
	"http2-errors": {
		Label: "H2O HTTP2 Errors",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "http2-errors_protocol", Label: "Error protocol", Diff: true},
			{Name: "http2-errors_internal", Label: "Error internal", Diff: true},
			{Name: "http2-errors_flow-control", Label: "Error Flow Control", Diff: true},
			{Name: "http2-errors_settings-timeout", Label: "Error Setting Timeout", Diff: true},
			{Name: "http2-errors_frame-size", Label: "Error Frame Size", Diff: true},
			{Name: "http2-errors_refused-stream", Label: "Error Refused Stream", Diff: true},
			{Name: "http2-errors_cancel", Label: "Error Cancel", Diff: true},
			{Name: "http2-errors_compression", Label: "Error Compression", Diff: true},
			{Name: "http2-errors_connect", Label: "Error Connect", Diff: true},
			{Name: "http2-errors_enhance-your-calm", Label: "Error Enhance Your Calm", Diff: true},
			{Name: "http2-errors_inadequate-security", Label: "Error Inadequate Security", Diff: true},
		},
	},
	"read-closed": {
		Label: "H2O Read Closed",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "http2_read-closed", Label: "Read Closed", Diff: true},
		},
	},
	"write-closed": {
		Label: "H2O Write Closed",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "http2_write-closed", Label: "Write Closed", Diff: true},
		},
	},
	"connect-time": {
		Label: "H2O Connect Time",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "connect-time-0", Label: "0 Percentile"},
			{Name: "connect-time-25", Label: "25 Percentile"},
			{Name: "connect-time-50", Label: "50 Percentile"},
			{Name: "connect-time-75", Label: "75 Percentile"},
			{Name: "connect-time-99", Label: "99 Percentile"},
		},
	},
	"header-time": {
		Label: "H2O Header Time",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "header-time-0", Label: "0 Percentile"},
			{Name: "header-time-25", Label: "25 Percentile"},
			{Name: "header-time-50", Label: "50 Percentile"},
			{Name: "header-time-75", Label: "75 Percentile"},
			{Name: "header-time-99", Label: "99 Percentile"},
		},
	},
	"body-time": {
		Label: "H2O Body Time",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "body-time-0", Label: "0 Percentile"},
			{Name: "body-time-25", Label: "25 Percentile"},
			{Name: "body-time-50", Label: "50 Percentile"},
			{Name: "body-time-75", Label: "75 Percentile"},
			{Name: "body-time-99", Label: "99 Percentile"},
		},
	},
	"request-total-time": {
		Label: "H2O Request Total Time",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "request-total-time-0", Label: "0 Percentile"},
			{Name: "request-total-time-25", Label: "25 Percentile"},
			{Name: "request-total-time-50", Label: "50 Percentile"},
			{Name: "request-total-time-75", Label: "75 Percentile"},
			{Name: "request-total-time-99", Label: "99 Percentile"},
		},
	},
	"process-time": {
		Label: "H2O Process Time",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "process-time-0", Label: "0 Percentile"},
			{Name: "process-time-25", Label: "25 Percentile"},
			{Name: "process-time-50", Label: "50 Percentile"},
			{Name: "process-time-75", Label: "75 Percentile"},
			{Name: "process-time-99", Label: "99 Percentile"},
		},
	},
	"response-time": {
		Label: "H2O Response Time",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "response-time-0", Label: "0 Percentile"},
			{Name: "response-time-25", Label: "25 Percentile"},
			{Name: "response-time-50", Label: "50 Percentile"},
			{Name: "response-time-75", Label: "75 Percentile"},
			{Name: "response-time-99", Label: "99 Percentile"},
		},
	},
	"duration": {
		Label: "H2O Duration",
		Unit:  mp.UnitInteger,
		Metrics: []mp.Metrics{
			{Name: "duration-0", Label: "0 Percentile"},
			{Name: "duration-25", Label: "25 Percentile"},
			{Name: "duration-50", Label: "50 Percentile"},
			{Name: "duration-75", Label: "75 Percentile"},
			{Name: "duration-99", Label: "99 Percentile"},
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

	for k, v := range stat {
		switch k {
		case "server-version", "openssl-version", "current-time", "restart-time", "generation":
		case "requests":
			requests, ok := stat["requests"].([]interface{})
			if !ok {
				return nil, fmt.Errorf("cannot get \"%s\" value", "requests")
			}
			metrics["requests"] = float64(len(requests))
		default:
			f, ok := v.(float64)
			if !ok {
				return nil, fmt.Errorf("cannot get \"%s\" value", k)
			}
			metrics[strings.Replace(k, ".", "_", -1)] = f
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
