package main

import (
	"bufio"
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"log"
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

// GraphDefinitionで動的に定義するので，これは多分必要なくなる
var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"xentop.cpu": mp.Graphs{
		Label:   "Xentop CPU",
		Unit:    "float",
		Metrics: [](mp.Metrics){},
	},
	"xentop.memory": mp.Graphs{
		Label:   "Xentop Memory",
		Unit:    "float",
		Metrics: [](mp.Metrics){},
	},
	"xentop.network": mp.Graphs{
		Label:   "Xentop Network",
		Unit:    "float",
		Metrics: [](mp.Metrics){},
	},
	"xentop.io": mp.Graphs{},
}

type XentopMetrics struct {
	HostName string
	Metrics  mp.Metrics
}

type XentopPlugin struct {
	GraphName          string
	GraphUnit          string
	XentopMetricsSlice []XentopMetrics
}

func parseXentop() error {
	cmd := exec.Command("/bin/sh", "-c", "sudo xentop --batch -i 1 -f")
	s, err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return nil
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
		// TODO: みにくいのでなんとかしたい
		stat[fmt.Sprintf("cpu_%s", name)] = strconv.Atof(sf[index["CPU_PER"]])
		stat[fmt.Sprintf("memory_%s", name)] = strconv.Atof(sf[index["MEM_PER"]])
		stat[fmt.Sprintf("vbd_read_%s", name)] = strconv.Atof(sf[index["VBD_RD"]])
		stat[fmt.Sprintf("vbd_write_%s", name)] = strconv.Atof(sf[index["VBD_WR"]])
		fmt.Println(sf[index["NAME"]])
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return stat
}

// ここでグラフを定義する
func (m XentopPlugin) GraphDefinition() map[string](mp.Graphs) {
	metrics := []mp.Metrics{}
}

func main() {
	// flagの取得

	var xentop XentopPlugin

	helper := mp.NewMackerelPlugin(xentop)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
