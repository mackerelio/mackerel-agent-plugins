package mpunicorn

import (
	"flag"
	"fmt"
	"io/ioutil"

	"os"

	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metrics.plugin.unicorn")

var graphdef = map[string]mp.Graphs{
	"unicorn.memory": {
		Label: "Unicorn Memory",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "memory_workers", Label: "Workers", Diff: false, Stacked: true},
			{Name: "memory_master", Label: "Master", Diff: false, Stacked: true},
			{Name: "memory_workeravg", Label: "Worker Average", Diff: false, Stacked: false},
		},
	},
	"unicorn.workers": {
		Label: "Unicorn Workers",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "busy_workers", Label: "Busy Workers", Diff: false, Stacked: true},
			{Name: "idle_workers", Label: "Idle Workers", Diff: false, Stacked: true},
		},
	},
}

// UnicornPlugin mackerel plugin for Unicorn
type UnicornPlugin struct {
	MasterPid  string
	WorkerPids []string
	Tempfile   string
}

// FetchMetrics interface for mackerelplugin
func (u UnicornPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	workers := len(u.WorkerPids)
	idles, err := idleWorkerCount(u.WorkerPids)
	if err != nil {
		return stat, err
	}
	stat["idle_workers"] = fmt.Sprint(idles)
	stat["busy_workers"] = fmt.Sprint(workers - idles)

	workersM, err := workersMemory()
	if err != nil {
		return stat, err
	}
	stat["memory_workers"] = workersM

	masterM, err := masterMemory()
	if err != nil {
		return stat, err
	}
	stat["memory_master"] = masterM

	averageM, err := workersMemoryAvg()
	if err != nil {
		return stat, err
	}
	stat["memory_workeravg"] = averageM

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (u UnicornPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	optPidFile := flag.String("pidfile", "", "Pid file name")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()
	var unicorn UnicornPlugin

	command = RealCommand{}
	pipedCommands = RealPipedCommands{}

	if *optPidFile == "" {
		logger.Errorf("Required unicorn pidfile.")
		os.Exit(1)
	} else {
		pid, err := ioutil.ReadFile(*optPidFile)
		if err != nil {
			logger.Errorf("Failed to load unicorn pid file. %s", err)
			os.Exit(1)
		}
		unicorn.MasterPid = strings.Replace(string(pid), "\n", "", 1)
	}

	workerPids, err := fetchUnicornWorkerPids(unicorn.MasterPid)
	if err != nil {
		logger.Errorf("Failed to fetch unicorn worker pids. %s", err)
		os.Exit(1)
	}
	unicorn.WorkerPids = workerPids

	helper := mp.NewMackerelPlugin(unicorn)
	helper.Tempfile = *optTempfile

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
