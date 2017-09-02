package mpmemcached

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/Songmu/go-memcached-tool"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

// MemcachedPlugin mackerel plugin for memchached
type MemcachedPlugin struct {
	Target    string
	Socket    string
	SlabStats bool
	Tempfile  string
	Prefix    string
}

// MetricKeyPrefix interface for PluginWithPrefix
func (m MemcachedPlugin) MetricKeyPrefix() string {
	if m.Prefix == "" {
		m.Prefix = "memcached"
	}
	return m.Prefix
}

// FetchMetrics interface for mackerelplugin
func (m MemcachedPlugin) FetchMetrics() (map[string]float64, error) {
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

	ret, err := m.parseStats(conn)
	if err != nil {
		return nil, err
	}
	if m.SlabStats {
		slabStats, err := memdtool.GetSlabStats(conn)
		if err != nil {
			log.Printf("failed to get slab stats: %s", err.Error())
		}
		for _, ss := range slabStats {
			ret["nonzero_evictions"] += float64(ss.EvictedNonzero)
			ret[fmt.Sprintf("slab.items.%d", ss.ID)] = float64(ss.Number)
			ret[fmt.Sprintf("slab.evictions.%d.evicted", ss.ID)] = float64(ss.Evicted)
			ret[fmt.Sprintf("slab.evictions.%d.evicted_nonzero", ss.ID)] = float64(ss.EvictedNonzero)
			ret[fmt.Sprintf("slab.evictions.%d.reclaimed", ss.ID)] = float64(ss.Reclaimed)
		}
	}
	return ret, nil
}

func (m MemcachedPlugin) parseStats(conn io.Reader) (map[string]float64, error) {
	scanner := bufio.NewScanner(conn)
	stat := make(map[string]float64)

	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)
		if s == "END" {
			stat["new_items"] = stat["total_items"]
			return stat, nil
		}

		res := strings.Split(s, " ")
		if res[0] == "STAT" {
			f, err := strconv.ParseFloat(res[2], 64)
			if err == nil {
				stat[res[1]] = f
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return stat, err
	}
	return nil, nil
}

// GraphDefinition interface for mackerelplugin
func (m MemcachedPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(m.Prefix)

	// https://github.com/memcached/memcached/blob/master/doc/protocol.txt
	var graphdef = map[string]mp.Graphs{
		"connections": {
			Label: (labelPrefix + " Connections"),
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "curr_connections", Label: "Connections"},
			},
		},
		"cmd": {
			Label: (labelPrefix + " Command"),
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "cmd_get", Label: "Get", Diff: true},
				{Name: "cmd_set", Label: "Set", Diff: true},
				{Name: "cmd_flush", Label: "Flush", Diff: true},
				{Name: "cmd_touch", Label: "Touch", Diff: true},
			},
		},
		"hitmiss": {
			Label: (labelPrefix + " Hits/Misses"),
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "get_hits", Label: "Get Hits", Diff: true},
				{Name: "get_misses", Label: "Get Misses", Diff: true},
				{Name: "delete_hits", Label: "Delete Hits", Diff: true},
				{Name: "delete_misses", Label: "Delete Misses", Diff: true},
				{Name: "incr_hits", Label: "Incr Hits", Diff: true},
				{Name: "incr_misses", Label: "Incr Misses", Diff: true},
				{Name: "cas_hits", Label: "Cas Hits", Diff: true},
				{Name: "cas_misses", Label: "Cas Misses", Diff: true},
				{Name: "touch_hits", Label: "Touch Hits", Diff: true},
				{Name: "touch_misses", Label: "Touch Misses", Diff: true},
			},
		},
		"evictions": {
			Label: (labelPrefix + " Evictions"),
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "evictions", Label: "Evictions", Diff: true},
				{Name: "nonzero_evictions", Label: "Nonzero Evictions", Diff: true},
				{Name: "reclaimed", Label: "Reclaimed", Diff: true},
			},
		},
		"unfetched": {
			Label: (labelPrefix + " Unfetched"),
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "expired_unfetched", Label: "Expired unfetched", Diff: true},
				{Name: "evicted_unfetched", Label: "Evicted unfetched", Diff: true},
			},
		},
		"rusage": {
			Label: (labelPrefix + " Resouce Usage"),
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "rusage_user", Label: "User", Diff: true},
				{Name: "rusage_system", Label: "System", Diff: true},
			},
		},
		"bytes": {
			Label: (labelPrefix + " Traffics"),
			Unit:  mp.UnitBytes,
			Metrics: []mp.Metrics{
				{Name: "bytes_read", Label: "Read", Diff: true},
				{Name: "bytes_written", Label: "Write", Diff: true},
			},
		},
		"cachesize": {
			Label: (labelPrefix + " Cache Size"),
			Unit:  mp.UnitBytes,
			Metrics: []mp.Metrics{
				{Name: "limit_maxbytes", Label: "Total"},
				{Name: "bytes", Label: "Used"},
			},
		},
		"items": {
			Label: (labelPrefix + " Items"),
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "curr_items", Label: "Current Items"},
				{Name: "new_items", Label: "New Items", Diff: true},
			},
		},
		"slab.items": {
			Label: (labelPrefix + "Slab Items"),
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "*", Label: "%1", Stacked: true},
			},
		},
		"slab.evictions.#": {
			Label: (labelPrefix + "Slab Evictions"),
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "evicted", Label: "Evicted Items", Diff: true},
				{Name: "evicted_nonzero", Label: "Nonzero Evicted Items", Diff: true},
				{Name: "reclaimed", Label: "Reclaimed Items", Diff: true},
			},
		},
	}
	return graphdef
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "11211", "Port")
	optSocket := flag.String("socket", "", "Server socket (overrides hosts and port)")
	optSlabStat := flag.Bool("enable-slab-stats", false, "enable slab stat or not")
	optPrefix := flag.String("metric-key-prefix", "memcached", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var memcached MemcachedPlugin

	memcached.Prefix = *optPrefix
	memcached.SlabStats = *optSlabStat

	if *optSocket != "" {
		memcached.Socket = *optSocket
	} else {
		memcached.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	}
	helper := mp.NewMackerelPlugin(memcached)
	helper.Tempfile = *optTempfile
	helper.Run()
}
