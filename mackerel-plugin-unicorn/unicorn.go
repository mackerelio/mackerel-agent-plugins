package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"os"
	"os/exec"

	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mattn/go-pipeline"
)

var graphdef = map[string](mp.Graphs){
	"unicorn.memory": mp.Graphs{
		Label: "Unicorn Memory",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "memory_used", Label: "Memory Used", Diff: false},
		},
	},
	"unicorn.workers": mp.Graphs{
		Label: "Unicorn Workers",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "worker_total", Label: "Worker Total", Diff: false},
			mp.Metrics{Name: "worker_idles", Label: "Worker Idles", Diff: false},
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

	stat["worker_total"] = fmt.Sprint(workers)
	stat["worker_idles"] = fmt.Sprint(idles)
	stat["memory_used"] = used

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (n UnicornPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func usedMemory() (string, error) {
	out, err := pipeline.Output(
		[]string{"ps", "auxw"},
		[]string{"grep", "[u]nicorn"},
		[]string{"awk", "{m+=$6*1024} END{print m;}"},
	)
	if err != nil {
		return "", fmt.Errorf("Cannot get unicorn used memory")
	}
	return strings.Trim(string(out), "\n"), nil
}

func idleWorkerCount(pids []string) (int, error) {
	var beforeCpu []string
	var afterCpu []string
	idles := 0

	for _, pid := range pids {
		cputime, err := cpuTime(pid)
		if err != nil {
			return idles, err
		}
		beforeCpu = append(beforeCpu, cputime)
	}
	time.Sleep(1 * time.Second)
	for _, pid := range pids {
		cputime, err := cpuTime(pid)
		if err != nil {
			return idles, err
		}
		afterCpu = append(afterCpu, cputime)
	}
	for i, _ := range pids {
		b, _ := strconv.Atoi(beforeCpu[i])
		a, _ := strconv.Atoi(afterCpu[i])
		if (a - b) == 0 {
			idles++
		}
	}

	return idles, nil
}

func cpuTime(pid string) (string, error) {
	out, err := pipeline.Output(
		[]string{"cat", fmt.Sprintf("/proc/%s/stat", pid)},
		[]string{"awk", "{print $14+$15}"},
	)
	if err != nil {
		return "", fmt.Errorf("Failed to cat /proc/%s/stat: %s", pid, err)
	}
	return string(out), nil
}

func fetchUnicornWorkerPids(m string) ([]string, error) {
	var workerPids []string

	out, err := exec.Command("ps", "w", "--ppid", m).Output()
	if err != nil {
		return workerPids, fmt.Errorf("Failed to ps of command: %s", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		words := strings.SplitN(line, " ", 5)
		if len(words) < 5 {
			continue
		}
		pid := words[0]
		cmd := words[4]
		if strings.Contains(cmd, "worker") {
			workerPids = append(workerPids, pid)
		}
	}

	if len(workerPids) > 0 {
		return workerPids, nil
	}

	return workerPids, fmt.Errorf("Cannot get unicorn worker pids")
}

func main() {
	optPidFile := flag.String("pidfile", "", "pidfile path")
	flag.Parse()
	var unicorn UnicornPlugin

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
