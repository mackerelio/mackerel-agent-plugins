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
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// metric value structure
var graphdef = map[string](mp.Graphs){
	"apache2.workers": mp.Graphs{
		Label: "Apache Workers",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "busy_workers", Label: "Busy Workers", Diff: false, Stacked: true},
			mp.Metrics{Name: "idle_workers", Label: "Idle Workers", Diff: false, Stacked: true},
		},
	},
	"apache2.bytes": mp.Graphs{
		Label: "Apache Bytes",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "bytes_sent", Label: "Bytes Sent", Diff: true, Type: "uint64"},
		},
	},
	"apache2.cpu": mp.Graphs{
		Label: "Apache CPU Load",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "cpu_load", Label: "CPU Load", Diff: false},
		},
	},
	"apache2.req": mp.Graphs{
		Label: "Apache Requests",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "requests", Label: "Requests", Diff: true, Type: "uint64"},
		},
	},
	"apache2.scoreboard": mp.Graphs{
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

// Apache2Plugin for fetching metrics
type Apache2Plugin struct {
	Host     string
	Port     uint16
	Path     string
	Header   []string
	Tempfile string
}

// GraphDefinition Graph definition
func (c Apache2Plugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

// main function
func doMain(c *cli.Context) {

	var apache2 Apache2Plugin

	apache2.Host = c.String("http_host")
	apache2.Port = uint16(c.Int("http_port"))
	apache2.Path = c.String("status_page")
	apache2.Header = c.StringSlice("header")

	helper := mp.NewMackerelPlugin(apache2)
	helper.Tempfile = c.String("tempfile")

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

// FetchMetrics fetch the metrics
func (c Apache2Plugin) FetchMetrics() (map[string]interface{}, error) {
	data, err := getApache2Metrics(c.Host, c.Port, c.Path, c.Header)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]interface{})
	errStat := parseApache2Status(data, &stat)
	if errStat != nil {
		return nil, errStat
	}
	errScore := parseApache2Scoreboard(data, &stat)
	if errScore != nil {
		return nil, errScore
	}

	return stat, nil
}

// parsing scoreboard from server-status?auto
func parseApache2Scoreboard(str string, p *map[string]interface{}) error {
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
			c, assert := (*p)[name].(float64)
			if !assert {
				c = 0.0
			}
			(*p)[name] = c + 1.0
		}
		return nil
	}

	return errors.New("Scoreboard data is not found.")
}

// parsing metrics from server-status?auto
func parseApache2Status(str string, p *map[string]interface{}) error {
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
		var errParse error
		(*p)[Params[record[0]]], errParse = strconv.ParseFloat(strings.Trim(record[1], " "), 64)
		if errParse != nil {
			return errParse
		}
	}

	if len(*p) == 0 {
		return errors.New("Status data not found.")
	}

	return nil
}

// Getting apache2 status from server-status module data.
func getApache2Metrics(host string, port uint16, path string, header []string) (string, error) {
	uri := "http://" + host + ":" + strconv.FormatUint(uint64(port), 10) + path
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}
	for _, h := range header {
		kv := strings.SplitN(h, ":", 2)
		var k, v string
		k = strings.TrimSpace(kv[0])
		if len(kv) == 2 {
			v = strings.TrimSpace(kv[1])
		}
		if http.CanonicalHeaderKey(k) == "Host" {
			req.Host = v
		} else {
			req.Header.Set(k, v)
		}
	}
	resp, err := http.DefaultClient.Do(req)
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

// main
func main() {
	app := cli.NewApp()
	app.Name = "apache2_metrics"
	app.Version = version
	app.Usage = "Get metrics from apache2."
	app.Author = "Yuichiro Saito"
	app.Email = "saito@heartbeats.jp"
	app.Flags = flags
	app.Action = doMain

	app.Run(os.Args)
}
