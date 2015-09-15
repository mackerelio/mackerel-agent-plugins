package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
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
	"docker.blkio.io_queued.#": mp.Graphs{
		Label: "Docker BlkIO Queued",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "read", Label: "Read", Diff: false, Stacked: true},
			mp.Metrics{Name: "write", Label: "Write", Diff: false, Stacked: true},
			mp.Metrics{Name: "sync", Label: "Sync", Diff: false, Stacked: true},
			mp.Metrics{Name: "async", Label: "Async", Diff: false, Stacked: true},
		},
	},
	"docker.blkio.io_serviced.#": mp.Graphs{
		Label: "Docker BlkIO IOPS",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "read", Label: "Read", Diff: true, Stacked: true, Type: "uint64"},
			mp.Metrics{Name: "write", Label: "Write", Diff: true, Stacked: true, Type: "uint64"},
			mp.Metrics{Name: "sync", Label: "Sync", Diff: true, Stacked: true, Type: "uint64"},
			mp.Metrics{Name: "async", Label: "Async", Diff: true, Stacked: true, Type: "uint64"},
		},
	},
	"docker.blkio.io_service_bytes.#": mp.Graphs{
		Label: "Docker BlkIO Bytes",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "read", Label: "Read", Diff: true, Stacked: true, Type: "uint64"},
			mp.Metrics{Name: "write", Label: "Write", Diff: true, Stacked: true, Type: "uint64"},
			mp.Metrics{Name: "sync", Label: "Sync", Diff: true, Stacked: true, Type: "uint64"},
			mp.Metrics{Name: "async", Label: "Async", Diff: true, Stacked: true, Type: "uint64"},
		},
	},
}

type DockerPlugin struct {
	Host          string
	DockerCommand string
	Tempfile      string
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

var normalizeMetricRe = regexp.MustCompile(`[^-a-zA-Z0-9_]`)

func normalizeMetricName(str string) string {
	return normalizeMetricRe.ReplaceAllString(str, "_")
}

func (m DockerPlugin) getDockerPs() (string, error) {
	cmd := exec.Command(m.DockerCommand, "--host", m.Host, "ps", "--no-trunc")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
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
	data, err := m.getDockerPs()
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

		// blkio statistics
		for _, blkioType := range []string{"io_queued", "io_serviced", "io_service_bytes"} {
			data, err := getFile(fmt.Sprintf("%s/blkio/docker/%s/blkio.%s", prefixPath, id, blkioType))
			if err != nil {
				return nil, err
			}
			for _, stat := range []string{"Read", "Write", "Sync", "Async"} {
				re := regexp.MustCompile(stat + " (\\d+)")
				matchs := re.FindAllStringSubmatch(data, -1)
				v := 0.0
				for _, m := range matchs {
					if m != nil {
						ret, _ := strconv.ParseFloat(m[1], 64)
						//fmt.Println(ret)
						v += ret
					}
				}
				res[fmt.Sprintf("docker.blkio.%s.%s_%s.%s", blkioType, normalizeMetricName(name[0]), id[0:6], strings.ToLower(stat))] = v
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
	optHost := flag.String("host", "unix:///var/run/docker.sock", "Host for socket")
	optCommand := flag.String("command", "docker", "Command path to docker")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var docker DockerPlugin

	docker.Host = fmt.Sprintf("%s", *optHost)
	docker.DockerCommand = *optCommand
	_, err := exec.LookPath(docker.DockerCommand)
	if err != nil {
		log.Fatalf("Docker command is not found: %s", docker.DockerCommand)
	}

	helper := mp.NewMackerelPlugin(docker)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-docker-%s", normalizeMetricName(*optHost))
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
