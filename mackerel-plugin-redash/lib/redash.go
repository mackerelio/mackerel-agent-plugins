package mpredash

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/golib/logging"
)

var logger = logging.GetLogger("metrics.plugin.redash")

// RedashPlugin mackerel plugin
type RedashPlugin struct {
	URI     string
	Prefix  string
	Timeout uint
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p RedashPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "redash"
	}
	return p.Prefix
}

// GraphDefinition interface for mackerelplugin
func (p RedashPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(p.Prefix)

	var graphdef = map[string]mp.Graphs{
		"task_queues_count": {
			Label: (labelPrefix + " Task Queues Count"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "wait", Label: "Wait", Diff: false},
				{Name: "done", Label: "Done", Diff: false},
				{Name: "in_progress", Label: "InProgress", Diff: false},
			},
		},
		"task_scheduled_count.#": {
			Label: (labelPrefix + " Task Scheduled Count"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "adhoc", Label: "Adhoc", Diff: false},
				{Name: "scheduled", Label: "Scheduled", Diff: false},
			},
		},
		"task_states_count": {
			Label: (labelPrefix + " Task States Count"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "waiting", Label: "Waiting", Diff: true},
				{Name: "finished", Label: "Finishing", Diff: true},
				{Name: "executing_query", Label: "Executing Query", Diff: true},
				{Name: "processing", Label: "Processing", Diff: true},
				{Name: "checking_alerts", Label: "Checking Alerts", Diff: true},
				{Name: "failed", Label: "Failed", Diff: true},
				{Name: "other", Label: "Other", Diff: true},
			},
		},
	}
	return graphdef
}

// FetchMetrics interface for mackerelplugin
func (p RedashPlugin) FetchMetrics() (map[string]interface{}, error) {
	stats, err := getUnsafeStats(p)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch redash metrics: %s", err)
	}

	metrics := make(map[string]interface{})
	metrics["wait"] = uint64(len(stats.WaitTasks))
	metrics["done"] = uint64(len(stats.DoneTasks))
	metrics["in_progress"] = uint64(len(stats.InProgressTasks))

	metrics["task_scheduled_count.wait.adhoc"] = filterCount(stats.WaitTasks, isAdhoc)
	metrics["task_scheduled_count.done.adhoc"] = filterCount(stats.DoneTasks, isAdhoc)
	metrics["task_scheduled_count.in_progress.adhoc"] = filterCount(stats.InProgressTasks, isAdhoc)
	metrics["task_scheduled_count.wait.scheduled"] = filterCount(stats.WaitTasks, isScheduled)
	metrics["task_scheduled_count.done.scheduled"] = filterCount(stats.DoneTasks, isScheduled)
	metrics["task_scheduled_count.in_progress.scheduled"] = filterCount(stats.InProgressTasks, isScheduled)

	for _, state := range UnsafeAllTaskStates {
		metrics[state] =
			filterCount(stats.WaitTasks, isState(state)) +
				filterCount(stats.DoneTasks, isState(state)) +
				filterCount(stats.InProgressTasks, isState(state))
	}
	return metrics, nil
}

// Do the plugin
func Do() {
	optURI := flag.String("uri", "http://localhost/api/admin/queries/tasks", "stats URI")
	apiKey := flag.String("api-key", os.Getenv("REDASH_API_KEY"), "API key")
	optPrefix := flag.String("metric-key-prefix", "redash", "Metric key prefix")
	optTimeout := flag.Uint("timeout", 5, "Timeout")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	// We recommend using -api-key option, but it's still possible specify ?api_key=XXX directly in -uri option.
	if *apiKey != "" {
		u, _ := url.Parse(*optURI)
		query := u.Query()
		query.Set("api_key", *apiKey)
		u.RawQuery = query.Encode()
		newURIString := u.String()
		optURI = &newURIString
	}

	p := RedashPlugin{
		URI:     *optURI,
		Prefix:  *optPrefix,
		Timeout: *optTimeout,
	}

	helper := mp.NewMackerelPlugin(p)
	helper.Tempfile = *optTempfile
	helper.Run()
}
