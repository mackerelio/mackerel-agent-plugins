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
	const PathDiskstats = "/proc/diskstats"
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

	data, err = getProcDiskstats(PathDiskstats)
	if err != nil {
		return nil, err
	}
	err = parseProcDiskstats(data, &stat)
	if err != nil {
		return nil, err
	}

	return stat, nil
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

// Getting diskstats
func getProcDiskstats(path string) (string, error) {
	cmd := exec.Command("cat", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
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
		(*p)[record[0]], err_parse = _atof(record[1])
		if err_parse != nil {
			return err_parse
		}
	}

	return nil
}

// Getting /proc/vmstat.
func getProcVmstat(path string) (string, error) {
	cmd := exec.Command("cat", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
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
