package main

import (
	"github.com/codegangsta/cli"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// metric value structure
var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"rpi.cpu": mp.Graphs{
		Label: "Raspberry PI CPU Temperature",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "temperature", Label: "CPU Temperature", Diff: false},
		},
	},
}

// for fetching metrics
type RpiPlugin struct {
	TemperaturePath string
	Tempfile        string
}

// Graph definition
func (c RpiPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

// main function
func doMain(c *cli.Context) {

	var rpi RpiPlugin

	rpi.TemperaturePath = c.String("temperature_path")

	helper := mp.NewMackerelPlugin(rpi)
	helper.Tempfile = c.String("tempfile")

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

// fetch metrics
func (c RpiPlugin) FetchMetrics() (map[string]float64, error) {
	data, err := ioutil.ReadFile(c.TemperaturePath)
	if err != nil {
		return nil, err
	}

	temp_str := strings.TrimSpace(string(data))
	temperature, err := strconv.ParseFloat(temp_str, 64)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]float64)
	stat["temperature"] = temperature / 1000

	return stat, nil
}

// main
func main() {
	app := cli.NewApp()
	app.Name = "rpi_metrics"
	app.Version = Version
	app.Usage = "Get metrics from Raspberry PI."
	app.Author = "Takuya Arita"
	app.Email = "takuya.arita@gmail.com"
	app.Flags = Flags
	app.Action = doMain

	app.Run(os.Args)
}
