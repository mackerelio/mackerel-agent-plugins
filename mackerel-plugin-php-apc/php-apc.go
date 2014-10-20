package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

// metric value structure
var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"php-apc.workers": mp.Graphs{
		Label: "Apache Workers",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "busy_workers", Label: "Busy Workers", Diff: false, Stacked: true},
			mp.Metrics{Name: "idle_workers", Label: "Idle Workers", Diff: false, Stacked: true},
		},
	},
	"php-apc.bytes": mp.Graphs{
		Label: "Apache Bytes",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "bytes_sent", Label: "Bytes Sent", Diff: false},
		},
	},
	"php-apc.cpu": mp.Graphs{
		Label: "Apache CPU Load",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "cpu_load", Label: "CPU Load", Diff: false},
		},
	},
	"php-apc.req": mp.Graphs{
		Label: "Apache Requests",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "requests", Label: "Requests", Diff: false},
		},
	},
	"php-apc.scoreboard": mp.Graphs{
		Label: "Apache Scoreboard",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "score-_", Label: "Waiting for connection", Diff: false, Stacked: true},
			mp.Metrics{Name: "score-S", Label: "Starting up", Diff: false, Stacked: true},
			mp.Metrics{Name: "score-R", Label: "Reading request", Diff: false, Stacked: true},
			mp.Metrics{Name: "scpre-W", Label: "Sending reply", Diff: false, Stacked: true},
			mp.Metrics{Name: "score-K", Label: "Keepalive", Diff: false, Stacked: true},
			mp.Metrics{Name: "score-D", Label: "DNS lookup", Diff: false, Stacked: true},
			mp.Metrics{Name: "score-C", Label: "Closing connection", Diff: false, Stacked: true},
			mp.Metrics{Name: "score-L", Label: "Logging", Diff: false, Stacked: true},
			mp.Metrics{Name: "score-G", Label: "Gracefully finishing", Diff: false, Stacked: true},
			mp.Metrics{Name: "score-I", Label: "Idle cleanup", Diff: false, Stacked: true},
			mp.Metrics{Name: "score-.", Label: "Open slot", Diff: false, Stacked: true},
		},
	},
}

// for fetching metrics
type Apache2Plugin struct {
	Host     string
	Port     uint16
	Path     string
	Tempfile string
}

// Graph definition
func (c Apache2Plugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

// main function
func doMain(c *cli.Context) {

	var php-apc Apache2Plugin

	php-apc.Host = c.String("http_host")
	php-apc.Port = uint16(c.Int("http_port"))
	php-apc.Path = c.String("status_page")
	php-apc.Tempfile = c.String("tempfile")

	helper := mp.NewMackerelPlugin(php-apc)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

// fetch metrics
func (c Apache2Plugin) FetchMetrics() (map[string]float64, error) {
	data, err := getApache2Metrics(c.Host, c.Port, c.Path)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]float64)
	err_stat := parseApache2Status(data, &stat)
	if err_stat != nil {
		return nil, err_stat
	}
	err_score := parseApache2Scoreboard(data, &stat)
	if err_score != nil {
		return nil, err_score
	}

	return stat, nil
}

// parsing scoreboard from server-status?auto
func parseApache2Scoreboard(str string, p *map[string]float64) error {
	for _, line := range strings.Split(str, "\n") {
		matched, err := regexp.MatchString("Scoreboard(.*)", line)
		if err != nil {
			return err
		}
		if !matched {
			continue
		}
		record := strings.Split(line, ":")
		for _, sb := range strings.Split(strings.Trim(record[1], " "), "") {
			name := fmt.Sprintf("score-%s", sb)
			c, assert := (*p)[name]
			if !assert {
				c = 0
			}
			(*p)[name] = c + 1
		}
		return nil
	}

	return errors.New("Scoreboard data is not found.")
}

// parsing metrics from server-status?auto
func parseApache2Status(str string, p *map[string]float64) error {
	Params := map[string]string{
		"Total Accesses": "requests",
		"Total kBytes":   "bytes_sent",
		"CPULoad":        "cpu_load",
		"BusyWorkers":    "busy_workers",
		"IdleWorkers":    "idle_workers"}

	for _, line := range strings.Split(str, "\n") {
		record := strings.Split(line, ":")
		_, assert := Params[record[0]]
		if !assert {
			continue
		}
		var err_parse error
		(*p)[Params[record[0]]], err_parse = strconv.ParseFloat(strings.Trim(record[1], " "), 64)
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
func getApache2Metrics(host string, port uint16, path string) (string, error) {
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
