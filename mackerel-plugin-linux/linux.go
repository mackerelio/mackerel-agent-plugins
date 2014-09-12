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

	var apache2 LinuxPlugin

	apache2.Tempfile = c.String("tempfile")

	helper := mp.NewMackerelPlugin(apache2)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

// fetch metrics
func (c Plugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)

	data, err := getProcVmstat()
	if err != nil {
		return nil, err
	}
	err_stat := parseProcVmstat(data, &stat)
	if err_stat != nil {
		return nil, err_stat
	}

	return stat, nil
}


// parsing metrics from server-status?auto
func parseProcVmstat(str string, p *map[string]float64) error {
	for _, line := range strings.Split(str, "\n") {
		record := strings.Split(line, " ")
		var err_parse error
		(*p)[Params[record[0]]], err_parse = strconv.ParseFloat(strings.Trim(record[1], " "), 64)
		if err_parse != nil {
			return err_parse
		}
	}

	return nil
}

// Getting apache2 status from server-status module data.
func getProcVmstat( path string ) ( string, error ) {
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
