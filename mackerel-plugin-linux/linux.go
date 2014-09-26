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
	"linux.ss": mp.Graphs{
		Label: "Network Connection States",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "ESTAB", Label: "Established", Diff: false},
			mp.Metrics{Name: "SYN-SENT", Label: "Syn Sent", Diff: false},
			mp.Metrics{Name: "SYN-RECV", Label: "Syn Received", Diff: false},
			mp.Metrics{Name: "FIN-WAIT-1", Label: "Fin Wait 1", Diff: false},
			mp.Metrics{Name: "FIN-WAIT-2", Label: "Fin Wait 2", Diff: false},
			mp.Metrics{Name: "TIME-WAIT", Label: "Time Wait", Diff: false},
			mp.Metrics{Name: "UNCONN", Label: "Close", Diff: false},
			mp.Metrics{Name: "CLOSE-WAIT", Label: "Close Wait", Diff: false},
			mp.Metrics{Name: "LAST-ACK", Label: "Last Ack", Diff: false},
			mp.Metrics{Name: "LISTEN", Label: "Listen", Diff: false},
			mp.Metrics{Name: "CLOSING", Label: "Closing", Diff: false},
			mp.Metrics{Name: "UNKNOWN", Label: "Unknown", Diff: false},
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
	var data string

	stat := make(map[string]float64)

	data, err = getProcVmstat(PathVmstat)
	if err != nil {
		return nil, err
	}
	err = parseProcVmstat(data, &stat)
	if err != nil {
		return nil, err
	}

	data, err = getSs()
	if err != nil {
		return nil, err
	}
	err = parseSs(data, &stat)
	if err != nil {
		return nil, err
	}

	return stat, nil
}

// parsing metrics from ss
func parseSs(str string, p *map[string]float64) error {
	for i, line := range strings.Split(str, "\n") {
		if i < 1 {
			continue
		}
		record := strings.Fields(line)
		if len(record) != 5 {
			continue
		}
		(*p)[record[0]] = (*p)[record[0]] + 1
	}

	return nil
}

// Getting ss
func getSs() (string, error) {
	cmd := exec.Command("ss", "-na")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
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
