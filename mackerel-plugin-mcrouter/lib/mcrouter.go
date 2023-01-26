package mpmcrouter

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var cmdMetricNames = []string{
	"cmd_add_count",
	"cmd_cas_count",
	"cmd_decr_count",
	"cmd_delete_count",
	"cmd_get_count",
	"cmd_gets_count",
	"cmd_incr_count",
	"cmd_lease_get_count",
	"cmd_lease_set_count",
	"cmd_meta_count",
	"cmd_other_count",
	"cmd_replace_count",
	"cmd_set_count",
	"cmd_stats_count",
}

var resultMetricNames = []string{
	"result_busy_all_count",
	"result_busy_count",
	"result_connect_error_all_count",
	"result_connect_error_count",
	"result_connect_timeout_all_count",
	"result_connect_timeout_count",
	"result_data_timeout_all_count",
	"result_data_timeout_count",
	"result_error_all_count",
	"result_error_count",
	"result_local_error_all_count",
	"result_local_error_count",
	"result_tko_all_count",
	"result_tko_count",
}

// McrouterPlugin mackerel plugin
type McrouterPlugin struct {
	Prefix    string
	StatsFile string
}

// MetricKeyPrefix interface for mackerelplugin
func (p McrouterPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "mcrouter"
	}
	return p.Prefix
}

// GraphDefinition interface for mackerelplugin
func (p McrouterPlugin) GraphDefinition() map[string]mp.Graphs {
	var cmdMetrics []mp.Metrics
	for _, name := range cmdMetricNames {
		cmdMetrics = append(cmdMetrics, mp.Metrics{
			Name:  name,
			Label: name,
			Diff:  true,
		})
	}

	var resultMetrics []mp.Metrics
	for _, name := range resultMetricNames {
		resultMetrics = append(resultMetrics, mp.Metrics{
			Name:  name,
			Label: name,
			Diff:  true,
		})
	}

	labelPrefix := cases.Title(language.Und, cases.NoLower).String(p.Prefix)
	return map[string]mp.Graphs{
		"cmd_count": {
			Label:   (labelPrefix + " Receive Command"),
			Unit:    "integer",
			Metrics: cmdMetrics,
		},
		"result_count": {
			Label:   (labelPrefix + " Reply Result"),
			Unit:    "integer",
			Metrics: resultMetrics,
		},
		"request_processing_time": {
			Label: (labelPrefix + " Request Processing Time"),
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "duration_us", Label: "Duration(us)", Diff: false},
			},
		},
	}
}

// FetchMetrics interface for mackerelplugin
func (p McrouterPlugin) FetchMetrics() (map[string]interface{}, error) {
	stats, err := readStatsFile(p.StatsFile)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{})
	statsPrefix := strings.TrimSuffix(path.Base(p.StatsFile), ".stats")

	// Get cmd_[operation]_count stats
	for _, name := range cmdMetricNames {
		key := fmt.Sprintf("%s.%s", statsPrefix, name)
		ret[name] = stats[key]
	}

	// Get result_[reply result]_count stats
	for _, name := range resultMetricNames {
		key := fmt.Sprintf("%s.%s", statsPrefix, name)
		ret[name] = stats[key]
	}

	// Get duration_us stats
	for _, name := range []string{"duration_us"} {
		key := fmt.Sprintf("%s.%s", statsPrefix, name)
		ret[name] = stats[key]
	}

	return ret, nil
}

func readStatsFile(statsFile string) (map[string]float64, error) {
	data, err := os.ReadFile(statsFile)
	if err != nil {
		return nil, err
	}

	var stats map[string]float64
	err = json.Unmarshal(data, &stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// Do the plugin
func Do() {
	var (
		optStatsFile = flag.String("stats-file", "", "Mcrouter stats file")
		optPrefix    = flag.String("metric-key-prefix", "mcrouter", "Metric key prefix")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -stats-file /path/to/mcrouter.stats [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *optStatsFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	if _, err := os.Stat(*optStatsFile); err != nil {
		fmt.Fprintf(os.Stderr, "Stats file not found\n")
		flag.Usage()
		os.Exit(1)
	}

	helper := mp.NewMackerelPlugin(McrouterPlugin{
		Prefix:    *optPrefix,
		StatsFile: *optStatsFile,
	})
	helper.Run()
}
