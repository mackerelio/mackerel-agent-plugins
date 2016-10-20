package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metrics.plugin.elasticsearch")

var metricPlace = map[string][]string{
	"http_opened":                 {"http", "total_opened"},
	"total_indexing_index":        {"indices", "indexing", "index_total"},
	"total_indexing_delete":       {"indices", "indexing", "delete_total"},
	"total_get":                   {"indices", "get", "total"},
	"total_search_query":          {"indices", "search", "query_total"},
	"total_search_fetch":          {"indices", "search", "fetch_total"},
	"total_merges":                {"indices", "merges", "total"},
	"total_refresh":               {"indices", "refresh", "total"},
	"total_flush":                 {"indices", "flush", "total"},
	"total_warmer":                {"indices", "warmer", "total"},
	"total_percolate":             {"indices", "percolate", "total"},
	"total_suggest":               {"indices", "suggest", "total"},
	"docs_count":                  {"indices", "docs", "count"},
	"docs_deleted":                {"indices", "docs", "deleted"},
	"fielddata_size":              {"indices", "fielddata", "memory_size_in_bytes"},
	"filter_cache_size":           {"indices", "filter_cache", "memory_size_in_bytes"},
	"segments_size":               {"indices", "segments", "memory_in_bytes"},
	"segments_index_writer_size":  {"indices", "segments", "index_writer_memory_in_bytes"},
	"segments_version_map_size":   {"indices", "segments", "version_map_memory_in_bytes"},
	"segments_fixed_bit_set_size": {"indices", "segments", "fixed_bit_set_memory_in_bytes"},
	"evictions_fielddata":         {"indices", "fielddata", "evictions"},
	"evictions_filter_cache":      {"indices", "filter_cache", "evictions"},
	"heap_used":                   {"jvm", "mem", "heap_used_in_bytes"},
	"heap_max":                    {"jvm", "mem", "heap_max_in_bytes"},
	"threads_generic":             {"thread_pool", "generic", "threads"},
	"threads_index":               {"thread_pool", "index", "threads"},
	"threads_snapshot_data":       {"thread_pool", "snapshot_data", "threads"},
	"threads_get":                 {"thread_pool", "get", "threads"},
	"threads_bench":               {"thread_pool", "bench", "threads"},
	"threads_snapshot":            {"thread_pool", "snapshot", "threads"},
	"threads_merge":               {"thread_pool", "merge", "threads"},
	"threads_suggest":             {"thread_pool", "suggest", "threads"},
	"threads_bulk":                {"thread_pool", "bulk", "threads"},
	"threads_optimize":            {"thread_pool", "optimize", "threads"},
	"threads_warmer":              {"thread_pool", "warmer", "threads"},
	"threads_flush":               {"thread_pool", "flush", "threads"},
	"threads_search":              {"thread_pool", "search", "threads"},
	"threads_percolate":           {"thread_pool", "percolate", "threads"},
	"threads_refresh":             {"thread_pool", "refresh", "threads"},
	"threads_management":          {"thread_pool", "management", "threads"},
	"threads_fetch_shard_started": {"thread_pool", "fetch_shard_started", "threads"},
	"threads_fetch_shard_store":   {"thread_pool", "fetch_shard_store", "threads"},
	"threads_listener":            {"thread_pool", "listener", "threads"},
	"count_rx":                    {"transport", "rx_count"},
	"count_tx":                    {"transport", "tx_count"},
	"open_file_descriptors":       {"process", "open_file_descriptors"},
}

func getFloatValue(s map[string]interface{}, keys []string) (float64, error) {
	var val float64
	sm := s
	for i, k := range keys {
		if i+1 < len(keys) {
			switch sm[k].(type) {
			case map[string]interface{}:
				sm = sm[k].(map[string]interface{})
			default:
				return 0, errors.New("Cannot handle as a hash")
			}
		} else {
			switch sm[k].(type) {
			case float64:
				val = sm[k].(float64)
			default:
				return 0, errors.New("Not float64")
			}
		}
	}

	return val, nil
}

// ElasticsearchPlugin mackerel plugin for Elasticsearch
type ElasticsearchPlugin struct {
	URI         string
	Prefix      string
	LabelPrefix string
}

// FetchMetrics interface for mackerelplugin
func (p ElasticsearchPlugin) FetchMetrics() (map[string]float64, error) {
	resp, err := http.Get(p.URI + "/_nodes/_local/stats")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	stat := make(map[string]float64)
	decoder := json.NewDecoder(resp.Body)

	var s map[string]interface{}
	err = decoder.Decode(&s)
	if err != nil {
		return nil, err
	}

	nodes := s["nodes"].(map[string]interface{})
	n := ""
	for k := range nodes {
		if n != "" {
			return nil, errors.New("Multiple node found")
		}
		n = k
	}
	node := nodes[n].(map[string]interface{})

	for k, v := range metricPlace {
		val, err := getFloatValue(node, v)
		if err != nil {
			logger.Errorf("Failed to find '%s': %s", k, err)
			continue
		}

		stat[k] = val
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (p ElasticsearchPlugin) GraphDefinition() map[string](mp.Graphs) {
	var graphdef = map[string](mp.Graphs){
		p.Prefix + ".http": mp.Graphs{
			Label: (p.LabelPrefix + " HTTP"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "http_opened", Label: "Opened", Diff: true},
			},
		},
		p.Prefix + ".indices": mp.Graphs{
			Label: (p.LabelPrefix + " Indices"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "total_indexing_index", Label: "Indexing-Index", Diff: true, Stacked: true},
				mp.Metrics{Name: "total_indexing_delete", Label: "Indexing-Delete", Diff: true, Stacked: true},
				mp.Metrics{Name: "total_get", Label: "Get", Diff: true, Stacked: true},
				mp.Metrics{Name: "total_search_query", Label: "Search-Query", Diff: true, Stacked: true},
				mp.Metrics{Name: "total_search_fetch", Label: "Search-fetch", Diff: true, Stacked: true},
				mp.Metrics{Name: "total_merges", Label: "Merges", Diff: true, Stacked: true},
				mp.Metrics{Name: "total_refresh", Label: "Refresh", Diff: true, Stacked: true},
				mp.Metrics{Name: "total_flush", Label: "Flush", Diff: true, Stacked: true},
				mp.Metrics{Name: "total_warmer", Label: "Warmer", Diff: true, Stacked: true},
				mp.Metrics{Name: "total_percolate", Label: "Percolate", Diff: true, Stacked: true},
				mp.Metrics{Name: "total_suggest", Label: "Suggest", Diff: true, Stacked: true},
			},
		},
		p.Prefix + ".indices.docs": mp.Graphs{
			Label: (p.LabelPrefix + " Indices Docs"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "docs_count", Label: "Count", Stacked: true},
				mp.Metrics{Name: "docs_deleted", Label: "Deleted", Stacked: true},
			},
		},
		p.Prefix + ".indices.memory_size": mp.Graphs{
			Label: (p.LabelPrefix + " Indices Memory Size"),
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "fielddata_size", Label: "Fielddata", Stacked: true},
				mp.Metrics{Name: "filter_cache_size", Label: "Filter Cache", Stacked: true},
				mp.Metrics{Name: "segments_size", Label: "Lucene Segments", Stacked: true},
				mp.Metrics{Name: "segments_index_writer_size", Label: "Lucene Segments Index Writer", Stacked: true},
				mp.Metrics{Name: "segments_version_map_size", Label: "Lucene Segments Version Map", Stacked: true},
				mp.Metrics{Name: "segments_fixed_bit_set_size", Label: "Lucene Segments Fixed Bit Set", Stacked: true},
			},
		},
		p.Prefix + ".indices.evictions": mp.Graphs{
			Label: (p.LabelPrefix + " Indices Evictions"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "evictions_fielddata", Label: "Fielddata", Diff: true},
				mp.Metrics{Name: "evictions_filter_cache", Label: "Filter Cache", Diff: true},
			},
		},
		p.Prefix + ".jvm.heap": mp.Graphs{
			Label: (p.LabelPrefix + " JVM Heap Mem"),
			Unit:  "bytes",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "heap_used", Label: "Used"},
				mp.Metrics{Name: "heap_max", Label: "Max"},
			},
		},
		p.Prefix + ".thread_pool.threads": mp.Graphs{
			Label: (p.LabelPrefix + " Thread-Pool Threads"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "threads_generic", Label: "Generic", Stacked: true},
				mp.Metrics{Name: "threads_index", Label: "Index", Stacked: true},
				mp.Metrics{Name: "threads_snapshot_data", Label: "Snapshot Data", Stacked: true},
				mp.Metrics{Name: "threads_get", Label: "Get", Stacked: true},
				mp.Metrics{Name: "threads_bench", Label: "Bench", Stacked: true},
				mp.Metrics{Name: "threads_snapshot", Label: "Snapshot", Stacked: true},
				mp.Metrics{Name: "threads_merge", Label: "Merge", Stacked: true},
				mp.Metrics{Name: "threads_suggest", Label: "Suggest", Stacked: true},
				mp.Metrics{Name: "threads_bulk", Label: "Bulk", Stacked: true},
				mp.Metrics{Name: "threads_optimize", Label: "Optimize", Stacked: true},
				mp.Metrics{Name: "threads_warmer", Label: "Warmer", Stacked: true},
				mp.Metrics{Name: "threads_flush", Label: "Flush", Stacked: true},
				mp.Metrics{Name: "threads_search", Label: "Search", Stacked: true},
				mp.Metrics{Name: "threads_percolate", Label: "Percolate", Stacked: true},
				mp.Metrics{Name: "threads_refresh", Label: "Refresh", Stacked: true},
				mp.Metrics{Name: "threads_management", Label: "Management", Stacked: true},
				mp.Metrics{Name: "threads_fetch_shard_started", Label: "Fetch Shard Started", Stacked: true},
				mp.Metrics{Name: "threads_fetch_shard_store", Label: "Fetch Shard Store", Stacked: true},
				mp.Metrics{Name: "threads_listener", Label: "Listener", Stacked: true},
			},
		},
		p.Prefix + ".transport.count": mp.Graphs{
			Label: (p.LabelPrefix + " Transport Count"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "count_rx", Label: "TX", Diff: true},
				mp.Metrics{Name: "count_tx", Label: "RX", Diff: true},
			},
		},
		p.Prefix + ".process": mp.Graphs{
			Label: (p.LabelPrefix + " Process"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "open_file_descriptors", Label: "Open File Descriptors"},
			},
		},
	}

	return graphdef
}

func main() {
	optScheme := flag.String("scheme", "http", "Scheme")
	optHost := flag.String("host", "localhost", "Host")
	optPort := flag.String("port", "9200", "Port")
	optPrefix := flag.String("metric-key-prefix", "elasticsearch", "Metric key prefix")
	optLabelPrefix := flag.String("metric-label-prefix", "", "Metric Label prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var elasticsearch ElasticsearchPlugin
	elasticsearch.URI = fmt.Sprintf("%s://%s:%s", *optScheme, *optHost, *optPort)
	elasticsearch.Prefix = *optPrefix
	if *optLabelPrefix == "" {
		elasticsearch.LabelPrefix = strings.Title(*optPrefix)
	} else {
		elasticsearch.LabelPrefix = *optLabelPrefix
	}

	helper := mp.NewMackerelPlugin(elasticsearch)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-elasticsearch-%s-%s", *optHost, *optPort)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
