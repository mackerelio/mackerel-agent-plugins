package main

import (
	"bufio"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var index map[string](int) = map[string](int){
	"NAME":       0,
	"STATE":      1,
	"CPU_SEC":    2,
	"CPU_PER":    3,
	"MEM_K":      4,
	"MEM_PER":    5,
	"MAXMEM_K":   6,
	"MAXMEM_PER": 7,
	"VCPUS":      8,
	"NETS":       9,
	"NETTX":      10,
	"NETRX":      11,
	"VBDS":       12,
	"VBD_OO":     13,
	"VBD_RD":     14,
	"VBD_WR":     15,
	"VBD_RSECT":  16,
	"VBD_WSECT":  17,
	"SSID":       18,
}

// All metrics are added dinamically at GraphDefinition
var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){}

type XentopMetrics struct {
	HostName string
	Metrics  mp.Metrics
}

type XentopPlugin struct {
	GraphName          string
	GraphUnit          string
	XentopMetricsSlice []XentopMetrics
}

func (m XentopPlugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)
	cmd := exec.Command("/bin/sh", "-c", "sudo xentop --batch -i 1 -f")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	scanner.Scan()
	for scanner.Scan() {
		sf := strings.Fields(string(scanner.Text()))
		name := sf[index["NAME"]]

		var err_parse error
		stat[fmt.Sprintf("cpu_%s", name)], err_parse = strconv.ParseFloat(sf[index["CPU_PER"]], 64)
		if err_parse != nil {
			return nil, err_parse
		}
		stat[fmt.Sprintf("memory_%s", name)], err_parse = strconv.ParseFloat(sf[index["MEM_PER"]], 64)
		if err_parse != nil {
			return nil, err_parse
		}
		stat[fmt.Sprintf("nettx_%s", name)], err_parse = strconv.ParseFloat(sf[index["NETTX"]], 64) * 1000
		if err_parse != nil {
			return nil, err_parse
		}
		stat[fmt.Sprintf("netrx_%s", name)], err_parse = strconv.ParseFloat(sf[index["NETRX"]], 64) * 1000
		if err_parse != nil {
			return nil, err_parse
		}
		stat[fmt.Sprintf("vbdrd_%s", name)], err_parse = strconv.ParseFloat(sf[index["VBD_RD"]], 64)
		if err_parse != nil {
			return nil, err_parse
		}
		stat[fmt.Sprintf("vbdwr_%s", name)], err_parse = strconv.ParseFloat(sf[index["VBD_WR"]], 64)
		if err_parse != nil {
			return nil, err_parse
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return stat, nil
}

func DefineCpuMetrics(names []string) []mp.Metrics {
	cpu_metrics := make([]mp.Metrics, 0)
	for _, name := range names {
		cpu_metrics = append(cpu_metrics, mp.Metrics{Name: fmt.Sprintf("cpu_%s", name), Label: name, Stacked: true})
	}
	return cpu_metrics
}

func DefineMemoryMetrics(names []string) []mp.Metrics {
	memory_metrics := make([]mp.Metrics, 0)
	for _, name := range names {
		memory_metrics = append(memory_metrics, mp.Metrics{Name: fmt.Sprintf("memory_%s", name), Label: name, Stacked: true})
	}
	return memory_metrics
}

func DefineNettxMetrics(names []string) []mp.Metrics {
	nettx_metrics := make([]mp.Metrics, 0)
	for _, name := range names {
		nettx_metrics = append(nettx_metrics, mp.Metrics{Name: fmt.Sprintf("nettx_%s", name), Label: name, Stacked: true})
	}
	return nettx_metrics
}

func DefineNetrxMetrics(names []string) []mp.Metrics {
	netrx_metrics := make([]mp.Metrics, 0)
	for _, name := range names {
		netrx_metrics = append(netrx_metrics, mp.Metrics{Name: fmt.Sprintf("netrx_%s", name), Label: name, Stacked: true})
	}
	return netrx_metrics
}

func DefineVbdrdMetrics(names []string) []mp.Metrics {
	vbdrd_metrics := make([]mp.Metrics, 0)
	for _, name := range names {
		vbdrd_metrics = append(vbdrd_metrics, mp.Metrics{Name: fmt.Sprintf("vbdrd_%s", name), Label: name, Stacked: true})
	}
	return vbdrd_metrics
}

func DefineVbdwrMetrics(names []string) []mp.Metrics {
	vbdwr_metrics := make([]mp.Metrics, 0)
	for _, name := range names {
		vbdwr_metrics = append(vbdwr_metrics, mp.Metrics{Name: fmt.Sprintf("vbdwr_%s", name), Label: name, Stacked: true})
	}
	return vbdwr_metrics
}

func DefineGraphs(names []string) {
	graphdef["xentop.cpu"] = mp.Graphs{
		Label:   "Xentop CPU",
		Unit:    "percentage",
		Metrics: DefineCpuMetrics(names),
	}
	graphdef["xentop.memory"] = mp.Graphs{
		Label:   "Xentop Memory",
		Unit:    "percentage",
		Metrics: DefineMemoryMetrics(names),
	}
	graphdef["xentop.nettx"] = mp.Graphs{
		Label:   "Xentop Nettx",
		Unit:    "bytes",
		Metrics: DefineNettxMetrics(names),
	}
	graphdef["xentop.netrx"] = mp.Graphs{
		Label:   "Xentop Netrx",
		Unit:    "bytes",
		Metrics: DefineNetrxMetrics(names),
	}
	graphdef["xentop.vbdrd"] = mp.Graphs{
		Label:   "Xentop VBD_RD",
		Unit:    "iops",
		Metrics: DefineVbdrdMetrics(names),
	}
	graphdef["xentop.vbdwr"] = mp.Graphs{
		Label:   "Xentop VBD_WR",
		Unit:    "iops",
		Metrics: DefineVbdwrMetrics(names),
	}
}

func (m XentopPlugin) GraphDefinition() map[string](mp.Graphs) {
	cmd := exec.Command("/bin/sh", "-c", "xentop --batch -i 1 -f")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd.Start()

	names := make([]string, 0)
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		sf := strings.Fields(string(scanner.Text()))
		if sf[index["NAME"]] != "NAME" {
			name := sf[index["NAME"]]
			names = append(names, name)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	DefineGraphs(names)

	return graphdef

}

func main() {
	// TODO: flagの取得

	var xentop XentopPlugin

	helper := mp.NewMackerelPlugin(xentop)
	helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-xentop")

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
