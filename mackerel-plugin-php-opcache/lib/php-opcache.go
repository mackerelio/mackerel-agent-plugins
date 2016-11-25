package mpphpopcache

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

var graphdef = map[string]mp.Graphs{
	"php-opcache.memory_size": {
		Label: "PHP OPCache Memory Size",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "used_memory", Label: "Used Memory", Diff: false, Stacked: false},
			{Name: "free_memory", Label: "Free Memory", Diff: false, Stacked: false},
			{Name: "wasted_memory", Label: "Wasted Memory", Diff: false, Stacked: false},
		},
	},
	"php-opcache.memory": {
		Label: "PHP OPCache Memory Statistics",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "opcache_hit_rate", Label: "OPCache hit rate", Diff: false, Stacked: false},
			{Name: "current_wasted_percentage", Label: "Used Memory", Diff: false, Stacked: false},
		},
	},

	"php-opcache.cache_size": {
		Label: "PHP OPCache Cache Size",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "num_cached_scripts", Label: "Cached Script", Diff: true, Stacked: false},
			{Name: "num_cached_keys", Label: "Num Cached Key", Diff: true, Stacked: false},
			{Name: "max_cached_keys", Label: "Max Cached Key", Diff: true, Stacked: false},
		},
	},
	"php-opcache.stats": {
		Label: "PHP OPCache Cache Statistics",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "hits", Label: "Hits", Diff: true, Stacked: false},
			{Name: "misses", Label: "Misses", Diff: true, Stacked: false},
			{Name: "blacklist_misses", Label: "Blacklist Misses", Diff: true, Stacked: false},
		},
	},
}

// PhpOpcachePlugin mackerel plugin for php-opcache
type PhpOpcachePlugin struct {
	Host     string
	Port     uint16
	Path     string
	Tempfile string
}

// GraphDefinition interface for mackerelplugin
func (c PhpOpcachePlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// FetchMetrics interface for mackerelplugin
func (c PhpOpcachePlugin) FetchMetrics() (map[string]float64, error) {
	data, err := getPhpOpcacheMetrics(c.Host, c.Port, c.Path)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]float64)
	errStat := parsePhpOpcacheStatus(data, &stat)
	if errStat != nil {
		return nil, errStat
	}

	return stat, nil
}

func parsePhpOpcacheStatus(str string, p *map[string]float64) error {
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

func getPhpOpcacheMetrics(host string, port uint16, path string) (string, error) {
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

func doMain(c *cli.Context) error {
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
	return nil
}

// Do the plugin
func Do() {
	app := cli.NewApp()
	app.Name = "php-opcache_metrics"
	app.Version = version
	app.Usage = "Get metrics from php-opcache."
	app.Author = "Yuichiro Mukai"
	app.Email = "y.iky917@gmail.com"
	app.Flags = flags
	app.Action = doMain

	app.Run(os.Args)
}
