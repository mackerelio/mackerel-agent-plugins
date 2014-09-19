package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

// metric value structure
var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"linux.swap": mp.Graphs{
		Label: "Linux Swap Usage",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "pswpin", Label: "Swap In", Diff: false},
			mp.Metrics{Name: "pswpout", Label: "Swap Out", Diff: false},
		},
	},
}

// for fetching metrics
type LinuxPlugin struct {
	Tempfile string
}

// Graph definition
func (c LinuxPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

// main function
func doMain(c *cli.Context) {

	var linux LinuxPlugin

	linux.Tempfile = c.String("tempfile")

	helper := mp.NewMackerelPlugin(linux)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

// fetch metrics
func (c LinuxPlugin) FetchMetrics() (map[string]float64, error) {
	const PathVmstat = "/proc/vmstat"
	var err error

	stat := make(map[string]float64)

	var data string
	data, err = getProcVmstat(PathVmstat)
	if err != nil {
		return nil, err
	}
	err = parseProcVmstat(data, &stat)
	if err != nil {
		return nil, err
	}

	return stat, nil
}

// parsing metrics from /proc/vmstat
func parseProcVmstat(str string, p *map[string]float64) error {
	for _, line := range strings.Split(str, "\n") {
		record := strings.Split(line, " ")
		if len(record) != 2 {
			continue
		}
		var err_parse error
		(*p)[record[0]], err_parse = strconv.ParseFloat(strings.Trim(record[1], " "), 64)
		if err_parse != nil {
			return err_parse
		}
	}

	return nil
}

// Getting /proc/vmstat.
func getProcVmstat(path string) (string, error) {
	cmd := exec.Command("cat", "/proc/vmstat")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// main
func main() {
	app := cli.NewApp()
	app.Name = "mackerel-plugin-linux"
	app.Version = Version
	app.Usage = "Get metrics from apache2."
	app.Author = "Yuichiro Saito"
	app.Email = "saito@heartbeats.jp"
	app.Flags = Flags
	app.Action = doMain

	app.Run(os.Args)
}
