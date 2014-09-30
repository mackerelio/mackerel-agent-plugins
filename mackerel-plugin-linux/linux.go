package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

// metric value structure
// note: diskstats are add dynamic at parseProcDiskstats().
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
		Label: "Linux Network Connection States",
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
	"linux.interrupts": mp.Graphs{
		Label: "Linux Interrupts",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "interrupts", Label: "Interrupts", Diff: true},
		},
	},
	"linux.context_switches": mp.Graphs{
		Label: "Linux Context Switches",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "context_switches", Label: "Context Switches", Diff: true},
		},
	},
	"linux.forks": mp.Graphs{
		Label: "Linux Forks",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "forks", Label: "Forks", Diff: true},
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
	const PathDiskstats = "/proc/diskstats"
	const PathStat = "/proc/stat"
	var err error
	var data string

	stat := make(map[string]float64)

	data, err = getProc(PathVmstat)
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

	data, err = getProc(PathDiskstats)
	if err != nil {
		return nil, err
	}
	err = parseProcDiskstats(data, &stat)
	if err != nil {
		return nil, err
	}

	data, err = getProc(PathStat)
	if err != nil {
		return nil, err
	}
	err = parseProcStat(data, &stat)
	if err != nil {
		return nil, err
	}

	return stat, nil
}

// parsing metrics from /proc/stat
func parseProcStat(str string, p *map[string]float64) error {
	for _, line := range strings.Split(str, "\n") {
		record := strings.Fields(line)
		if len(record) < 2 {
			continue
		}
		name := record[0]
		value, err_parse := _atof(record[1])
		if err_parse != nil {
			return err_parse
		}

		if name == "intr" {
			(*p)["interrupts"] = value
		} else if name == "ctxt" {
			(*p)["context_switches"] = value
		} else if name == "processes" {
			(*p)["forks"] = value
		}
	}

	return nil
}

// Getting /proc/*
func getProc(path string) (string, error) {
	cmd := exec.Command("cat", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// parsing metrics from diskstats
func parseProcDiskstats(str string, p *map[string]float64) error {
	for _, line := range strings.Split(str, "\n") {
		// See also: https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
		record := strings.Fields(line)
		if len(record) < 14 {
			continue
		}
		device := record[2]
		matched, err := regexp.MatchString("[0-9]$", device)
		if matched || err != nil {
			continue
		}

		(*p)[fmt.Sprintf("iotime_%s", device)], _ = _atof(record[12])
		(*p)[fmt.Sprintf("iotime_weighted_%s", device)], _ = _atof(record[13])
		graphdef[fmt.Sprintf("linux.disk.elapsed.%s", device)] = mp.Graphs{
			Label: fmt.Sprintf("Disk Elapsed IO Time %s", device),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: fmt.Sprintf("iotime_%s", device), Label: "IO Time", Diff: true},
				mp.Metrics{Name: fmt.Sprintf("iotime_weighted_%s", device), Label: "IO Time Weighted", Diff: true},
			},
		}

		(*p)[fmt.Sprintf("tsreading_%s", device)], _ = _atof(record[6])
		(*p)[fmt.Sprintf("tswriting_%s", device)], _ = _atof(record[10])
		graphdef[fmt.Sprintf("linux.disk.rwtime.%s", device)] = mp.Graphs{
			Label: fmt.Sprintf("Disk Read/Write Time %s", device),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: fmt.Sprintf("tsreading_%s", device), Label: "Read", Diff: true},
				mp.Metrics{Name: fmt.Sprintf("tswriting_%s", device), Label: "Write", Diff: true},
			},
		}
	}

	return nil
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
		record := strings.Fields(line)
		if len(record) != 2 {
			continue
		}
		var err_parse error
		(*p)[record[0]], err_parse = _atof(record[1])
		if err_parse != nil {
			return err_parse
		}
	}

	return nil
}

// atof
func _atof(str string) (float64, error) {
	return strconv.ParseFloat(strings.Trim(str, " "), 64)
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
