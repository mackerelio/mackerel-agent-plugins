package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mattn/go-pipeline"
)

var logger = logging.GetLogger("metrics.plugin.proc-fd")

type ProcfdPlugin struct {
	Process           string
	NormalizedProcess string
	MetricName        string
}

func (p ProcfdPlugin) FetchMetrics() (map[string]interface{}, error) {
	fds, err := p.getNumOpenFileDesc()
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

func (p ProcfdPlugin) getNumOpenFileDesc() (map[string]uint64, error) {
	fds := make(map[string]uint64)

	// Fetch all pids which contain specified process name
	out, err := pipeline.Output(
		[]string{"ps", "aux"},
		[]string{"grep", p.Process},
		[]string{"grep", "-v", "grep"},
		[]string{"grep", "-v", "mackerel-plugin-proc-fd"},
		[]string{"awk", "{print $2}"},
	)
	if err != nil {
		// No matching with p.Process invokes this case
		logger.Errorf("No matching with process")
		return nil, errors.New("No matching with process")
	}

	// List the number of all open files beloging to each pid
	for _, pid := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		out, err = pipeline.Output(
			[]string{"ls", "-l", fmt.Sprintf("/proc/%s/fd", pid)},
			[]string{"grep", "-v", "total"},
			[]string{"wc", "-l"},
		)
		if err != nil {
			// The process with pid terminates"
			logger.Errorf("The process terminates")
			return nil, errors.New("The process terminates")
		}

		num, err := strconv.ParseUint(strings.TrimSpace(string(out)), 10, 32)
		if err != nil {
			return nil, err
		}
		fds[pid] = num
	}

	return fds, nil
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
