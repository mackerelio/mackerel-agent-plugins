package mpunicorn

import (
	"flag"
	"fmt"
	"os"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/golib/logging"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var logger = logging.GetLogger("metrics.plugin.unicorn")

// UnicornPlugin mackerel plugin for Unicorn
type UnicornPlugin struct {
	MasterPid  string
	WorkerPids []string
	Tempfile   string
	Prefix     string
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

// MetricKeyPrefix interface for PluginWithPrefix
func (u UnicornPlugin) MetricKeyPrefix() string {
	if u.Prefix == "" {
		u.Prefix = "unicorn"
	}
	return u.Prefix
}

// GraphDefinition interface for mackerelplugin
func (u UnicornPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := cases.Title(language.Und, cases.NoLower).String(u.MetricKeyPrefix())
	var graphdef = map[string]mp.Graphs{
		"memory": {
			Label: (labelPrefix + " Memory"),
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "memory_workers", Label: "Workers", Diff: false, Stacked: true},
				{Name: "memory_master", Label: "Master", Diff: false, Stacked: true},
				{Name: "memory_workeravg", Label: "Worker Average", Diff: false, Stacked: false},
			},
		},
		"workers": {
			Label: (labelPrefix + " Workers"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "busy_workers", Label: "Busy Workers", Diff: false, Stacked: true},
				{Name: "idle_workers", Label: "Idle Workers", Diff: false, Stacked: true},
			},
		},
	}

	return graphdef
}

// Do the plugin
func Do() {
	optPidFile := flag.String("pidfile", "", "Pid file name")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optPrefix := flag.String("metric-key-prefix", "unicorn", "Prefix")
	flag.Parse()
	var unicorn UnicornPlugin

	command = RealCommand{}
	pipedCommands = RealPipedCommands{}

	if *optPidFile == "" {
		logger.Errorf("Required unicorn pidfile.")
		os.Exit(1)
	} else {
		pid, err := os.ReadFile(*optPidFile)
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

	unicorn.Prefix = *optPrefix

	helper := mp.NewMackerelPlugin(unicorn)
	helper.Tempfile = *optTempfile

	helper.Run()
}
