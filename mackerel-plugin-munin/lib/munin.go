package mpmunin

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

type serviceEnvs map[string]string

type services map[string]serviceEnvs

// MuninMetric metric of munin
type MuninMetric struct {
	Label string
	Type  string
	Draw  string
	Value string
}

// MuninPlugin mackerel plugin for munin
type MuninPlugin struct {
	PluginPath    string
	PluginConfDir string
	GraphTitle    string
	GraphName     string
	MuninMetrics  map[string](*MuninMetric)
}

var exp = map[string](*regexp.Regexp){}

func getExp(expstr string) *(regexp.Regexp) {
	if exp[expstr] == nil {
		exp[expstr] = regexp.MustCompile(expstr)
	}
	return exp[expstr]
}

func getEnvSettingsReader(s *services, plg string, reader io.Reader) {
	svc := ""
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\n")
		line = getExp("([^\\\\])#.*?$").ReplaceAllString(line, "$1")
		line = getExp("\\\\#").ReplaceAllString(line, "#")
		line = getExp("(^\\s+|\\s+$)").ReplaceAllString(line, "")

		svcM := getExp("^\\s*\\[([^\\]]*?)(\\*?)\\]\\s*$").FindStringSubmatch(line)
		if svcM != nil {
			if svcM[2] == "" && plg == svcM[1] { // perfect match
				svc = svcM[1]
			} else if svcM[2] == "*" && len(svcM[1]) <= len(plg) && plg[0:len(svcM[1])] == svcM[1] { // left-hand match
				svc = svcM[1] + svcM[2]
			} else {
				svc = ""
			}
			continue
		}

		if svc == "" {
			continue
		}

		envM := getExp("^\\s*env\\.(\\w+)\\s+(.+)$").FindStringSubmatch(line)
		if envM != nil && svc != "" {
			if (*s)[svc] == nil {
				(*s)[svc] = serviceEnvs{}
			}
			(*s)[svc][envM[1]] = envM[2]
		}
	}
}

func getEnvSettingsFile(s *services, plg string, file string) {
	fp, err := os.Open(file)
	if err != nil {
		return
	}
	defer fp.Close()

	getEnvSettingsReader(s, plg, fp)
}

func compileEnvPairs(s *services, plg string) *map[string]string {
	servenvs := *s

	// ordered services
	srvs := make([]string, 0, len(servenvs))
	for k := range servenvs {
		if k == plg {
			continue
		}
		srvs = append(srvs, k)
	}
	sort.Strings(srvs)
	if servenvs[plg] != nil {
		srvs = append(srvs, plg)
	}

	// apply envs
	envs := make(map[string]string)
	for _, srv := range srvs {
		for k, v := range servenvs[srv] {
			envs[k] = v
		}
	}

	return &envs
}

func setPluginEnvironments(plg string, confdir string) {
	files, err := ioutil.ReadDir(confdir)
	if err != nil {
		log.Fatalln(err)
	}

	filenames := make([]string, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filenames = append(filenames, f.Name())
	}
	sort.Strings(filenames)

	servenvs := make(services)
	for _, f := range filenames {
		getEnvSettingsFile(&servenvs, plg, path.Join(confdir, f))
	}

	for k, v := range *compileEnvPairs(&servenvs, plg) {
		os.Setenv(k, v)
	}
}

func parsePluginConfig(str string, m *map[string](*MuninMetric), title *string) {
	for _, line := range strings.Split(str, "\n") {
		lineM := getExp("^([^ ]+) +(.*)$").FindStringSubmatch(line)
		if lineM == nil {
			continue
		}

		if lineM[1] == "graph_title" {
			(*title) = lineM[2]
			continue
		}

		keyM := getExp("^([^\\.]+)\\.(.*)$").FindStringSubmatch(lineM[1])
		if keyM == nil {
			continue
		}

		if (*m)[keyM[1]] == nil {
			(*m)[keyM[1]] = &(MuninMetric{})
		}
		met := (*m)[keyM[1]]

		switch keyM[2] {
		case "label":
			met.Label = lineM[2]
		case "type":
			met.Type = lineM[2]
		case "draw":
			met.Draw = lineM[2]
		}
	}
}

func parsePluginVals(str string, m *map[string](*MuninMetric)) {
	for _, line := range strings.Split(str, "\n") {
		lineM := getExp("^([^ ]+) +(.*)$").FindStringSubmatch(line)
		if lineM == nil {
			continue
		}

		keyM := getExp("^([^\\.]+)\\.(.*)$").FindStringSubmatch(lineM[1])
		if keyM == nil {
			continue
		}

		if (*m)[keyM[1]] == nil {
			(*m)[keyM[1]] = &(MuninMetric{})
		}
		met := (*m)[keyM[1]]

		switch keyM[2] {
		case "value":
			met.Value = lineM[2]
		}
	}
}

func removeUselessMetrics(m *map[string](*MuninMetric)) {
	// remove metrics which have an empty Value
	for name, mmet := range *m {
		if mmet.Value == "" {
			delete(*m, name)
		}
	}
}

func (p *MuninPlugin) prepare() error {
	var err error

	if p.PluginConfDir != "" {
		setPluginEnvironments(path.Base(p.PluginPath), p.PluginConfDir)
	}

	outConfig, err := exec.Command(p.PluginPath, "config").Output()
	if err != nil {
		return fmt.Errorf("%s: %s", err, outConfig)
	}

	outVals, err := exec.Command(p.PluginPath).Output()
	if err != nil {
		return fmt.Errorf("%s: %s", err, outVals)
	}

	p.MuninMetrics = make(map[string](*MuninMetric))
	parsePluginConfig(string(outConfig), &p.MuninMetrics, &p.GraphTitle)
	parsePluginVals(string(outVals), &p.MuninMetrics)
	removeUselessMetrics(&p.MuninMetrics)

	return nil
}

// FetchMetrics interface for mackerelplugin
func (p MuninPlugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64, len(p.MuninMetrics))
	for name, mmet := range p.MuninMetrics {
		parsed, err := strconv.ParseFloat(mmet.Value, 64)
		if err != nil {
			log.Printf("Failed to parse value of %s: %s", name, err)
			continue
		}

		stat[name] = parsed
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (p MuninPlugin) GraphDefinition() map[string]mp.Graphs {
	metrics := make([]mp.Metrics, 0, len(p.MuninMetrics))
	for name, mmet := range p.MuninMetrics {
		met := mp.Metrics{Name: name}
		if mmet.Label == "" {
			met.Label = name
		} else {
			met.Label = mmet.Label
		}
		if mmet.Draw == "STACK" {
			met.Stacked = true
		}
		switch mmet.Type {
		case "COUNTER", "DERIVE", "ABSOLUTE":
			met.Diff = true
		}

		metrics = append(metrics, met)
	}

	return map[string]mp.Graphs{p.GraphName: {
		Label:   p.GraphTitle,
		Unit:    "float",
		Metrics: metrics,
	}}
}

// Do the plugin
func Do() {
	optPluginPath := flag.String("plugin", "", "Munin plugin path")
	optPluginConfDir := flag.String("plugin-conf-d", "", "Munin plugin-conf.d path")
	optGraphName := flag.String("name", "", "Graph name")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var munin MuninPlugin
	if *optPluginPath == "" {
		log.Fatalln("Munin plugin path is required")
	}
	munin.PluginPath = *optPluginPath
	munin.PluginConfDir = *optPluginConfDir
	if *optGraphName == "" {
		munin.GraphName = "munin." + path.Base(munin.PluginPath)
	} else {
		munin.GraphName = *optGraphName
	}

	err := munin.prepare()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(munin)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-munin-%s", munin.GraphName)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
