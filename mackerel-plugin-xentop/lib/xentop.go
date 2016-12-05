package mpxentop

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// index from table headers to array index
// This is generated dynamically at generateIndex
var index = map[string]int{}

var graphdef = map[string]mp.Graphs{
	"xentop.cpu.#": {
		Label: "Xentop CPU",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "cpu", Label: "cpu", Stacked: true, Diff: true},
		},
	},
	"xentop.memory.#": {
		Label: "Xentop Memory",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "memory", Label: "memory", Stacked: true},
		},
	},
	"xentop.nettx.#": {
		Label: "Xentop Nettx",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "nettx", Label: "nettx", Stacked: true, Diff: true},
		},
	},
	"xentop.netrx.#": {
		Label: "Xentop Netrx",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "netrx", Label: "netrx", Stacked: true, Diff: true},
		},
	},
	"xentop.vbdrd.#": {
		Label: "Xentop VBD_RD",
		Unit:  "iops",
		Metrics: []mp.Metrics{
			{Name: "vbdrd", Label: "vbdrd", Stacked: true, Diff: true},
		},
	},
	"xentop.vbdwr.#": {
		Label: "Xentop VBD_WR",
		Unit:  "iops",
		Metrics: []mp.Metrics{
			{Name: "vbdwr", Label: "vbdwr", Stacked: true, Diff: true},
		},
	},
}

// XentopPlugin mackerel plugin for xentop
type XentopPlugin struct {
	XenVersion int
}

// FetchMetrics interface for mackerelplugin
func (m XentopPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	var cmd *exec.Cmd
	if m.XenVersion == 4 {
		cmd = exec.Command("/bin/sh", "-c", "xentop --batch -i 1 -f")
	} else {
		cmd = exec.Command("/bin/sh", "-c", "xentop --batch -i 1")
	}
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	execErr := cmd.Start()
	if execErr != nil {
		fmt.Println(execErr)
		os.Exit(1)
	}

	dom0 := false
	scanner := bufio.NewScanner(stdout)
	hasIndex := false
	for scanner.Scan() {
		sf := strings.Fields(string(scanner.Text()))
		if sf[0] == "NAME" {
			generateIndex(sf, index)
			hasIndex = true
			continue
		}
		if !hasIndex {
			continue
		}
		if stringInSlice("n/a", sf) {
			changeIndex(&index)
			dom0 = true
		}
		name := normalizeXenName(sf[index["NAME"]])

		var errParse error
		var tmpval float64 // avoid `stat[*] *= 1000` because of go interface bug

		tmpval, errParse = strconv.ParseFloat(sf[index["CPU(sec)"]], 64)
		if errParse != nil {
			return nil, errParse
		}
		stat[fmt.Sprintf("xentop.cpu.%s.cpu", name)] = tmpval * 100 / 60
		stat[fmt.Sprintf("xentop.memory.%s.memory", name)], errParse = strconv.ParseFloat(sf[index["MEM(%)"]], 64)
		if errParse != nil {
			return nil, errParse
		}
		tmpval, errParse = strconv.ParseFloat(sf[index["NETTX(k)"]], 64)
		if errParse != nil {
			return nil, errParse
		}
		stat[fmt.Sprintf("xentop.nettx.%s.nettx", name)] = tmpval * 1000
		tmpval, errParse = strconv.ParseFloat(sf[index["NETRX(k)"]], 64)
		if errParse != nil {
			return nil, errParse
		}
		stat[fmt.Sprintf("xentop.netrx.%s.netrx", name)] = tmpval * 1000
		stat[fmt.Sprintf("xentop.vbdrd.%s.vbdrd", name)], errParse = strconv.ParseFloat(sf[index["VBD_RD"]], 64)
		if errParse != nil {
			return nil, errParse
		}
		stat[fmt.Sprintf("xentop.vbdwr.%s.vbdwr", name)], errParse = strconv.ParseFloat(sf[index["VBD_WR"]], 64)
		if errParse != nil {
			return nil, errParse
		}
		if dom0 {
			revertIndex(&index)
			dom0 = false
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (m XentopPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optXenVersion := flag.Int("xenversion", 4, "Xen Version")
	flag.Parse()

	var xentop XentopPlugin

	xentop.XenVersion = *optXenVersion

	helper := mp.NewMackerelPlugin(xentop)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}

func normalizeXenName(raw string) string {
	return strings.Replace(raw, ".", "_", -1)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func generateIndex(sf []string, index map[string]int) {
	i := 0
	for _, column := range sf {
		index[column] = i
		i++
	}
}

func changeIndex(p *map[string]int) {
	maxmemPer := (*p)["MAXMEM(%)"]
	for key, value := range *p {
		if value >= maxmemPer {
			(*p)[key]++
		}
	}
}

func revertIndex(p *map[string]int) {
	maxmemPer := (*p)["MAXMEM(%)"]
	for key, value := range *p {
		if value >= maxmemPer {
			(*p)[key]--
		}
	}
}
