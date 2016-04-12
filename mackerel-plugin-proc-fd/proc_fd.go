package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metrics.plugin.proc-fd")

// ProcfdPlugin for fetching metrics
type ProcfdPlugin struct {
	Process           string
	NormalizedProcess string
	MetricName        string
}

// FetchMetrics fetch the metrics
func (p ProcfdPlugin) FetchMetrics() (map[string]interface{}, error) {
	fds, err := openFd.getNumOpenFileDesc()
	if err != nil {
		return nil, err
	}

	stat := make(map[string]interface{})

	// Compute maximum open file descriptor
	var maxFD uint64
	for _, fd := range fds {
		if fd > maxFD {
			maxFD = fd
		}
	}
	stat["max_fd"] = maxFD

	return stat, nil
}

// GraphDefinition Graph definition
func (p ProcfdPlugin) GraphDefinition() map[string](mp.Graphs) {
	return map[string](mp.Graphs){
		fmt.Sprintf("proc-fd.%s", p.NormalizedProcess): mp.Graphs{
			Label: fmt.Sprintf("Opening fd by %s", p.NormalizedProcess),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "max_fd", Label: "Maximum opening fd", Diff: false, Type: "uint64"},
			},
		},
	}
}

func normalizeForMetricName(process string) string {
	// Mackerel accepts following characters in custom metric names
	// [-a-zA-Z0-9_.]
	re := regexp.MustCompile("[^-a-zA-Z0-9_.]")
	return re.ReplaceAllString(process, "_")
}

func main() {
	optProcess := flag.String("process", "", "Process name")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	if *optProcess == "" {
		logger.Warningf("Process name is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var fd ProcfdPlugin
	fd.Process = *optProcess
	openFd = RealOpenFd{fd.Process}
	fd.NormalizedProcess = normalizeForMetricName(*optProcess)

	helper := mp.NewMackerelPlugin(fd)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-proc-fd")
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
