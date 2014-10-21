package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

// metric value structure
var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"php-apc.purges": mp.Graphs{
		Label: "APC purge count",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "cache_full_count", Label: "File Cache", Diff: true, Stacked: false},
			mp.Metrics{Name: "user_cache_full_count", Label: "User Cache", Diff: true, Stacked: false},
		},
	},
	"php-apc.stats": mp.Graphs{
		Label: "APC file cache statistics",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "cache_hits", Label: "Hits", Diff: true, Stacked: false},
			mp.Metrics{Name: "cache_misses", Label: "Misses", Diff: true, Stacked: false},
		},
	},
	"php-apc.cache_size": mp.Graphs{
		Label: "APC cache size",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "Code Cache", Label: "cached_files_size", Diff: false, Stacked: true},
			mp.Metrics{Name: "User Items Cache", Label: "user_cache_vars_size", Diff: false, Stacked: true},
			mp.Metrics{Name: "Limit", Label: "total_memory", Diff: false, Stacked: false},
		},
	},
	"php-apc.user_stats": mp.Graphs{
		Label: "APC user cache statistics",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "user_cache_hits", Label: "Hits", Diff: true, Stacked: false},
			mp.Metrics{Name: "user_cache_misses", Label: "Misses", Diff: true, Stacked: false},
		},
	},
}

// for fetching metrics
type PhpApcPlugin struct {
	Host     string
	Port     uint16
	Path     string
	Tempfile string
}

// Graph definition
func (c PhpApcPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

// main function
func doMain(c *cli.Context) {

	var phpapc PhpApcPlugin

	phpapc.Host = c.String("http_host")
	phpapc.Port = uint16(c.Int("http_port"))
	phpapc.Path = c.String("status_page")
	phpapc.Tempfile = c.String("tempfile")

	helper := mp.NewMackerelPlugin(phpapc)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

// fetch metrics
func (c PhpApcPlugin) FetchMetrics() (map[string]float64, error) {
	data, err := getPhpApcMetrics(c.Host, c.Port, c.Path)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]float64)
	err_stat := parsePhpApcStatus(data, &stat)
	if err_stat != nil {
		return nil, err_stat
	}

	return stat, nil
}

// parsing metrics from server-status?auto
func parsePhpApcStatus(str string, p *map[string]float64) error {
	for _, line := range strings.Split(str, "\n") {
		record := strings.Split(line, ":")
		var err_parse error
		(*p)[record[0]], err_parse = strconv.ParseFloat(strings.Trim(record[1], " "), 64)
		if err_parse != nil {
			return err_parse
		}
	}

	if len(*p) == 0 {
		return errors.New("Status data not found.")
	}

	return nil
}

// Getting php-apc status from server-status module data.
func getPhpApcMetrics(host string, port uint16, path string) (string, error) {
	uri := "http://" + host + ":" + strconv.FormatUint(uint64(port), 10) + path
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("HTTP status error: %d", resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(body[:]), nil
}

// main
func main() {
	app := cli.NewApp()
	app.Name = "php-apc_metrics"
	app.Version = Version
	app.Usage = "Get metrics from php-apc."
	app.Author = "Yuichiro Saito"
	app.Email = "saito@heartbeats.jp"
	app.Flags = Flags
	app.Action = doMain

	app.Run(os.Args)
}
