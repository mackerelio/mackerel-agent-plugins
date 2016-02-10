package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

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

var graphdef = map[string](mp.Graphs){
	"jolokia.memory.heap_memory_usage": mp.Graphs{
		Label: "HeapMemoryUsage",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "init", Label: "init", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "committed", Label: "committed", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "max", Label: "max", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "used", Label: "used", Diff: false, Type: "uint64"},
		},
	},
	"jolokia.memory.non_heap_memory_usage": mp.Graphs{
		Label: "NonHeapMemoryUsage",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "nonInit", Label: "init", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "nonCommitted", Label: "committed", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "nonMax", Label: "max", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "nonUsed", Label: "used", Diff: false, Type: "uint64"},
		},
	},
	"jolokia.class_load": mp.Graphs{
		Label: "ClassLoading",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "LoadedClassCount", Label: "loaded", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "UnloadedClassCount", Label: "unloaded", Diff: false, Type: "uint64"},
			mp.Metrics{Name: "TotalLoadedClassCount", Label: "total", Diff: false, Type: "uint64"},
		},
	},
	"jolokia.thread": mp.Graphs{
		Label: "Threading",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "ThreadCount", Label: "thread", Diff: false, Type: "uint64"},
		},
	},
	"jolokia.ops.cpu_load": mp.Graphs{
		Label: "CpuLoad",
		Unit:  "percentage",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "ProcessCpuLoad", Label: "process %", Diff: false, Type: "float64", Scale: 100},
			mp.Metrics{Name: "SystemCpuLoad", Label: "system %", Diff: false, Type: "float64", Scale: 100},
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
	stat["init"] = heap["init"]
	stat["committed"] = heap["committed"]
	stat["max"] = heap["max"]
	stat["used"] = heap["used"]

	nonHeap := resp.Value["NonHeapMemoryUsage"].(map[string]interface{})
	stat["nonInit"] = nonHeap["init"]
	stat["nonCommitted"] = nonHeap["committed"]
	stat["nonMax"] = nonHeap["max"]
	stat["nonUsed"] = nonHeap["used"]

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
func (j JmxJolokiaPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
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

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
