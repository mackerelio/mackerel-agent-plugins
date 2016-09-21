package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// MemcachedPlugin mackerel plugin for memchached
type MemcachedPlugin struct {
	Target   string
	Socket   string
	Tempfile string
	Prefix   string
}

// MetricKeyPrefix interface for PluginWithPrefix
func (m MemcachedPlugin) MetricKeyPrefix() string {
	if m.Prefix == "" {
		m.Prefix = "memcached"
	}
	return m.Prefix
}

// FetchMetrics interface for mackerelplugin
func (m MemcachedPlugin) FetchMetrics() (map[string]interface{}, error) {
	network := "tcp"
	target := m.Target
	if m.Socket != "" {
		network = "unix"
		target = m.Socket
	}
	conn, err := net.Dial(network, target)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(conn, "stats")
	return m.parseStats(conn)
}

func (m MemcachedPlugin) parseStats(conn io.Reader) (map[string]interface{}, error) {
	scanner := bufio.NewScanner(conn)
	stat := make(map[string]interface{})

	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)
		if s == "END" {
			return stat, nil
		}

		res := strings.Split(s, " ")
		if res[0] == "STAT" {
			stat[res[1]] = res[2]
		}
	}
	if err := scanner.Err(); err != nil {
		return stat, err
	}
	return nil, nil
}

// GraphDefinition interface for mackerelplugin
func (m MemcachedPlugin) GraphDefinition() map[string](mp.Graphs) {
	labelPrefix := strings.Title(m.Prefix)

	// https://github.com/memcached/memcached/blob/master/doc/protocol.txt
	var graphdef = map[string](mp.Graphs){
		"connections": mp.Graphs{
			Label: (labelPrefix + " Connections"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "curr_connections", Label: "Connections", Diff: false},
			},
		},
		"cmd": mp.Graphs{
			Label: (labelPrefix + " Command"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "cmd_get", Label: "Get", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "cmd_set", Label: "Set", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "cmd_flush", Label: "Flush", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "cmd_touch", Label: "Touch", Diff: true, Type: "uint64"},
			},
		},
		"hitmiss": mp.Graphs{
			Label: (labelPrefix + " Hits/Misses"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "get_hits", Label: "Get Hits", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "get_misses", Label: "Get Misses", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "delete_hits", Label: "Delete Hits", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "delete_misses", Label: "Delete Misses", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "incr_hits", Label: "Incr Hits", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "incr_misses", Label: "Incr Misses", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "cas_hits", Label: "Cas Hits", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "cas_misses", Label: "Cas Misses", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "touch_hits", Label: "Touch Hits", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "touch_misses", Label: "Touch Misses", Diff: true, Type: "uint64"},
			},
		},
		"evictions": mp.Graphs{
			Label: (labelPrefix + " Evictions"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "evictions", Label: "Evictions", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "evicted_nonzero", Label: "Evictions prior to Expire", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "reclaimed", Label: "Reclaimed Items", Diff: true, Type: "uint64"},
			},
		},
		"unfetched": mp.Graphs{
			Label: (labelPrefix + " Unfetched"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "expired_unfetched", Label: "Expired unfetched", Diff: true, Type: "uint64"},
			},
		},
		"rusage": mp.Graphs{
			Label: (labelPrefix + " Resouce Usage"),
			Unit:  "float",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "rusage_user", Label: "User", Diff: true},
				mp.Metrics{Name: "rusage_system", Label: "System", Diff: true},
			},
		},
		"bytes": mp.Graphs{
			Label: (labelPrefix + " Traffics"),
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "bytes_read", Label: "Read", Diff: true, Type: "uint64"},
				mp.Metrics{Name: "bytes_written", Label: "Write", Diff: true, Type: "uint64"},
			},
		},
		"cachesize": mp.Graphs{
			Label: (labelPrefix + " Cache Size"),
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "limit_maxbytes", Label: "Total", Diff: false},
				mp.Metrics{Name: "bytes", Label: "Used", Diff: false, Type: "uint64"},
			},
		},
	}
	return graphdef
}

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "11211", "Port")
	optSocket := flag.String("socket", "", "Server socket (overrides hosts and port)")
	optPrefix := flag.String("metric-key-prefix", "memcached", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var memcached MemcachedPlugin

	memcached.Prefix = *optPrefix

	if *optSocket != "" {
		memcached.Socket = *optSocket
	} else {
		memcached.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	}
	helper := mp.NewMackerelPlugin(memcached)
	helper.Tempfile = *optTempfile
	helper.Run()
}
