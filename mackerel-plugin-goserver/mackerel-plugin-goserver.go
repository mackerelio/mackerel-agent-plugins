package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fukata/golang-stats-api-handler"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

type GoServerPlugin struct {
	URI    string
	Prefix string
}

/*
{
  "time": 1449124022112358000,
  "go_version": "go1.5.1",
  "go_os": "darwin",
  "go_arch": "amd64",
  "cpu_num": 4,
  "goroutine_num": 6,
  "gomaxprocs": 4,
  "cgo_call_num": 5,
  "memory_alloc": 213360,
  "memory_total_alloc": 213360,
  "memory_sys": 3377400,
  "memory_lookups": 15,
  "memory_mallocs": 1137,
  "memory_frees": 0,
  "memory_stack": 393216,
  "heap_alloc": 213360,
  "heap_sys": 655360,
  "heap_idle": 65536,
  "heap_inuse": 589824,
  "heap_released": 0,
  "heap_objects": 1137,
  "gc_next": 4194304,
  "gc_last": 0,
  "gc_num": 0,
  "gc_per_second": 0,
  "gc_pause_per_second": 0,
  "gc_pause": []
}
*/

// GraphDefinition interface for mackerelplugin
func (m GoServerPlugin) GraphDefinition() map[string](mp.Graphs) {
	labelPrefix := strings.Title(m.Prefix)
	return map[string](mp.Graphs){
		(m.Prefix + ".runtime"): mp.Graphs{
			Label: (labelPrefix + " Runtime"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "goroutine_num", Label: "Gorotine Num"},
				mp.Metrics{Name: "cgo_call_num", Label: "CGO Call Num", Diff: true},
			},
		},
		(m.Prefix + ".memory"): mp.Graphs{
			Label: (labelPrefix + " Memory"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "memory_alloc", Label: "Alloc"},
				mp.Metrics{Name: "memory_total_alloc", Label: "Total Alloc", Diff: true},
				mp.Metrics{Name: "memory_sys", Label: "Sys"},
				mp.Metrics{Name: "memory_lookups", Label: "Lookups"},
				mp.Metrics{Name: "memory_mallocs", Label: "Mallocs"},
				mp.Metrics{Name: "memory_frees", Label: "Frees"},
				mp.Metrics{Name: "memory_stack", Label: "Stack In Use"},
			},
		},
		(m.Prefix + ".heap"): mp.Graphs{
			Label: (labelPrefix + " Heap"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "heap_alloc", Label: "Alloc"},
				mp.Metrics{Name: "heap_sys", Label: "Sys"},
				mp.Metrics{Name: "heap_idle", Label: "Idle"},
				mp.Metrics{Name: "heap_inuse", Label: "In Use"},
				mp.Metrics{Name: "heap_released", Label: "Released"},
				mp.Metrics{Name: "heap_objects", Label: "Objects"},
			},
		},
		(m.Prefix + ".gc"): mp.Graphs{
			Label: (labelPrefix + " GC"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "gc_num", Label: "Num"},
				mp.Metrics{Name: "gc_per_second", Label: "GC Per Second"},
				mp.Metrics{Name: "gc_pause_per_second", Label: "GC Pause Per Second"},
			},
		},
	}
}

// FetchMetrics interface for mackerelplugin
func (m GoServerPlugin) FetchMetrics() (map[string]interface{}, error) {
	resp, err := http.Get(m.URI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return m.parseStats(resp.Body)
}

func (m GoServerPlugin) parseStats(body io.Reader) (map[string]interface{}, error) {
	stat := make(map[string]interface{})
	decoder := json.NewDecoder(body)

	s := stats_api.Stats{}
	err := decoder.Decode(&s)
	if err != nil {
		return nil, err
	}
	stat["goroutine_num"] = uint64(s.GoroutineNum)
	stat["cgo_call_num"] = uint64(s.CgoCallNum)
	stat["memory_alloc"] = s.MemoryAlloc
	stat["memory_total_alloc"] = s.MemoryTotalAlloc
	stat["memory_sys"] = s.MemorySys
	stat["memory_lookups"] = s.MemoryLookups
	stat["memory_frees"] = s.MemoryFrees
	stat["memory_stack"] = s.StackInUse
	stat["heap_alloc"] = s.HeapAlloc
	stat["heap_sys"] = s.HeapSys
	stat["heap_idle"] = s.HeapIdle
	stat["heap_inuse"] = s.HeapInuse
	stat["heap_released"] = s.HeapReleased
	stat["heap_objects"] = s.HeapObjects
	stat["gc_num"] = s.GcNum
	stat["gc_per_second"] = s.GcPerSecond
	stat["gc_pause_per_second"] = s.GcPausePerSecond

	return stat, nil
}

func main() {
	optURI := flag.String("uri", "", "URI")
	optScheme := flag.String("scheme", "http", "Scheme")
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "8080", "Port")
	optPath := flag.String("path", "/api/stats", "Path")
	optPrefix := flag.String("metric-key-prefix", "goserver", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	gosrv := GoServerPlugin{
		Prefix: *optPrefix,
	}
	if *optURI != "" {
		gosrv.URI = *optURI
	} else {
		gosrv.URI = fmt.Sprintf("%s://%s:%s%s", *optScheme, *optHost, *optPort, *optPath)
	}

	helper := mp.NewMackerelPlugin(gosrv)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-gosrv")
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
