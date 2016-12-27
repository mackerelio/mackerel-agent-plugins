package mpdocker

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
	"sync"
	"time"

	"github.com/fsouza/go-dockerclient"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string]mp.Graphs{
	"docker.cpuacct.#": {
		Label: "Docker CPU",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "user", Label: "User", Diff: true, Stacked: true, Type: "uint64"},
			{Name: "system", Label: "System", Diff: true, Stacked: true, Type: "uint64"},
		},
	},
	"docker.memory.#": {
		Label: "Docker Memory",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "cache", Label: "Cache", Diff: false, Stacked: true},
			{Name: "rss", Label: "RSS", Diff: false, Stacked: true},
		},
	},
	"docker.blkio.io_queued.#": {
		Label: "Docker BlkIO Queued",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "read", Label: "Read", Diff: false, Stacked: true},
			{Name: "write", Label: "Write", Diff: false, Stacked: true},
			{Name: "sync", Label: "Sync", Diff: false, Stacked: true},
			{Name: "async", Label: "Async", Diff: false, Stacked: true},
		},
	},
	"docker.blkio.io_serviced.#": {
		Label: "Docker BlkIO IOPS",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "read", Label: "Read", Diff: true, Stacked: true, Type: "uint64"},
			{Name: "write", Label: "Write", Diff: true, Stacked: true, Type: "uint64"},
			{Name: "sync", Label: "Sync", Diff: true, Stacked: true, Type: "uint64"},
			{Name: "async", Label: "Async", Diff: true, Stacked: true, Type: "uint64"},
		},
	},
	"docker.blkio.io_service_bytes.#": {
		Label: "Docker BlkIO Bytes",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "read", Label: "Read", Diff: true, Stacked: true, Type: "uint64"},
			{Name: "write", Label: "Write", Diff: true, Stacked: true, Type: "uint64"},
			{Name: "sync", Label: "Sync", Diff: true, Stacked: true, Type: "uint64"},
			{Name: "async", Label: "Async", Diff: true, Stacked: true, Type: "uint64"},
		},
	},
}

// DockerPlugin mackerel plugin for docker
type DockerPlugin struct {
	Host          string
	DockerCommand string
	Tempfile      string
	Method        string
	NameFormat    string
	Label         string
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
		if result {
			return path, nil
		}
	}
	return "", errors.New("no prefix path is found")
}

type pathBuilder struct {
	prefix   string
	pathType pathType
}

type pathType uint8

const (
	pathUnknown pathType = iota
	pathDocker
	pathDockerShort
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
	case pathDockerShort:
		return fmt.Sprintf("%s/docker/%s/%s.%s", pb.prefix, id, metric, postfix)
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
	if ok, err := exists(prefix + "/docker/"); ok && err == nil {
		return pathDockerShort, nil
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

	re := regexp.MustCompile(`Server API version: ([0-9]+)(?:\.([0-9]+))?`)
	res := re.FindAllStringSubmatch(out.String(), 1)
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

	if majorVer == 1 && minorVer < 17 {
		//log.Printf("Use File for fetching Metrics")
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

func (m DockerPlugin) generateName(container docker.APIContainers) string {
	switch m.NameFormat {
	case "name_id":
		return fmt.Sprintf("%s_%s", strings.Replace(container.Names[0], "/", "", 1), container.ID[0:6])
	case "name":
		return strings.Replace(container.Names[0], "/", "", 1)
	case "id":
		return container.ID
	case "image":
		return container.Image
	case "image_id":
		return fmt.Sprintf("%s_%s", container.Image, container.ID[0:6])
	case "image_name":
		return fmt.Sprintf("%s_%s", container.Image, strings.Replace(container.Names[0], "/", "", 1))
	case "label":
		return container.Labels[m.Label]
	}
	return strings.Replace(container.Names[0], "/", "", 1)
}

// FetchMetricsWithAPI use docker API to fetch metrics
func (m DockerPlugin) FetchMetricsWithAPI(containers []docker.APIContainers) (map[string]interface{}, error) {
	var wg sync.WaitGroup
	res := map[string]interface{}{}
	for _, container := range containers {
		wg.Add(1)
		go func(cont docker.APIContainers) {
			defer wg.Done()
			name := strings.Replace(cont.Names[0], "/", "", 1)
			metricName := normalizeMetricName(m.generateName(cont))
			client, _ := docker.NewClient(m.Host)
			errC := make(chan error, 1)
			statsC := make(chan *docker.Stats)
			done := make(chan bool)
			go func() {
				errC <- client.Stats(docker.StatsOptions{ID: name, Stats: statsC, Stream: false, Done: done, Timeout: time.Duration(20) * time.Second})
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
			m.parseStats(&res, metricName, resultStats[0])
		}(container)
	}
	wg.Wait()
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
		"cpuacct": {"user", "system"},
		"memory":  {"cache", "rss"},
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
func (m DockerPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	candidateNameFormat := []string{"name", "name_id", "id", "image", "image_id", "image_name", "label"}
	setCandidateNameFormat := make(map[string]bool)
	for _, v := range candidateNameFormat {
		setCandidateNameFormat[v] = true
	}

	optHost := flag.String("host", "unix:///var/run/docker.sock", "Host for socket")
	optCommand := flag.String("command", "docker", "Command path to docker")
	optUseAPI := flag.String("method", "", "Specify the method to collect stats, 'API' or 'File'. If not specified, an appropriate method is chosen.")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optNameFormat := flag.String("name-format", "name_id", "Set the name format from "+strings.Join(candidateNameFormat, ", "))
	optLabel := flag.String("label", "", "Use the value of the key as name in case that name-format is label.")
	flag.Parse()

	var docker DockerPlugin

	docker.Host = fmt.Sprintf("%s", *optHost)
	docker.DockerCommand = *optCommand
	_, err := exec.LookPath(docker.DockerCommand)
	if err != nil {
		log.Fatalf("Docker command is not found: %s", docker.DockerCommand)
	}

	docker.NameFormat = *optNameFormat
	docker.Label = *optLabel
	if !setCandidateNameFormat[docker.NameFormat] {
		log.Fatalf("Name flag should be each of '%s'", strings.Join(candidateNameFormat, ","))
	}
	if docker.NameFormat == "label" && docker.Label == "" {
		log.Fatalf("Label flag should be set when name flag is 'label'.")
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
