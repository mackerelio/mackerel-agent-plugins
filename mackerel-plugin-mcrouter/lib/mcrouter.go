package mpmcrouter

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
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

	labelPrefix := strings.Title(p.Prefix)
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
		value := getStats(stats, statsPrefix, name)
		ret[name], err = strconv.ParseUint(string(value), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	// Get result_[reply result]_count stats
	for _, name := range resultMetricNames {
		value := getStats(stats, statsPrefix, name)
		ret[name], err = strconv.ParseUint(string(value), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	// Get duration_us stats
	for _, name := range []string{"duration_us"} {
		value := getStats(stats, statsPrefix, name)
		ret[name], err = strconv.ParseFloat(string(value), 64)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func readStatsFile(statsFile string) (interface{}, error) {
	data, err := ioutil.ReadFile(statsFile)
	if err != nil {
		return nil, err
	}

	var stats interface{}
	d := json.NewDecoder(bytes.NewBuffer(data))
	d.UseNumber()
	if err := d.Decode(&stats); err != nil {
		return nil, err
	}

	return stats, nil
}

func getStats(stats interface{}, statsPrefix string, metricName string) json.Number {
	statsKey := fmt.Sprintf("%s.%s", statsPrefix, metricName)
	return stats.(map[string]interface{})[statsKey].(json.Number)
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
