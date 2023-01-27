package mptwemproxy

import (
	"flag"
	"fmt"
	"regexp"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TwemproxyPlugin mackerel plugin
type TwemproxyPlugin struct {
	Address           string
	Prefix            string
	Timeout           uint
	EachServerMetrics bool
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p TwemproxyPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "twemproxy"
	}
	return p.Prefix
}

// GraphDefinition interface for mackerelplugin
func (p TwemproxyPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := cases.Title(language.Und, cases.NoLower).String(p.Prefix)

	var graphdef = map[string]mp.Graphs{
		"connections": {
			Label: (labelPrefix + " Connections"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "total_connections", Label: "New Connections", Diff: true},
				{Name: "curr_connections", Label: "Current Connections", Diff: false},
			},
		},
		"total_server_error": {
			Label: (labelPrefix + " Total Server Error"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "total_pool_client_error", Label: "Pool Client Error", Diff: true},
				{Name: "total_pool_server_ejects", Label: "Pool Server Ejects", Diff: true},
				{Name: "total_pool_forward_error", Label: "Pool Forward Error", Diff: true},
				{Name: "total_server_timeout", Label: "Server Error", Diff: true},
				{Name: "total_server_error", Label: "Server Timeout", Diff: true},
			},
		},
		"pool_error.#": {
			Label: (labelPrefix + " Pool Error"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "client_err", Label: "Client Error", Diff: true},
				{Name: "server_ejects", Label: "Server Ejects", Diff: true},
				{Name: "forward_error", Label: "Forward Error", Diff: true},
			},
		},
		"pool_client_connections.#": {
			Label: (labelPrefix + " Pool Client Connections"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "client_connections", Label: "Client Connections", Diff: false},
				{Name: "client_eof", Label: "Client EOF", Diff: true},
			},
		},
		"server_error.#": {
			Label: (labelPrefix + " Server Error"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "server_err", Label: "Server Error", Diff: true},
				{Name: "server_timedout", Label: "Server Timedout", Diff: true},
			},
		},
		"server_connections.#": {
			Label: (labelPrefix + " Server Connections"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "server_connections", Label: "Server Connections", Diff: false},
				{Name: "server_eof", Label: "Server EOF", Diff: true},
			},
		},
		"server_queue.#": {
			Label: (labelPrefix + " Server Queue"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "out_queue", Label: "Out Queue", Diff: false},
				{Name: "in_queue", Label: "In Queue", Diff: false},
			},
		},
		"server_queue_bytes.#": {
			Label: (labelPrefix + " Server Queue Bytes"),
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "out_queue_bytes", Label: "Out Queue Bytes", Diff: false},
				{Name: "in_queue_bytes", Label: "In Queue Bytes", Diff: false},
			},
		},
		"server_communications.#": {
			Label: (labelPrefix + " Server Communications"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "requests", Label: "Requests", Diff: true},
				{Name: "responses", Label: "Responses", Diff: true},
			},
		},
		"server_communication_bytes.#": {
			Label: (labelPrefix + " Server Communication Bytes"),
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "request_bytes", Label: "Request Bytes", Diff: true},
				{Name: "response_bytes", Label: "Response Bytes", Diff: true},
			},
		},
	}
	return graphdef
}

// FetchMetrics interface for mackerelplugin
func (p TwemproxyPlugin) FetchMetrics() (map[string]interface{}, error) {
	stats, err := getStats(p)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch twemproxy metrics: %s", err)
	}

	metrics := make(map[string]interface{})

	if stats.TotalConnections != nil {
		metrics["total_connections"] = *stats.TotalConnections
	}
	if stats.CurrConnections != nil {
		metrics["curr_connections"] = *stats.CurrConnections
	}

	totalPoolClientErr := uint64(0)
	totalPoolServerEjects := uint64(0)
	totalPoolForwardErr := uint64(0)
	totalServerTimeout := uint64(0)
	totalServerErr := uint64(0)

	// NOTE: Each custom metric name contains a wildcard.
	for pName, po := range stats.Pools {
		// A normalized pool name corresponds a wildcard
		np := normalizeMetricName(pName)
		wp := "." + np + "."
		metrics["pool_error"+wp+"client_err"] = *po.ClientErr
		metrics["pool_error"+wp+"server_ejects"] = *po.ServerEjects
		metrics["pool_error"+wp+"forward_error"] = *po.ForwardError
		metrics["pool_client_connections"+wp+"client_eof"] = *po.ClientEOF
		metrics["pool_client_connections"+wp+"client_connections"] = *po.ClientConnections
		totalPoolClientErr += *po.ClientErr
		totalPoolServerEjects += *po.ServerEjects
		totalPoolForwardErr += *po.ForwardError

		for sName, s := range po.Servers {
			if p.EachServerMetrics {
				// A concat of normalized pool and server names corresponds a wildcard
				ns := normalizeMetricName(sName)
				ws := "." + np + "_" + ns + "."
				metrics["server_error"+ws+"server_err"] = *s.ServerErr
				metrics["server_error"+ws+"server_timedout"] = *s.ServerTimedout
				metrics["server_connections"+ws+"server_eof"] = *s.ServerEOF
				metrics["server_connections"+ws+"server_connections"] = *s.ServerConnections
				metrics["server_queue"+ws+"out_queue"] = *s.OutQueue
				metrics["server_queue"+ws+"in_queue"] = *s.InQueue
				metrics["server_queue_bytes"+ws+"out_queue_bytes"] = *s.OutQueueBytes
				metrics["server_queue_bytes"+ws+"in_queue_bytes"] = *s.InQueueBytes
				metrics["server_communications"+ws+"requests"] = *s.Requests
				metrics["server_communications"+ws+"responses"] = *s.Responses
				metrics["server_communication_bytes"+ws+"request_bytes"] = *s.RequestBytes
				metrics["server_communication_bytes"+ws+"response_bytes"] = *s.ResponseBytes
			}
			totalServerTimeout += *s.ServerTimedout
			totalServerErr += *s.ServerErr
		}
	}

	metrics["total_pool_client_error"] = totalPoolClientErr
	metrics["total_pool_server_ejects"] = totalPoolServerEjects
	metrics["total_pool_forward_error"] = totalPoolForwardErr
	metrics["total_server_timeout"] = totalServerTimeout
	metrics["total_server_error"] = totalServerErr
	return metrics, nil
}

var normalizeMetricNameRe = regexp.MustCompile(`[^-a-zA-Z0-9_]`)

func normalizeMetricName(name string) string {
	return normalizeMetricNameRe.ReplaceAllString(name, "_")
}

// Do the plugin
func Do() {
	optAddress := flag.String("address", "localhost:22222", "twemproxy stats Address")
	optPrefix := flag.String("metric-key-prefix", "twemproxy", "Metric key prefix")
	optTimeout := flag.Uint("timeout", 5, "Timeout")
	optEachServerMetrics := flag.Bool("enable-each-server-metrics", false, "Enable metric collection for each server")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	p := TwemproxyPlugin{
		Address:           *optAddress,
		Prefix:            *optPrefix,
		Timeout:           *optTimeout,
		EachServerMetrics: *optEachServerMetrics,
	}

	helper := mp.NewMackerelPlugin(p)
	helper.Tempfile = *optTempfile
	helper.Run()
}
