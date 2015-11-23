package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"os"
	"os/exec"

	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mattn/go-pipeline"
)

var graphdef = map[string](mp.Graphs){
	"unicorn.memory": mp.Graphs{
		Label: "Unicorn Memory",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "used", Label: "Memory Used", Diff: false},
		},
	},
	"unicorn.workers": mp.Graphs{
		Label: "Unicorn Workers",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "total", Label: "Worker Total", Diff: false},
			mp.Metrics{Name: "idles", Label: "Worker Idles", Diff: false},
		},
	},
}

// UnicornPlugin mackerel plugin for Unicorn
type UnicornPlugin struct {
	MasterPid  string
	WorkerPids []string
}

// FetchMetrics interface for mackerelplugin
func (n UnicornPlugin) FetchMetrics() (map[string]interface{}, error) {

	return n.parseStats(resp.Body)
}

// GraphDefinition interface for mackerelplugin
func (n UnicornPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func totalMemory() ([]string, error) {
	var m []string
	return m, nil
}

func memoryUsage() ([]string, error) {
	var m []string
	out, err := exec.Command("ps", "auxw").Output()
	if err != nil {
		return m, fmt.Errorf("Failed to ps of command: %s", err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		words := strings.Split(line, " ")
		if len(words) != 2 {
			continue
		}
		masterPid, name := words[0], words[1]
		if name == "unicorn" {
			m = append(m, masterPid)
		}
	}
	return m, fmt.Errorf("Cannot get unicorn master pid")
}

func workerCount() (int, error) {
	return 0, nil
}

func idleWorkerCount() (int, error) {
	return 0, nil
}

func cpuTime() (int, error) {
	return 0, nil
}

func fetchUnicornMasterPid() (string, error) {
	out, err := pipeline.Output(
		[]string{"ps", "ax"},
		[]string{"grep", "unicorn"},
		[]string{"grep", "master"},
	)
	if err != nil {
		return "", fmt.Errorf("Failed to ps and grep: %s", err)
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 1 {
		return "", fmt.Errorf("At least two unicorn master processes.")
	}
	return lines[0], nil
}

func fetchUnicornWorkerPids(m string) ([]string, error) {
	var workerPids []string

	out, err := exec.Command("ps", "w", "--ppid", m).Output()
	if err != nil {
		return workerPids, fmt.Errorf("Failed to ps of command: %s", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		words := strings.SplitN(line, " ", 5)
		pid, cmd := words[0], words[4]
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
		masterPid, err := fetchUnicornMasterPid()
		if err != nil {
			fmt.Errorf("Failed to fetch unicorn pid. %s", err)
			os.Exit(1)
		}
		unicorn.MasterPid = masterPid
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
