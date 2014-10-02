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

const PathVmstat = "/proc/vmstat"
const PathDiskstats = "/proc/diskstats"
const PathStat = "/proc/stat"

// metric value structure
// note: all metrics are add dynamic at collect*().
var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){}

// for fetching metrics
type LinuxPlugin struct {
	Tempfile string
	Type     string
}

// Graph definition
func (c LinuxPlugin) GraphDefinition() map[string](mp.Graphs) {
	var err error

	p := make(map[string]float64)

	if c.Type == "all" || c.Type == "swap" {
		err = collectProcVmstat(PathVmstat, &p)
		if err != nil {
			return nil
		}
	}

	if c.Type == "all" || c.Type == "netstat" {
		err = collectSs(&p)
		if err != nil {
			return nil
		}
	}

	if c.Type == "all" || c.Type == "diskstats" {
		err = collectProcDiskstats(PathDiskstats, &p)
		if err != nil {
			return nil
		}
	}

	if c.Type == "all" || c.Type == "proc_stat" {
		err = collectProcStat(PathStat, &p)
		if err != nil {
			return nil
		}
	}

	if c.Type == "all" || c.Type == "users" {
		err = collectWho(&p)
		if err != nil {
			return nil
		}
	}

	return graphdef
}

// main function
func doMain(c *cli.Context) {
	var linux LinuxPlugin

	linux.Type = c.String("type")
	helper := mp.NewMackerelPlugin(linux)
	helper.Tempfile = c.String("tempfile")

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

// fetch metrics
func (c LinuxPlugin) FetchMetrics() (map[string]float64, error) {
	var err error

	p := make(map[string]float64)

	if c.Type == "all" || c.Type == "swap" {
		err = collectProcVmstat(PathVmstat, &p)
		if err != nil {
			return nil, err
		}
	}

	if c.Type == "all" || c.Type == "netstat" {
		err = collectSs(&p)
		if err != nil {
			return nil, err
		}
	}

	if c.Type == "all" || c.Type == "diskstats" {
		err = collectProcDiskstats(PathDiskstats, &p)
		if err != nil {
			return nil, err
		}
	}

	if c.Type == "all" || c.Type == "proc_stat" {
		err = collectProcStat(PathStat, &p)
		if err != nil {
			return nil, err
		}
	}

	if c.Type == "all" || c.Type == "users" {
		err = collectWho(&p)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

// collect who
func collectWho(p *map[string]float64) error {
	var err error
	var data string

	graphdef["linux.users"] = mp.Graphs{
		Label: "Linux Users",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "users", Label: "Users", Diff: false},
		},
	}

	data, err = getWho()
	if err != nil {
		return err
	}
	err = parseWho(data, p)
	if err != nil {
		return err
	}

	return nil
}

// parsing metrics from /proc/stat
func parseWho(str string, p *map[string]float64) error {
	if strings.TrimSpace(str) == "" {
		(*p)["users"] = 0
		return nil
	}
	line := strings.Split(str, "\n")
	(*p)["users"] = float64(len(line))

	return nil
}

// Getting who
func getWho() (string, error) {
	cmd := exec.Command("who")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// collect /proc/stat
func collectProcStat(path string, p *map[string]float64) error {
	var err error
	var data string

	graphdef["linux.interrupts"] = mp.Graphs{
		Label: "Linux Interrupts",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "interrupts", Label: "Interrupts", Diff: true},
		},
	}
	graphdef["linux.context_switches"] = mp.Graphs{
		Label: "Linux Context Switches",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "context_switches", Label: "Context Switches", Diff: true},
		},
	}
	graphdef["linux.forks"] = mp.Graphs{
		Label: "Linux Forks",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "forks", Label: "Forks", Diff: true},
		},
	}

	data, err = getProc(path)
	if err != nil {
		return err
	}
	err = parseProcStat(data, p)
	if err != nil {
		return err
	}

	return nil
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

// collect /proc/diskstats
func collectProcDiskstats(path string, p *map[string]float64) error {
	var err error
	var data string

	data, err = getProc(path)
	if err != nil {
		return err
	}
	err = parseProcDiskstats(data, p)
	if err != nil {
		return err
	}

	return nil
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

// collect ss
func collectSs(p *map[string]float64) error {
	var err error
	var data string

	graphdef["linux.ss"] = mp.Graphs{
		Label: "Linux Network Connection States",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "ESTAB", Label: "Established", Diff: false, Stacked: true},
			mp.Metrics{Name: "SYN-SENT", Label: "Syn Sent", Diff: false, Stacked: true},
			mp.Metrics{Name: "SYN-RECV", Label: "Syn Received", Diff: false, Stacked: true},
			mp.Metrics{Name: "FIN-WAIT-1", Label: "Fin Wait 1", Diff: false, Stacked: true},
			mp.Metrics{Name: "FIN-WAIT-2", Label: "Fin Wait 2", Diff: false, Stacked: true},
			mp.Metrics{Name: "TIME-WAIT", Label: "Time Wait", Diff: false, Stacked: true},
			mp.Metrics{Name: "UNCONN", Label: "Close", Diff: false, Stacked: true},
			mp.Metrics{Name: "CLOSE-WAIT", Label: "Close Wait", Diff: false, Stacked: true},
			mp.Metrics{Name: "LAST-ACK", Label: "Last Ack", Diff: false, Stacked: true},
			mp.Metrics{Name: "LISTEN", Label: "Listen", Diff: false, Stacked: true},
			mp.Metrics{Name: "CLOSING", Label: "Closing", Diff: false, Stacked: true},
			mp.Metrics{Name: "UNKNOWN", Label: "Unknown", Diff: false, Stacked: true},
		},
	}
	data, err = getSs()
	if err != nil {
		return err
	}
	err = parseSs(data, p)
	if err != nil {
		return err
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

// collect /proc/vmstat
func collectProcVmstat(path string, p *map[string]float64) error {
	var err error
	var data string

	graphdef["linux.swap"] = mp.Graphs{
		Label: "Linux Swap Usage",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "pswpin", Label: "Swap In", Diff: true},
			mp.Metrics{Name: "pswpout", Label: "Swap Out", Diff: true},
		},
	}

	data, err = getProc(path)
	if err != nil {
		return err
	}
	err = parseProcVmstat(data, p)
	if err != nil {
		return err
	}

	return nil
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

// atof
func _atof(str string) (float64, error) {
	return strconv.ParseFloat(strings.Trim(str, " "), 64)
}

// main
func main() {
	app := cli.NewApp()
	app.Name = "mackerel-plugin-linux"
	app.Version = Version
	app.Usage = "Get metrics from Linux."
	app.Author = "Yuichiro Saito"
	app.Email = "saito@heartbeats.jp"
	app.Flags = Flags
	app.Action = doMain

	app.Run(os.Args)
}
