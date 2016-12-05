package mptwemproxy

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// TwemproxyPlugin mackerel plugin
type TwemproxyPlugin struct {
	Address string
	Prefix  string
	Timeout uint
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
	labelPrefix := strings.Title(p.Prefix)

	var graphdef = map[string]mp.Graphs{
		"connections": {
			Label: (labelPrefix + " Connections"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "total_connections", Label: "New Connections", Diff: true},
				{Name: "curr_connections", Label: "Current Connections", Diff: false},
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

	metrics := map[string]interface{}{
		"total_connections": *stats.TotalConnections,
		"curr_connections":  *stats.CurrConnections,
	}

	// NOTE: Each custom metric name contains a wildcard.
	for pName, p := range stats.Pools {
		// A normalized pool name corresponds a wildcard
		np := normalizeMetricName(pName)
		wp := "." + np + "."
		metrics["pool_error"+wp+"client_err"] = *p.ClientErr
		metrics["pool_error"+wp+"server_ejects"] = *p.ServerEjects
		metrics["pool_error"+wp+"forward_error"] = *p.ForwardError
		metrics["pool_client_connections"+wp+"client_eof"] = *p.ClientEOF
		metrics["pool_client_connections"+wp+"client_connections"] = *p.ClientConnections

		for sName, s := range p.Servers {
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
	}
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
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	p := TwemproxyPlugin{
		Address: *optAddress,
		Prefix:  *optPrefix,
		Timeout: *optTimeout,
	}

	helper := mp.NewMackerelPlugin(p)
	helper.Tempfile = *optTempfile
	helper.Run()
}
