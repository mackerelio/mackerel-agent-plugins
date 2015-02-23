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

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"php-opcache.memory_size": mp.Graphs{
		Label: "PHP OPCache Memory Size",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "used_memory", Label: "Used Memory", Diff: false, Stacked: false},
			mp.Metrics{Name: "free_memory", Label: "Free Memory", Diff: false, Stacked: false},
			mp.Metrics{Name: "wasted_memory", Label: "Wasted Memory", Diff: false, Stacked: false},
		},
	},
	"php-opcache.memory": mp.Graphs{
		Label: "PHP OPCache Memory Statistics",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "opcache_hit_rate", Label: "OPCache hit rate", Diff: false, Stacked: false},
			mp.Metrics{Name: "current_wasted_percentage", Label: "Used Memory", Diff: false, Stacked: false},
		},
	},

	"php-opcache.cache_size": mp.Graphs{
		Label: "PHP OPCache Cache Size",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "num_cached_scripts", Label: "Cached Script", Diff: true, Stacked: false},
			mp.Metrics{Name: "num_cached_keys", Label: "Num Cached Key", Diff: true, Stacked: false},
			mp.Metrics{Name: "max_cached_keys", Label: "Max Cached Key", Diff: true, Stacked: false},
		},
	},
	"php-opcache.stats": mp.Graphs{
		Label: "PHP OPCache Cache Statistics",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "hits", Label: "Hits", Diff: true, Stacked: false},
			mp.Metrics{Name: "misses", Label: "Misses", Diff: true, Stacked: false},
			mp.Metrics{Name: "blacklist_misses", Label: "Blacklist Misses", Diff: true, Stacked: false},
		},
	},
}

type PhpOpcachePlugin struct {
	Host     string
	Port     uint16
	Path     string
	Tempfile string
}

func (c PhpOpcachePlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func (c PhpOpcachePlugin) FetchMetrics() (map[string]float64, error) {
	data, err := getPhpOpcacheMetrics(c.Host, c.Port, c.Path)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]float64)
	err_stat := parsePhpOpcacheStatus(data, &stat)
	if err_stat != nil {
		return nil, err_stat
	}

	return stat, nil
}

func parsePhpOpcacheStatus(str string, p *map[string]float64) error {
	for _, line := range strings.Split(str, "\n") {
		record := strings.Split(line, ":")
		if len(record) != 2 {
			continue
		}
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

func getPhpOpcacheMetrics(host string, port uint16, path string) (string, error) {
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

func doMain(c *cli.Context) {
	var phpopcache PhpOpcachePlugin

	phpopcache.Host = c.String("http_host")
	phpopcache.Port = uint16(c.Int("http_port"))
	phpopcache.Path = c.String("status_page")

	helper := mp.NewMackerelPlugin(phpopcache)
	helper.Tempfile = c.String("tempfile")

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "php-opcache_metrics"
	app.Version = Version
	app.Usage = "Get metrics from php-opcache."
	app.Author = "Yuichiro Mukai"
	app.Email = "y.iky917@gmail.com"
	app.Flags = Flags
	app.Action = doMain

	app.Run(os.Args)
}
