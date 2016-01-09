package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metrics.plugin.elasticsearch")

var graphdef = map[string](mp.Graphs){
	"elasticsearch.http": mp.Graphs{
		Label: "Elasticsearch HTTP",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "http_opened", Label: "Opened", Diff: true},
		},
	},
	"elasticsearch.indices": mp.Graphs{
		Label: "Elasticsearch Indices",
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
	"elasticsearch.indices.docs": mp.Graphs{
		Label: "Elasticsearch Indices Docs",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "docs_count", Label: "Count", Stacked: true},
			mp.Metrics{Name: "docs_deleted", Label: "Deleted", Stacked: true},
		},
	},
	"elasticsearch.indices.memory_size": mp.Graphs{
		Label: "Elasticsearch Indices Memory Size",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "fielddata_size", Label: "Fielddata", Stacked: true},
			mp.Metrics{Name: "segments_size", Label: "Lucene Segments", Stacked: true},
			mp.Metrics{Name: "segments_index_writer_size", Label: "Lucene Segments Index Writer", Stacked: true},
			mp.Metrics{Name: "segments_version_map_size", Label: "Lucene Segments Version Map", Stacked: true},
			mp.Metrics{Name: "segments_fixed_bit_set_size", Label: "Lucene Segments Fixed Bit Set", Stacked: true},
		},
	},
	"elasticsearch.indices.evictions": mp.Graphs{
		Label: "Elasticsearch Indices Evictions",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "evictions_fielddata", Label: "Fielddata", Diff: true},
		},
	},
	"elasticsearch.jvm.heap": mp.Graphs{
		Label: "Elasticsearch JVM Heap Mem",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "heap_used", Label: "Used"},
			mp.Metrics{Name: "heap_max", Label: "Max"},
		},
	},
	"elasticsearch.thread_pool.threads": mp.Graphs{
		Label: "Elasticsearch Thread-Pool Threads",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "threads_bulk", Label: "Bulk", Stacked: true},
			mp.Metrics{Name: "threads_fetch_shard_started", Label: "Fetch Shard Started", Stacked: true},
			mp.Metrics{Name: "threads_fetch_shard_store", Label: "Fetch Shard Store", Stacked: true},
			mp.Metrics{Name: "threads_flush", Label: "Flush", Stacked: true},
			mp.Metrics{Name: "threads_generic", Label: "Generic", Stacked: true},
			mp.Metrics{Name: "threads_get", Label: "Get", Stacked: true},
			mp.Metrics{Name: "threads_index", Label: "Index", Stacked: true},
			mp.Metrics{Name: "threads_listener", Label: "Listener", Stacked: true},
			mp.Metrics{Name: "threads_management", Label: "Management", Stacked: true},
			mp.Metrics{Name: "threads_optimize", Label: "Optimize", Stacked: true},
			mp.Metrics{Name: "threads_percolate", Label: "Percolate", Stacked: true},
			mp.Metrics{Name: "threads_refresh", Label: "Refresh", Stacked: true},
			mp.Metrics{Name: "threads_search", Label: "Search", Stacked: true},
			mp.Metrics{Name: "threads_snapshot", Label: "Snapshot", Stacked: true},
			mp.Metrics{Name: "threads_suggest", Label: "Suggest", Stacked: true},
			mp.Metrics{Name: "threads_warmer", Label: "Warmer", Stacked: true},
		},
	},
	"elasticsearch.transport.count": mp.Graphs{
		Label: "Elasticsearch Transport Count",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "count_rx", Label: "TX", Diff: true},
			mp.Metrics{Name: "count_tx", Label: "RX", Diff: true},
		},
	},
}

var metricPlace = map[string][]string{
	"http_opened":                 []string{"http", "total_opened"},
	"total_indexing_index":        []string{"indices", "indexing", "index_total"},
	"total_indexing_delete":       []string{"indices", "indexing", "delete_total"},
	"total_get":                   []string{"indices", "get", "total"},
	"total_search_query":          []string{"indices", "search", "query_total"},
	"total_search_fetch":          []string{"indices", "search", "fetch_total"},
	"total_merges":                []string{"indices", "merges", "total"},
	"total_refresh":               []string{"indices", "refresh", "total"},
	"total_flush":                 []string{"indices", "flush", "total"},
	"total_warmer":                []string{"indices", "warmer", "total"},
	"total_percolate":             []string{"indices", "percolate", "total"},
	"total_suggest":               []string{"indices", "suggest", "total"},
	"docs_count":                  []string{"indices", "docs", "count"},
	"docs_deleted":                []string{"indices", "docs", "deleted"},
	"fielddata_size":              []string{"indices", "fielddata", "memory_size_in_bytes"},
	"segments_size":               []string{"indices", "segments", "memory_in_bytes"},
	"segments_index_writer_size":  []string{"indices", "segments", "index_writer_memory_in_bytes"},
	"segments_version_map_size":   []string{"indices", "segments", "version_map_memory_in_bytes"},
	"segments_fixed_bit_set_size": []string{"indices", "segments", "fixed_bit_set_memory_in_bytes"},
	"evictions_fielddata":         []string{"indices", "fielddata", "evictions"},
	"heap_used":                   []string{"jvm", "mem", "heap_used_in_bytes"},
	"heap_max":                    []string{"jvm", "mem", "heap_max_in_bytes"},
	"threads_bulk":                []string{"thread_pool", "bulk", "threads"},
	"threads_fetch_shard_started": []string{"thread_pool", "fetch_shard_started", "threads"},
	"threads_fetch_shard_store":   []string{"thread_pool", "fetch_shard_store", "threads"},
	"threads_flush":               []string{"thread_pool", "flush", "threads"},
	"threads_generic":             []string{"thread_pool", "generic", "threads"},
	"threads_get":                 []string{"thread_pool", "get", "threads"},
	"threads_index":               []string{"thread_pool", "index", "threads"},
	"threads_listener":            []string{"thread_pool", "listener", "threads"},
	"threads_management":          []string{"thread_pool", "management", "threads"},
	"threads_optimize":            []string{"thread_pool", "optimize", "threads"},
	"threads_percolate":           []string{"thread_pool", "percolate", "threads"},
	"threads_refresh":             []string{"thread_pool", "refresh", "threads"},
	"threads_search":              []string{"thread_pool", "search", "threads"},
	"threads_snapshot":            []string{"thread_pool", "snapshot", "threads"},
	"threads_suggest":             []string{"thread_pool", "suggest", "threads"},
	"threads_warmer":              []string{"thread_pool", "warmer", "threads"},
	"count_rx":                    []string{"transport", "rx_count"},
	"count_tx":                    []string{"transport", "tx_count"},
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
	URI string
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
	return graphdef
}

func main() {
	optHost := flag.String("host", "localhost", "Host")
	optPort := flag.String("port", "9200", "Port")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var elasticsearch ElasticsearchPlugin
	elasticsearch.URI = fmt.Sprintf("http://%s:%s", *optHost, *optPort)

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
