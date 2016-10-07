package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/urfave/cli"
)

const (
	pathVmstat    = "/proc/vmstat"
	pathDiskstats = "/proc/diskstats"
	pathStat      = "/proc/stat"
)

// metric value structure
// note: all metrics are add dynamic at collect*().
var graphdef = map[string](mp.Graphs){}

// LinuxPlugin mackerel plugin for linux
type LinuxPlugin struct {
	Tempfile string
	Typemap  map[string]bool
}

// GraphDefinition interface for mackerelplugin
func (c LinuxPlugin) GraphDefinition() map[string](mp.Graphs) {
	var err error

	p := make(map[string]interface{})

	if c.Typemap["all"] || c.Typemap["swap"] {
		err = collectProcVmstat(pathVmstat, &p)
		if err != nil {
			return nil
		}
	}

	if c.Typemap["all"] || c.Typemap["netstat"] {
		err = collectSs(&p)
		if err != nil {
			return nil
		}
	}

	if c.Typemap["all"] || c.Typemap["diskstats"] {
		err = collectProcDiskstats(pathDiskstats, &p)
		if err != nil {
			return nil
		}
	}

	if c.Typemap["all"] || c.Typemap["proc_stat"] {
		err = collectProcStat(pathStat, &p)
		if err != nil {
			return nil
		}
	}

	if c.Typemap["all"] || c.Typemap["users"] {
		err = collectWho(&p)
		if err != nil {
			return nil
		}
	}

	return graphdef
}

// main function
func doMain(c *cli.Context) error {
	var linux LinuxPlugin

	typemap := map[string]bool{}
	types := c.StringSlice("type")
	// If no `type` is specified, fetch all metrics
	if len(types) == 0 {
		typemap["all"] = true
	} else {
		for _, t := range types {
			typemap[t] = true
		}
	}
	linux.Typemap = typemap
	helper := mp.NewMackerelPlugin(linux)
	helper.Tempfile = c.String("tempfile")

	helper.Run()
	return nil
}

// FetchMetrics interface for mackerelplugin
func (c LinuxPlugin) FetchMetrics() (map[string]interface{}, error) {
	var err error

	p := make(map[string]interface{})

	if c.Typemap["all"] || c.Typemap["swap"] {
		err = collectProcVmstat(pathVmstat, &p)
		if err != nil {
			return nil, err
		}
	}

	if c.Typemap["all"] || c.Typemap["netstat"] {
		err = collectSs(&p)
		if err != nil {
			return nil, err
		}
	}

	if c.Typemap["all"] || c.Typemap["diskstats"] {
		err = collectProcDiskstats(pathDiskstats, &p)
		if err != nil {
			return nil, err
		}
	}

	if c.Typemap["all"] || c.Typemap["proc_stat"] {
		err = collectProcStat(pathStat, &p)
		if err != nil {
			return nil, err
		}
	}

	if c.Typemap["all"] || c.Typemap["users"] {
		err = collectWho(&p)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

// collect who
func collectWho(p *map[string]interface{}) error {
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
func parseWho(str string, p *map[string]interface{}) error {
	str = strings.TrimSpace(str)
	if str == "" {
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
func collectProcStat(path string, p *map[string]interface{}) error {
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
func parseProcStat(str string, p *map[string]interface{}) error {
	for _, line := range strings.Split(str, "\n") {
		record := strings.Fields(line)
		if len(record) < 2 {
			continue
		}
		name := record[0]
		value, errParse := atof(record[1])
		if errParse != nil {
			return errParse
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
func collectProcDiskstats(path string, p *map[string]interface{}) error {
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
func parseProcDiskstats(str string, p *map[string]interface{}) error {

	var elapsedData []mp.Metrics
	var rwtimeData []mp.Metrics

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

		(*p)[fmt.Sprintf("iotime_%s", device)], _ = atof(record[12])
		(*p)[fmt.Sprintf("iotime_weighted_%s", device)], _ = atof(record[13])
		elapsedData = append(elapsedData, mp.Metrics{Name: fmt.Sprintf("iotime_%s", device), Label: fmt.Sprintf("%s IO Time", device), Diff: true})
		elapsedData = append(elapsedData, mp.Metrics{Name: fmt.Sprintf("iotime_weighted_%s", device), Label: fmt.Sprintf("%s IO Time Weighted", device), Diff: true})

		(*p)[fmt.Sprintf("tsreading_%s", device)], _ = atof(record[6])
		(*p)[fmt.Sprintf("tswriting_%s", device)], _ = atof(record[10])
		rwtimeData = append(rwtimeData, mp.Metrics{Name: fmt.Sprintf("tsreading_%s", device), Label: fmt.Sprintf("%s Read", device), Diff: true})
		rwtimeData = append(rwtimeData, mp.Metrics{Name: fmt.Sprintf("tswriting_%s", device), Label: fmt.Sprintf("%s Write", device), Diff: true})
	}

	graphdef["linux.disk.elapsed"] = mp.Graphs{
		Label:   "Disk Elapsed IO Time",
		Unit:    "integer",
		Metrics: elapsedData,
	}

	graphdef["linux.disk.rwtime"] = mp.Graphs{
		Label:   "Disk Read/Write Time",
		Unit:    "integer",
		Metrics: rwtimeData,
	}

	return nil
}

// collect ss
func collectSs(p *map[string]interface{}) error {
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
func parseSs(str string, p *map[string]interface{}) error {
	status := 0
	for i, line := range strings.Split(str, "\n") {
		record := strings.Fields(line)
		if len(record) < 5 {
			continue
		}
		if i == 0 {
			if record[0] == "State" {
				// for RHEL6
				status = 0
			} else if record[1] == "State" {
				// for RHEL7
				status = 1
			}
		}
		v, _ := (*p)[record[status]].(float64)
		(*p)[record[status]] = v + 1
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
func collectProcVmstat(path string, p *map[string]interface{}) error {
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
func parseProcVmstat(str string, p *map[string]interface{}) error {
	for _, line := range strings.Split(str, "\n") {
		record := strings.Fields(line)
		if len(record) != 2 {
			continue
		}
		var errParse error
		(*p)[record[0]], errParse = atof(record[1])
		if errParse != nil {
			return errParse
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
func atof(str string) (float64, error) {
	return strconv.ParseFloat(strings.Trim(str, " "), 64)
}

// main
func main() {
	app := cli.NewApp()
	app.Name = "mackerel-plugin-linux"
	app.Version = version
	app.Usage = "Get metrics from Linux."
	app.Author = "Yuichiro Saito"
	app.Email = "saito@heartbeats.jp"
	app.Flags = flags
	app.Action = doMain

	app.Run(os.Args)
}
