package mpphpapc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/urfave/cli"
)

// metric value structure
var graphdef = map[string]mp.Graphs{
	"php-apc.purges": {
		Label: "PHP APC Cache Purge Count",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "cache_full_count", Label: "File Cache", Diff: true, Stacked: false},
			{Name: "user_cache_full_count", Label: "User Cache", Diff: true, Stacked: false},
		},
	},
	"php-apc.stats": {
		Label: "PHP APC File Cache Statistics",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "cache_hits", Label: "Hits", Diff: true, Stacked: false},
			{Name: "cache_misses", Label: "Misses", Diff: true, Stacked: false},
		},
	},
	"php-apc.cache_size": {
		Label: "PHP APC Cache Size",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "cached_files_size", Label: "File Cache", Diff: false, Stacked: true},
			{Name: "user_cache_vars_size", Label: "User Cache", Diff: false, Stacked: true},
			{Name: "total_memory", Label: "Total", Diff: false, Stacked: false},
		},
	},
	"php-apc.user_stats": {
		Label: "PHP APC User Cache Statistics",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "user_cache_hits", Label: "Hits", Diff: true, Stacked: false},
			{Name: "user_cache_misses", Label: "Misses", Diff: true, Stacked: false},
		},
	},
}

// PhpApcPlugin mackerel plugin for php-apc
type PhpApcPlugin struct {
	Host     string
	Port     uint16
	Path     string
	Tempfile string
}

// GraphDefinition interface for mackerelplugin
func (c PhpApcPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// main function
func doMain(c *cli.Context) error {

	var phpapc PhpApcPlugin

	phpapc.Host = c.String("http_host")
	phpapc.Port = uint16(c.Int("http_port"))
	phpapc.Path = c.String("status_page")

	helper := mp.NewMackerelPlugin(phpapc)
	helper.Tempfile = c.String("tempfile")

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
	return nil
}

// FetchMetrics interface for mackerelplugin
func (c PhpApcPlugin) FetchMetrics() (map[string]float64, error) {
	data, err := getPhpApcMetrics(c.Host, c.Port, c.Path)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]float64)
	errStat := parsePhpApcStatus(data, &stat)
	if errStat != nil {
		return nil, errStat
	}

	return stat, nil
}

// parsing metrics from server-status?auto
func parsePhpApcStatus(str string, p *map[string]float64) error {
	for _, line := range strings.Split(str, "\n") {
		record := strings.Split(line, ":")
		if len(record) != 2 {
			continue
		}
		var errParse error
		(*p)[record[0]], errParse = strconv.ParseFloat(strings.Trim(record[1], " "), 64)
		if errParse != nil {
			return errParse
		}
	}

	if len(*p) == 0 {
		return errors.New("status data not found")
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
		return "", fmt.Errorf("HTTP status error: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(body[:]), nil
}

// Do the plugin
func Do() {
	app := cli.NewApp()
	app.Name = "php-apc_metrics"
	app.Version = version
	app.Usage = "Get metrics from php-apc."
	app.Author = "Yuichiro Saito"
	app.Email = "saito@heartbeats.jp"
	app.Flags = flags
	app.Action = doMain

	app.Run(os.Args)
}
