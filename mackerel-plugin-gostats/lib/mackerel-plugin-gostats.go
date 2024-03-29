package mpgostats

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

	stats_api "github.com/fukata/golang-stats-api-handler"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// GostatsPlugin mackerel plugin for go server
type GostatsPlugin struct {
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
func (m GostatsPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := cases.Title(language.Und, cases.NoLower).String(m.Prefix)
	return map[string]mp.Graphs{
		(m.Prefix + ".runtime"): {
			Label: (labelPrefix + " Runtime"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "goroutine_num", Label: "Goroutine Num"},
				{Name: "cgo_call_num", Label: "CGO Call Num", Diff: true},
			},
		},
		(m.Prefix + ".memory"): {
			Label: (labelPrefix + " Memory"),
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "memory_alloc", Label: "Alloc"},
				{Name: "memory_sys", Label: "Sys"},
				{Name: "memory_stack", Label: "Stack In Use"},
			},
		},
		(m.Prefix + ".operation"): {
			Label: (labelPrefix + " Operation"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "memory_lookups", Label: "Pointer Lookups", Diff: true},
				{Name: "memory_mallocs", Label: "Mallocs", Diff: true},
				{Name: "memory_frees", Label: "Frees", Diff: true},
			},
		},
		(m.Prefix + ".heap"): {
			Label: (labelPrefix + " Heap"),
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "heap_sys", Label: "Sys"},
				{Name: "heap_idle", Label: "Idle"},
				{Name: "heap_inuse", Label: "In Use"},
				{Name: "heap_released", Label: "Released", Diff: true},
			},
		},
		(m.Prefix + ".gc"): {
			Label: (labelPrefix + " GC"),
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "gc_num", Label: "GC Num", Diff: true},
				{Name: "gc_per_second", Label: "GC Per Second"},
				{Name: "gc_pause_per_second", Label: "GC Pause Per Second"},
			},
		},
	}
}

// FetchMetrics interface for mackerelplugin
func (m GostatsPlugin) FetchMetrics() (map[string]float64, error) {
	req, err := http.NewRequest(http.MethodGet, m.URI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "mackerel-plugin-gostats")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return m.parseStats(resp.Body)
}

func (m GostatsPlugin) parseStats(body io.Reader) (map[string]float64, error) {
	stat := make(map[string]float64)
	decoder := json.NewDecoder(body)

	s := stats_api.Stats{}
	err := decoder.Decode(&s)
	if err != nil {
		return nil, err
	}
	stat["goroutine_num"] = float64(s.GoroutineNum)
	stat["cgo_call_num"] = float64(s.CgoCallNum)
	stat["memory_sys"] = float64(s.MemorySys)
	stat["memory_alloc"] = float64(s.MemoryAlloc)
	stat["memory_stack"] = float64(s.StackInUse)
	stat["memory_lookups"] = float64(s.MemoryLookups)
	stat["memory_frees"] = float64(s.MemoryFrees)
	stat["memory_mallocs"] = float64(s.MemoryMallocs)
	stat["heap_sys"] = float64(s.HeapSys)
	stat["heap_idle"] = float64(s.HeapIdle)
	stat["heap_inuse"] = float64(s.HeapInuse)
	stat["heap_released"] = float64(s.HeapReleased)
	stat["gc_num"] = float64(s.GcNum)
	stat["gc_per_second"] = float64(s.GcPerSecond)
	stat["gc_pause_per_second"] = float64(s.GcPausePerSecond)

	return stat, nil
}

// Do the plugin
func Do() {
	optURI := flag.String("uri", "", "URI")
	optScheme := flag.String("scheme", "http", "Scheme")
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "8080", "Port")
	optPath := flag.String("path", "/api/stats", "Path")
	optPrefix := flag.String("metric-key-prefix", "gostats", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	gosrv := GostatsPlugin{
		Prefix: *optPrefix,
	}
	if *optURI != "" {
		gosrv.URI = *optURI
	} else {
		gosrv.URI = fmt.Sprintf("%s://%s:%s%s", *optScheme, *optHost, *optPort, *optPath)
	}

	helper := mp.NewMackerelPlugin(gosrv)
	helper.Tempfile = *optTempfile

	helper.Run()
}
