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

	"github.com/fsouza/go-dockerclient"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string](mp.Graphs){
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

// DockerPlugin mackerel plugin for docker
type DockerPlugin struct {
	Host          string
	DockerCommand string
	Tempfile      string
	Method        string
	pathBuilder   *pathBuilder
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

func (m DockerPlugin) listContainer() ([]docker.APIContainers, error) {
	client, _ := docker.NewClient(m.Host)
	containers, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func findPrefixPath() (string, error) {
	pathCandidate := []string{"/host/cgroup", "/cgroup", "/host/sys/fs/cgroup", "/sys/fs/cgroup"}
	for _, path := range pathCandidate {
		result, err := exists(path)
		if err != nil {
			return "", err
		}
		resultDeep, err := exists(path + "/cpuacct")
		if err != nil {
			return "", err
		}
		if result && resultDeep {
			return path, nil
		}
	}
	return "", errors.New("No prefix path is found.")
}

type pathBuilder struct {
	prefix   string
	pathType pathType
}

type pathType uint8

const (
	pathUnknown pathType = iota
	pathDocker
	pathLxc
	pathSlice
)

func newPathBuilder() (*pathBuilder, error) {
	prefixPath, err := findPrefixPath()
	if err != nil {
		return nil, err
	}
	pathT, err := guessPathType(prefixPath)
	if err != nil {
		return nil, err
	}
	return &pathBuilder{
		prefix:   prefixPath,
		pathType: pathT,
	}, nil
}

func (pb *pathBuilder) build(id, metric, postfix string) string {
	switch pb.pathType {
	case pathDocker:
		return fmt.Sprintf("%s/%s/docker/%s/%s.%s", pb.prefix, metric, id, metric, postfix)
	case pathLxc:
		return fmt.Sprintf("%s/%s/lxc/%s/%s.%s", pb.prefix, metric, id, metric, postfix)
	case pathSlice:
		return fmt.Sprintf("%s/%s/system.slice/docker-%s.scope/%s.%s", pb.prefix, metric, id, metric, postfix)
	default:
		return ""
	}
}

func guessPathType(prefix string) (pathType, error) {
	if ok, err := exists(prefix + "/memory/system.slice/"); ok && err == nil {
		return pathSlice, nil
	}
	if ok, err := exists(prefix + "/memory/docker/"); ok && err == nil {
		return pathDocker, nil
	}
	if ok, err := exists(prefix + "/memory/lxc/"); ok && err == nil {
		return pathLxc, nil
	}
	return pathUnknown, fmt.Errorf("can't resolve runtime metrics path")
}

func guessMethod(docker string) (string, error) {
	cmd := exec.Command(docker, "version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	lines := strings.SplitAfterN(out.String(), "\n", 2)
	// Client version: 1.6.2
	re := regexp.MustCompile(`Client version: ([0-9]+)(?:\.([0-9]+))?(?:\.([0-9]+))?`)
	res := re.FindAllStringSubmatch(lines[0], 1)
	if len(res) < 1 || len(res[0]) < 2 {
		log.Printf("Use API because of failing to recognize version")
		return "API", nil
	}

	majorVer, err := strconv.Atoi(res[0][1])
	if err != nil {
		log.Printf("Use API because of failing to recognize version")
		return "API", nil
	}

	minorVer, err := strconv.Atoi(res[0][2])
	if err != nil {
		log.Printf("Use API because of failing to recognize version")
		return "API", nil
	}

	if majorVer == 1 && minorVer < 9 {
		return "File", nil
	}
	return "API", nil
}

// FetchMetrics interface for mackerel plugin
func (m DockerPlugin) FetchMetrics() (map[string]interface{}, error) {
	dockerStats := map[string][]string{}
	if m.Method == "API" {
		containers, err := m.listContainer()
		if err != nil {
			return nil, err
		}
		return m.FetchMetricsWithAPI(containers)
	}

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
		if len(fields) > 3 {
			dockerStats[fields[0]] = []string{fields[1], fields[len(fields)-2]}
		}
	}

	return m.FetchMetricsWithFile(&dockerStats)
}

// FetchMetricsWithAPI use docker API to fetch metrics
func (m DockerPlugin) FetchMetricsWithAPI(containers []docker.APIContainers) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	for _, container := range containers {
		name := strings.Replace(container.Names[0], "/", "", 1)
		client, _ := docker.NewClient(m.Host)
		errC := make(chan error, 1)
		statsC := make(chan *docker.Stats)
		done := make(chan bool)
		go func() {
			errC <- client.Stats(docker.StatsOptions{name, statsC, false, done, 0})
			close(errC)
		}()
		var resultStats []*docker.Stats
		for {
			stats, ok := <-statsC
			if !ok {
				break
			}
			resultStats = append(resultStats, stats)
		}
		err := <-errC
		if err != nil {
			log.Fatal(err)
		}
		if len(resultStats) == 0 {
			log.Fatalf("Stats: Expected 1 result. Got %d.", len(resultStats))
		}
		m.parseStats(&res, name, resultStats[0])
	}
	return res, nil
}

func (m DockerPlugin) parseStats(stats *map[string]interface{}, name string, result *docker.Stats) error {
	(*stats)["docker.cpuacct."+name+".user"] = (*result).CPUStats.CPUUsage.UsageInUsermode
	(*stats)["docker.cpuacct."+name+".system"] = (*result).CPUStats.CPUUsage.UsageInKernelmode
	(*stats)["docker.memory."+name+".cache"] = (*result).MemoryStats.Stats.TotalCache
	(*stats)["docker.memory."+name+".rss"] = (*result).MemoryStats.Stats.TotalRss
	fields := []string{"read", "write", "sync", "async"}
	for _, field := range fields {
		for _, s := range (*result).BlkioStats.IOQueueRecursive {
			if s.Op == strings.Title(field) {
				(*stats)["docker.blkio.io_queued."+name+"."+field] = s.Value
			}
		}
		for _, s := range (*result).BlkioStats.IOServicedRecursive {
			if s.Op == strings.Title(field) {
				(*stats)["docker.blkio.io_serviced."+name+"."+field] = s.Value
			}
		}
		for _, s := range (*result).BlkioStats.IOServiceBytesRecursive {
			if s.Op == strings.Title(field) {
				(*stats)["docker.blkio.io_service_bytes."+name+"."+field] = s.Value
			}
		}
	}
	return nil
}

// FetchMetricsWithFile use cgroup stats files to fetch metrics
func (m DockerPlugin) FetchMetricsWithFile(dockerStats *map[string][]string) (map[string]interface{}, error) {
	pb := m.pathBuilder

	metrics := map[string][]string{
		"cpuacct": []string{"user", "system"},
		"memory":  []string{"cache", "rss"},
	}

	res := map[string]interface{}{}
	for id, name := range *dockerStats {
		for metric, stats := range metrics {
			if ok, err := exists(pb.build(id, metric, "stat")); !ok || err != nil {
				continue
			}
			data, err := getFile(pb.build(id, metric, "stat"))
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
			if ok, err := exists(pb.build(id, "blkio", blkioType)); !ok || err != nil {
				continue
			}
			data, err := getFile(pb.build(id, "blkio", blkioType))
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
						v += ret
					}
				}
				res[fmt.Sprintf("docker.blkio.%s.%s_%s.%s", blkioType, normalizeMetricName(name[0]), id[0:6], strings.ToLower(stat))] = v
			}
		}

	}

	return res, nil
}

// GraphDefinition interface for mackerel plugin
func (m DockerPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optHost := flag.String("host", "unix:///var/run/docker.sock", "Host for socket")
	optCommand := flag.String("command", "docker", "Command path to docker")
	optUseAPI := flag.String("method", "", "Specify the method to collect stats, 'API' or 'File'. If not specified, an appropriate method is choosen.")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var docker DockerPlugin

	docker.Host = fmt.Sprintf("%s", *optHost)
	docker.DockerCommand = *optCommand
	_, err := exec.LookPath(docker.DockerCommand)
	if err != nil {
		log.Fatalf("Docker command is not found: %s", docker.DockerCommand)
	}

	if *optUseAPI == "" {
		docker.Method, err = guessMethod(docker.DockerCommand)
		if err != nil {
			log.Fatalf("Fail to guess stats method: %s", err.Error())
		}
	} else {
		if *optUseAPI != "API" && *optUseAPI != "File" {
			log.Fatalf("Method should be 'API', 'File' or an empty string.")
		}
		docker.Method = *optUseAPI
	}

	if docker.Method == "File" {
		pb, err := newPathBuilder()
		if err != nil {
			log.Fatalf("failed to resolve docker metrics path: %s. It may be no Docker containers exists.", err)
		}
		docker.pathBuilder = pb
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
