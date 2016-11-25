package mpjmxjolokia

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metrics.plugin.jmx-jolokia")

// JmxJolokiaPlugin mackerel plugin for Jolokia
type JmxJolokiaPlugin struct {
	Target   string
	Tempfile string
}

// JmxJolokiaResponse response for Jolokia
type JmxJolokiaResponse struct {
	Status    uint32
	Timestamp uint32
	Request   map[string]interface{}
	Value     map[string]interface{}
	Error     string
}

var graphdef = map[string]mp.Graphs{
	"jmx.jolokia.memory.heap_memory_usage": {
		Label: "Jmx HeapMemoryUsage",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "HeapMemoryInit", Label: "init", Diff: false, Type: "uint64"},
			{Name: "HeapMemoryCommitted", Label: "committed", Diff: false, Type: "uint64"},
			{Name: "HeapMemoryMax", Label: "max", Diff: false, Type: "uint64"},
			{Name: "HeapMemoryUsed", Label: "used", Diff: false, Type: "uint64"},
		},
	},
	"jmx.jolokia.memory.non_heap_memory_usage": {
		Label: "Jmx NonHeapMemoryUsage",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "NonHeapMemoryInit", Label: "init", Diff: false, Type: "uint64"},
			{Name: "NonHeapMemoryCommitted", Label: "committed", Diff: false, Type: "uint64"},
			{Name: "NonHeapMemoryMax", Label: "max", Diff: false, Type: "uint64"},
			{Name: "NonHeapMemoryUsed", Label: "used", Diff: false, Type: "uint64"},
		},
	},
	"jmx.jolokia.class_load": {
		Label: "Jmx ClassLoading",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "LoadedClassCount", Label: "loaded", Diff: false, Type: "uint64"},
			{Name: "UnloadedClassCount", Label: "unloaded", Diff: false, Type: "uint64"},
			{Name: "TotalLoadedClassCount", Label: "total", Diff: false, Type: "uint64"},
		},
	},
	"jmx.jolokia.thread": {
		Label: "Jmx Threading",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "ThreadCount", Label: "thread", Diff: false, Type: "uint64"},
		},
	},
	"jmx.jolokia.ops.cpu_load": {
		Label: "Jmx CpuLoad",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "ProcessCpuLoad", Label: "process", Diff: false, Type: "float64", Scale: 100},
			{Name: "SystemCpuLoad", Label: "system", Diff: false, Type: "float64", Scale: 100},
		},
	},
}

// FetchMetrics interface for mackerelplugin
func (j JmxJolokiaPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})
	if err := j.fetchMemory(stat); err != nil {
		logger.Warningf(err.Error())
	}

	if err := j.fetchClassLoad(stat); err != nil {
		logger.Warningf(err.Error())
	}

	if err := j.fetchThread(stat); err != nil {
		logger.Warningf(err.Error())
	}

	if err := j.fetchOperatingSystem(stat); err != nil {
		logger.Warningf(err.Error())
	}

	return stat, nil
}

func (j JmxJolokiaPlugin) fetchMemory(stat map[string]interface{}) error {
	resp, err := j.executeGetRequest("java.lang:type=Memory")
	if err != nil {
		return err
	}
	heap := resp.Value["HeapMemoryUsage"].(map[string]interface{})
	stat["HeapMemoryInit"] = heap["init"]
	stat["HeapMemoryCommitted"] = heap["committed"]
	stat["HeapMemoryMax"] = heap["max"]
	stat["HeapMemoryUsed"] = heap["used"]

	nonHeap := resp.Value["NonHeapMemoryUsage"].(map[string]interface{})
	stat["NonHeapMemoryInit"] = nonHeap["init"]
	stat["NonHeapMemoryCommitted"] = nonHeap["committed"]
	stat["NonHeapMemoryMax"] = nonHeap["max"]
	stat["NonHeapMemoryUsed"] = nonHeap["used"]

	return nil
}

func (j JmxJolokiaPlugin) fetchClassLoad(stat map[string]interface{}) error {
	resp, err := j.executeGetRequest("java.lang:type=ClassLoading")
	if err != nil {
		return err
	}
	stat["LoadedClassCount"] = resp.Value["LoadedClassCount"]
	stat["UnloadedClassCount"] = resp.Value["UnloadedClassCount"]
	stat["TotalLoadedClassCount"] = resp.Value["TotalLoadedClassCount"]

	return nil
}

func (j JmxJolokiaPlugin) fetchThread(stat map[string]interface{}) error {
	resp, err := j.executeGetRequest("java.lang:type=Threading")
	if err != nil {
		return err
	}
	stat["ThreadCount"] = resp.Value["ThreadCount"]

	return nil
}

func (j JmxJolokiaPlugin) fetchOperatingSystem(stat map[string]interface{}) error {
	resp, err := j.executeGetRequest("java.lang:type=OperatingSystem")
	if err != nil {
		return err
	}
	stat["ProcessCpuLoad"] = resp.Value["ProcessCpuLoad"]
	stat["SystemCpuLoad"] = resp.Value["SystemCpuLoad"]

	return nil
}

func (j JmxJolokiaPlugin) executeGetRequest(mbean string) (*JmxJolokiaResponse, error) {
	resp, err := http.Get(j.Target + mbean)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var respJ JmxJolokiaResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&respJ); err != nil {
		return nil, err
	}
	return &respJ, nil
}

// GraphDefinition interface for mackerelplugin
func (j JmxJolokiaPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "8778", "Port")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var jmxJolokia JmxJolokiaPlugin
	jmxJolokia.Target = fmt.Sprintf("http://%s:%s/jolokia/read/", *optHost, *optPort)

	helper := mp.NewMackerelPlugin(jmxJolokia)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-jmx-jolokia-%s-%s", *optHost, *optPort)
	}
	helper.Run()
}
