//go:build linux

package mpdocker

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	"docker.cpuacct_percentage.#": {
		Label: "Docker CPU Percentage",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "user", Label: "User", Diff: false, Stacked: true, Type: "float64"},
			{Name: "system", Label: "System", Diff: false, Stacked: true, Type: "float64"},
		},
	},
	"docker.memory.#": {
		Label: "Docker Memory",
		Unit:  "bytes",
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
		Unit:  "iops",
		Metrics: []mp.Metrics{
			{Name: "read", Label: "Read", Diff: true, Stacked: true, Type: "uint64", Scale: (1.0 / 60.0)},
			{Name: "write", Label: "Write", Diff: true, Stacked: true, Type: "uint64", Scale: (1.0 / 60.0)},
			{Name: "sync", Label: "Sync", Diff: true, Stacked: true, Type: "uint64", Scale: (1.0 / 60.0)},
			{Name: "async", Label: "Async", Diff: true, Stacked: true, Type: "uint64", Scale: (1.0 / 60.0)},
		},
	},
	"docker.blkio.io_service_bytes.#": {
		Label: "Docker BlkIO Bytes",
		Unit:  "bytes",
		Metrics: []mp.Metrics{
			{Name: "read", Label: "Read", Diff: true, Stacked: true, Type: "uint64"},
			{Name: "write", Label: "Write", Diff: true, Stacked: true, Type: "uint64"},
			{Name: "sync", Label: "Sync", Diff: true, Stacked: true, Type: "uint64"},
			{Name: "async", Label: "Async", Diff: true, Stacked: true, Type: "uint64"},
		},
	},
	// some other fields also exist in metrics, but they're internal intermediate data
}

// DockerPlugin mackerel plugin for docker
type DockerPlugin struct {
	Host             string
	Tempfile         string
	Method           string
	NameFormat       string
	Label            string
	lastMetricValues mp.MetricValues
	UseCPUPercentage bool
}

var normalizeMetricRe = regexp.MustCompile(`[^-a-zA-Z0-9_]`)

func normalizeMetricName(str string) string {
	return normalizeMetricRe.ReplaceAllString(str, "_")
}

func (m DockerPlugin) listContainer() ([]docker.APIContainers, error) {
	client, _ := docker.NewClient(m.Host)
	containers, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// FetchMetrics interface for mackerel plugin
func (m DockerPlugin) FetchMetrics() (map[string]interface{}, error) {
	var stats map[string]interface{}

	if m.Method == "File" {
		return nil, errors.New("no longer supported")
	}
	containers, err := m.listContainer()
	if err != nil {
		return nil, err
	}
	stats, err = m.FetchMetricsWithAPI(containers)
	if err != nil {
		return nil, err
	}

	if m.UseCPUPercentage {
		if time.Since(m.lastMetricValues.Timestamp) <= 5*time.Minute {
			addCPUPercentageStats(&stats, m.lastMetricValues.Values)
		}
	}

	return stats, nil
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
	var mu sync.Mutex
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
			mu.Lock()
			err = m.parseStats(&res, metricName, resultStats[0])
			if err != nil {
				log.Fatal(err)
			}
			mu.Unlock()
		}(container)
	}
	wg.Wait()
	return res, nil
}

const internalCPUStatPrefix = "docker._internal.cpuacct."

func (m DockerPlugin) parseStats(stats *map[string]interface{}, name string, result *docker.Stats) error {
	if m.UseCPUPercentage {
		// intermediate data to calc CPU percentage
		(*stats)[internalCPUStatPrefix+name+".user"] = (*result).CPUStats.CPUUsage.UsageInUsermode
		(*stats)[internalCPUStatPrefix+name+".system"] = (*result).CPUStats.CPUUsage.UsageInKernelmode
		(*stats)[internalCPUStatPrefix+name+".host"] = (*result).CPUStats.SystemCPUUsage

		onlineCPUs := int((*result).CPUStats.OnlineCPUs)
		// if either `CPUStats.OnlineCPUs` or `PerCPUStats.OnlineCPUs` is zero,
		// use the length of CPUUsage.PerCPUUsage for onlineCPUs
		// ref. https://docs.docker.com/engine/api/v1.41/#operation/ContainerStats
		if onlineCPUs == 0 || (*result).PreCPUStats.OnlineCPUs == 0 {
			onlineCPUs = len((*result).CPUStats.CPUUsage.PercpuUsage)
		}
		(*stats)[internalCPUStatPrefix+name+".onlineCPUs"] = onlineCPUs
	} else {
		(*stats)["docker.cpuacct."+name+".user"] = (*result).CPUStats.CPUUsage.UsageInUsermode
		(*stats)["docker.cpuacct."+name+".system"] = (*result).CPUStats.CPUUsage.UsageInKernelmode
	}

	totalRss := (*result).MemoryStats.Stats.TotalRss
	if totalRss == 0 {
		// use `anon` and `file` for RSS and Cache usage on cgroup2 host
		// ref. https://github.com/google/cadvisor/blob/a9858972e75642c2b1914c8d5428e33e6392c08a/container/libcontainer/handler.go#L799-L800
		(*stats)["docker.memory."+name+".rss"] = (*result).MemoryStats.Stats.Anon
		(*stats)["docker.memory."+name+".cache"] = (*result).MemoryStats.Stats.File

	} else {
		// use `total_rss` and `total_cache` for RSS and Cache usage on cgroup host
		(*stats)["docker.memory."+name+".rss"] = totalRss
		(*stats)["docker.memory."+name+".cache"] = (*result).MemoryStats.Stats.TotalCache
	}

	fields := []string{"read", "write", "sync", "async"}
	for _, field := range fields {
		for _, s := range (*result).BlkioStats.IOQueueRecursive {
			if s.Op == cases.Title(language.Und, cases.NoLower).String(field) {
				(*stats)["docker.blkio.io_queued."+name+"."+field] = s.Value
			}
		}
		for _, s := range (*result).BlkioStats.IOServicedRecursive {
			if s.Op == cases.Title(language.Und, cases.NoLower).String(field) {
				(*stats)["docker.blkio.io_serviced."+name+"."+field] = s.Value
			}
		}
		for _, s := range (*result).BlkioStats.IOServiceBytesRecursive {
			if s.Op == cases.Title(language.Und, cases.NoLower).String(field) {
				(*stats)["docker.blkio.io_service_bytes."+name+"."+field] = s.Value
			}
		}
	}
	return nil
}

func addCPUPercentageStats(stats *map[string]interface{}, lastStat map[string]interface{}) {
	for k, v := range lastStat {
		if !strings.HasPrefix(k, internalCPUStatPrefix) || !strings.HasSuffix(k, ".host") {
			continue
		}
		name := strings.TrimSuffix(strings.TrimPrefix(k, internalCPUStatPrefix), ".host")
		currentHostUsage, ok1 := (*stats)[internalCPUStatPrefix+name+".host"]
		cpuNums, ok2 := (*stats)[internalCPUStatPrefix+name+".onlineCPUs"]
		if !ok1 || !ok2 {
			continue
		}
		hostUsage := float64(currentHostUsage.(uint64) - uint64(v.(float64)))
		cpuNumsInt := cpuNums.(int)
		if hostUsage < 0 {
			continue // counter seems reset
		}

		currentUserUsage, ok1 := (*stats)[internalCPUStatPrefix+name+".user"]
		prevUserUsage, ok2 := lastStat[internalCPUStatPrefix+name+".user"]
		if ok1 && ok2 {
			userUsage := float64(currentUserUsage.(uint64) - uint64(prevUserUsage.(float64)))
			if userUsage >= 0 {
				(*stats)["docker.cpuacct_percentage."+name+".user"] = userUsage / hostUsage * 100.0 * float64(cpuNumsInt)
			}
		}

		currentSystemUsage, ok1 := (*stats)[internalCPUStatPrefix+name+".system"]
		prevSystemUsage, ok2 := lastStat[internalCPUStatPrefix+name+".system"]
		if ok1 && ok2 {
			systemUsage := float64(currentSystemUsage.(uint64) - uint64(prevSystemUsage.(float64)))
			if systemUsage >= 0 {
				(*stats)["docker.cpuacct_percentage."+name+".system"] = systemUsage / hostUsage * 100.0 * float64(cpuNumsInt)
			}
		}
	}
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
	flag.String("command", "docker", "Command path to docker(deprecated)") // backward compatibility
	optMethod := flag.String("method", "", "Specify the method to collect stats, 'API' or 'File'. If not specified, an appropriate method is chosen.(deprecated)")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optNameFormat := flag.String("name-format", "name_id", "Set the name format from "+strings.Join(candidateNameFormat, ", "))
	optLabel := flag.String("label", "", "Use the value of the key as name in case that name-format is label.")
	optCPUFormat := flag.String("cpu-format", "", "Specify which CPU metrics format to use, 'percentage' or 'usage'. 'percentage' is default for 'API' method, and is not supported in 'File' method.")
	flag.Parse()

	var docker DockerPlugin

	docker.Host = *optHost
	docker.NameFormat = *optNameFormat
	docker.Label = *optLabel
	if !setCandidateNameFormat[docker.NameFormat] {
		log.Fatalf("Name flag should be each of '%s'", strings.Join(candidateNameFormat, ","))
	}
	if docker.NameFormat == "label" && docker.Label == "" {
		log.Fatalf("Label flag should be set when name flag is 'label'.")
	}

	switch *optMethod {
	case "", "API":
		docker.Method = "API"
	case "File":
		log.Fatalf("'File' method is no longer supported")
	default:
		log.Fatalf("Method should be 'API', 'File' or an empty string.")
	}

	switch *optCPUFormat {
	case "percentage":
		docker.UseCPUPercentage = true
	case "usage":
		docker.UseCPUPercentage = false
	default:
		docker.UseCPUPercentage = true
	}

	helper := mp.NewMackerelPlugin(docker)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.SetTempfileByBasename(fmt.Sprintf("mackerel-plugin-docker-%s", normalizeMetricName(*optHost)))
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		docker.lastMetricValues, _ = helper.FetchLastValues()
		helper.Plugin = docker
		helper.OutputValues()
	}
}
