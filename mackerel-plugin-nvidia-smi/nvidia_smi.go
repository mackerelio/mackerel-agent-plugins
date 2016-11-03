package main

import (
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"os/exec"
	"strconv"
	"strings"
)

var queryOptions = []string{
	"utilization.gpu",
	"utilization.memory",
	"temperature.gpu",
	"fan.speed",
	"memory.total",
	"memory.used",
	"memory.free",
}

var formatOptions = []string{
	"noheader",
	"nounits",
	"csv",
}

var nvidiaSmiOptions = []string{
	fmt.Sprintf("--format=%s", strings.Join(formatOptions, ",")),
	fmt.Sprintf("--query-gpu=%s", strings.Join(queryOptions, ",")),
}

var metricsKeyFormats = []string{
	"gpu.util.gpu%d",
	"memory.util.gpu%d",
	"temperature.gpu%d",
	"fanspeed.gpu%d",
	"memory.usage.gpu%d.total",
	"memory.usage.gpu%d.used",
	"memory.usage.gpu%d.free",
}

func (n NVidiaSMIPlugin) getMetricKey(index int, gpuIndex int) string {
	return fmt.Sprintf(metricsKeyFormats[index], gpuIndex)
}

// NVidiaSMIPlugin mackerel plugin for nvidia-smi
type NVidiaSMIPlugin struct {
	Prefix string
}

// GraphDefinition interface for mackerelplugin
func (n NVidiaSMIPlugin) GraphDefinition() map[string](mp.Graphs) {
	var graphdef = map[string](mp.Graphs){
		"gpu.util": mp.Graphs{
			Label: "GPU Utilization",
			Unit:  "percentage",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "#", Label: "util"},
			},
		},
		"memory.util": mp.Graphs{
			Label: "GPU Memory Utilization",
			Unit:  "percentage",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "#", Label: "util"},
			},
		},
		"temperature": mp.Graphs{
			Label: "GPU Temperature",
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "#", Label: "temp"},
			},
		},
		"fanspeed": mp.Graphs{
			Label: "GPU Fan Speed",
			Unit:  "percentage",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "#", Label: "fan speed"},
			},
		},
		"memory.usage.#": mp.Graphs{
			Label: "GPU Memory Usage",
			Unit:  "bytes",

			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "total", Label: "total", Scale: 1024 * 1024},
				mp.Metrics{Name: "used", Label: "used", Scale: 1024 * 1024, Stacked: true},
				mp.Metrics{Name: "free", Label: "free", Scale: 1024 * 1024, Stacked: true},
			},
		},
	}
	return graphdef
}

// FetchMetrics interface for mackerelplugin
func (n NVidiaSMIPlugin) FetchMetrics() (map[string]interface{}, error) {
	ret, err := exec.Command("nvidia-smi", nvidiaSmiOptions...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, ret)
	}
	return n.parseStats(string(ret))
}

// MetricKeyPrefix interface for mackerelplugin
func (n NVidiaSMIPlugin) MetricKeyPrefix() string {
	return n.Prefix
}

func (n NVidiaSMIPlugin) parseStats(ret string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	for id, line := range strings.Split(ret, "\n") {
		err := n.parseLine(id, line, &stats)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", err, ret)
		}
	}
	return stats, nil
}

func (n NVidiaSMIPlugin) parseLine(id int, line string, stats *map[string]interface{}) error {
	if strings.TrimSpace(line) == "" {
		return nil
	}

	for i, value := range strings.Split(line, ",") {
		value, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
		if err != nil {
			return err
		}
		(*stats)[n.getMetricKey(i, id)] = value
	}
	return nil
}

func main() {
	optPrefix := flag.String("prefix", "nvidia.gpu", "Metric key prefix")
	flag.Parse()
	var plugin NVidiaSMIPlugin
	plugin.Prefix = *optPrefix
	helper := mp.NewMackerelPlugin(plugin)
	helper.Run()
}
