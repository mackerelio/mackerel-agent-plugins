package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin"
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
)

var exp map[string](*regexp.Regexp) = map[string](*regexp.Regexp){}

type ServiceEnvs map[string]string

type Services map[string]ServiceEnvs

type MuninMetric struct {
	Label string
	Type  string
	Draw  string
	Value string
}

type MuninPlugin struct {
	PluginPath    string
	PluginConfDir string
	GraphTitle    string
	GraphName     string
	MuninMetrics  map[string](*MuninMetric)
}

func getExp(expstr string) *(regexp.Regexp) {
	if exp[expstr] == nil {
		exp[expstr] = regexp.MustCompile(expstr)
	}
	return exp[expstr]
}

func getEnvSettingsReader(s *Services, plg string, reader io.Reader) {
	var service string
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\n")
		line = getExp("([^\\\\])#.*?$").ReplaceAllString(line, "$1")
		line = getExp("\\\\#").ReplaceAllString(line, "#")
		line = getExp("(^\\s+|\\s+$)").ReplaceAllString(line, "")

		service_m := getExp("^\\s*\\[([^\\]]*?)(\\*?)\\]\\s*$").FindStringSubmatch(line)
		if service_m != nil {
			if service_m[2] == "" && plg == service_m[1] { // perfect match
				service = service_m[1]
			} else if service_m[2] == "*" && len(service_m[1]) <= len(plg) && plg[0:len(service_m[1])] == service_m[1] { // left-hand match
				service = service_m[1] + service_m[2]
			} else {
				service = ""
			}
			continue
		}

		if service == "" {
			continue
		}

		env_m := getExp("^\\s*env\\.(\\w+)\\s+(.+)$").FindStringSubmatch(line)
		if env_m != nil && service != "" {
			if (*s)[service] == nil {
				(*s)[service] = ServiceEnvs{}
			}
			(*s)[service][env_m[1]] = env_m[2]
		}
	}
}

func getEnvSettingsFile(s *Services, plg string, file string) {
	fp, err := os.Open(file)
	if err != nil {
		return
	}
	defer fp.Close()

	getEnvSettingsReader(s, plg, fp)
}

func compileEnvPairs(s *Services, plg string) *map[string]string {
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

	servenvs := make(Services)
	for _, f := range filenames {
		getEnvSettingsFile(&servenvs, plg, path.Join(confdir, f))
	}

	for k, v := range *compileEnvPairs(&servenvs, plg) {
		os.Setenv(k, v)
	}
}

func parsePluginConfig(str string, m *map[string](*MuninMetric), title *string) {
	for _, line := range strings.Split(str, "\n") {
		line_m := getExp("^([^ ]+) +(.*)$").FindStringSubmatch(line)
		if line_m == nil {
			continue
		}

		if line_m[1] == "graph_title" {
			(*title) = line_m[2]
			continue
		}

		key_m := getExp("^([^\\.]+)\\.(.*)$").FindStringSubmatch(line_m[1])
		if key_m == nil {
			continue
		}

		if (*m)[key_m[1]] == nil {
			(*m)[key_m[1]] = &(MuninMetric{})
		}
		met := (*m)[key_m[1]]

		switch key_m[2] {
		case "label":
			met.Label = line_m[2]
		case "type":
			met.Type = line_m[2]
		case "draw":
			met.Draw = line_m[2]
		}
	}
}

func parsePluginVals(str string, m *map[string](*MuninMetric)) {
	for _, line := range strings.Split(str, "\n") {
		line_m := getExp("^([^ ]+) +(.*)$").FindStringSubmatch(line)
		if line_m == nil {
			continue
		}

		key_m := getExp("^([^\\.]+)\\.(.*)$").FindStringSubmatch(line_m[1])
		if key_m == nil {
			continue
		}

		if (*m)[key_m[1]] == nil {
			(*m)[key_m[1]] = &(MuninMetric{})
		}
		met := (*m)[key_m[1]]

		switch key_m[2] {
		case "value":
			met.Value = line_m[2]
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

func (p *MuninPlugin) Prepare() error {
	var err error

	if p.PluginConfDir != "" {
		setPluginEnvironments(path.Base(p.PluginPath), p.PluginConfDir)
	}

	out_config, err := exec.Command(p.PluginPath, "config").Output()
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %s", err, out_config))
	}

	out_vals, err := exec.Command(p.PluginPath).Output()
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %s", err, out_vals))
	}

	p.MuninMetrics = make(map[string](*MuninMetric))
	parsePluginConfig(string(out_config), &p.MuninMetrics, &p.GraphTitle)
	parsePluginVals(string(out_vals), &p.MuninMetrics)
	removeUselessMetrics(&p.MuninMetrics)

	return nil
}

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

func (p MuninPlugin) GraphDefinition() map[string](mp.Graphs) {
	metrics := make([](mp.Metrics), 0, len(p.MuninMetrics))
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

	return map[string](mp.Graphs){p.GraphName: mp.Graphs{
		Label:   p.GraphTitle,
		Unit:    "float",
		Metrics: metrics,
	}}
}

func main() {
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

	err := munin.Prepare()
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
