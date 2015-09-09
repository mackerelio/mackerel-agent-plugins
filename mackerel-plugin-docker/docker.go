package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"docker.cpuacct.#": mp.Graphs{
		Label: "Docker CPU",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "user", Label: "User", Diff: true, Stacked: true, Type: "uint64"},
			mp.Metrics{Name: "system", Label: "System", Diff: true, Stacked: true, Type: "uint64"},
		},
	},
	"docker.memory.#": mp.Graphs{
		Label: "Docker Memory",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "cache", Label: "Cache", Diff: false, Stacked: true},
			mp.Metrics{Name: "rss", Label: "RSS", Diff: false, Stacked: true},
		},
	},
}

type DockerPlugin struct {
	Target   string
	Tempfile string
}

func getDockerPs() (string, error) {
	cmd := exec.Command("docker", "ps", "--no-trunc")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func getFile(path string) (string, error) {
	cmd := exec.Command("cat", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func normalizeMetricName(str string) string {
	re := regexp.MustCompile("[-a-zA-Z0-9_]")
	runes := []rune(str)
	for i, c := range str {
		if re.FindString(string(c)) == "" {
			runes[i] = '_'
		}
	}
	return string(runes)
}

func (m DockerPlugin) findPrefixPath() (string, error) {
	pathCandidate := []string{"/host/sys/fs/cgroup", "/sys/fs/cgroup"}
	for _, path := range pathCandidate {
		result, err := exists(path)
		if err != nil {
			return "", err
		}
		if result {
			return path, nil
		}
	}
	return "", errors.New("No prefix path is found.")
}

func (m DockerPlugin) FetchMetrics() (map[string]interface{}, error) {
	dockerStats := map[string][]string{}
	data, err := getDockerPs()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(data, "\n")
	for n, line := range lines {
		if n == 0 {
			continue
		}
		fields := regexp.MustCompile(" +").Split(line, -1)
		//fmt.Println(fields)
		if len(fields) > 3 {
			dockerStats[fields[0]] = []string{fields[1], fields[len(fields)-2]}
		}
	}

	prefixPath, err := m.findPrefixPath()
	if err != nil {
		return nil, err
	}

	metrics := map[string][]string{
		"cpuacct": []string{"user", "system"},
		"memory":  []string{"cache", "rss"},
	}

	res := map[string]interface{}{}
	for id, name := range dockerStats {
		for metric, stats := range metrics {
			data, err := getFile(fmt.Sprintf("%s/%s/docker/%s/%s.stat", prefixPath, metric, id, metric))
			if err != nil {
				return nil, err
			}
			for _, stat := range stats {
				re := regexp.MustCompile(stat + " (\\d+)")
				m := re.FindStringSubmatch(data)
				if m != nil {
					res[fmt.Sprintf("docker.%s.%s_%s.%s", metric, normalizeMetricName(name[0]), id[0:6], stat)] = m[1]
				}
			}
		}
	}
	//fmt.Println(res)

	return res, nil
}

func (m DockerPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "4243", "Port")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var docker DockerPlugin

	docker.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	helper := mp.NewMackerelPlugin(docker)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-docker-%s-%s", *optHost, *optPort)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
