package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"os"

	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string](mp.Graphs){
	"unicorn.memory": mp.Graphs{
		Label: "Unicorn Memory",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "memory_used", Label: "Memory Used", Diff: false, Stacked: true},
			mp.Metrics{Name: "memory_average", Label: "Memory Average", Diff: false, Stacked: true},
		},
	},
	"unicorn.workers": mp.Graphs{
		Label: "Unicorn Workers",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "worker_total", Label: "Worker Total", Diff: false, Stacked: true},
			mp.Metrics{Name: "worker_idles", Label: "Worker Idles", Diff: false, Stacked: true},
		},
	},
}

// UnicornPlugin mackerel plugin for Unicorn
type UnicornPlugin struct {
	MasterPid  string
	WorkerPids []string
}

// FetchMetrics interface for mackerelplugin
func (u UnicornPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	workers := len(u.WorkerPids)
	idles, err := idleWorkerCount(u.WorkerPids)
	if err != nil {
		return stat, err
	}
	used, err := usedMemory()
	if err != nil {
		return stat, err
	}
	average, err := averageMemory()
	if err != nil {
		return stat, err
	}

	stat["worker_total"] = fmt.Sprint(workers)
	stat["worker_idles"] = fmt.Sprint(idles)
	stat["memory_used"] = used
	stat["memory_average"] = average

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (n UnicornPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optPidFile := flag.String("pidfile", "", "pidfile path")
	flag.Parse()
	var unicorn UnicornPlugin

	command = RealCommand{}
	pipedCommands = RealPipedCommands{}

	if *optPidFile == "" {
		fmt.Errorf("Required unicorn pidfile.")
		os.Exit(1)
	} else {
		pid, err := ioutil.ReadFile(*optPidFile)
		if err != nil {
			fmt.Errorf("Failed to load unicorn pid file. %s", err)
			os.Exit(1)
		}
		unicorn.MasterPid = strings.Replace(string(pid), "\n", "", 1)
	}

	workerPids, err := fetchUnicornWorkerPids(unicorn.MasterPid)
	if err != nil {
		fmt.Errorf("Failed to fetch unicorn worker pids. %s", err)
		os.Exit(1)
	}
	unicorn.WorkerPids = workerPids

	helper := mp.NewMackerelPlugin(unicorn)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
