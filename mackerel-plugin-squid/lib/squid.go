package mpsquid

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string]mp.Graphs{
	"squid.requests": {
		Label: "Squid Client Requests",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "requests", Label: "Requests", Diff: true},
		},
	},
	"squid.cache_hit_ratio.5min": {
		Label: "Squid Client Cache Hit Ratio (5min)",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "request_ratio", Label: "Request Ratio", Diff: false},
			{Name: "byte_ratio", Label: "Byte Ratio", Diff: false},
		},
	},
	"squid.cpu_usage_ratio.5min": {
		Label: "Squid CPU Usage Ratio (5min)",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "cpu_usage", Label: "CPU Usage Ratio", Diff: false},
		},
	},
	"squid.cache_storage_usage": {
		Label: "Squid Cache Storage Usage",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: "swap_used_ratio", Label: "Swap capacity (used)", Diff: false},
			{Name: "memory_used_ratio", Label: "Memory capacity (used)", Diff: false},
		},
	},
	"squid.file_descriptor_usage": {
		Label: "Squid File descriptor usage",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "total_fd", Label: "Maximum number of file descriptors", Diff: false},
			{Name: "max_fd", Label: "Largest file desc currently in use", Diff: false},
			{Name: "current_fd", Label: "Number of file desc currently in use", Diff: false},
			{Name: "avail_fd", Label: "Available number of file descriptors", Diff: false},
			{Name: "reserved_fd", Label: "Reserved number of file descriptors", Diff: false},
			{Name: "open_files", Label: "Store Disk files open", Diff: false},
			{Name: "queued_files", Label: "Files queued for open", Diff: false},
		},
	},
	"squid.memory_account_for": {
		Label: "Squid Memory accounted for",
		Unit:  "interger",
		Metrics: []mp.Metrics{
			{Name: "memory_poll_alloc", Label: "memPoolAlloc calls", Diff: true},
			{Name: "memory_poll_free", Label: "memPoolFree calls", Diff: true},
		},
	},
}

// SquidPlugin mackerel plugin for squid
type SquidPlugin struct {
	Target   string
	Tempfile string
}

// FetchMetrics interface for mackerelplugin
func (m SquidPlugin) FetchMetrics() (map[string]interface{}, error) {
	conn, err := net.Dial("tcp", m.Target)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(conn, "GET cache_object://"+m.Target+"/info HTTP/1.0\n\n")
	return m.ParseMgrInfo(conn)
}

// ParseMgrInfo parser for squid mgr:info
func (m SquidPlugin) ParseMgrInfo(info io.Reader) (map[string]interface{}, error) {
	// https://wiki.squid-cache.org/Features/CacheManager/Info?highlight=%28Feature..Squid.Cache.Manager%29
	scanner := bufio.NewScanner(info)

	stat := make(map[string]interface{})
	//regexpmap := make(map[string]*regexp.Regexp)
	regexpmap := map[*regexp.Regexp]string{
		regexp.MustCompile("Number of HTTP requests received:\t([0-9]+)"): "requests",
		// version 2
		regexp.MustCompile("Request Hit Ratios:\t5min: ([0-9\\.]+)%"): "request_ratio",
		regexp.MustCompile("Byte Hit Ratios:\t5min: ([0-9\\.]+)%"):    "byte_ratio",
		// version 3
		regexp.MustCompile("Hits as % of all requests:\t5min: ([0-9\\.]+)%"):      "request_ratio",
		regexp.MustCompile("Hits as % of bytes sent:\t5min: ([0-9\\.]+)%"):        "byte_ratio",
		regexp.MustCompile("CPU Usage, 5 minute avg:\t([0-9\\.]+)%"):              "cpu_usage",
		regexp.MustCompile("Storage Swap capacity:[\t ]+([0-9\\.]+)% used"):       "swap_used_ratio",
		regexp.MustCompile("Storage Mem capacity:[\t ]+([0-9\\.]+)% used"):        "memory_used_ratio",
		regexp.MustCompile("Maximum number of file descriptors:[\t ]+([0-9]+)"):   "total_fd",
		regexp.MustCompile("Largest file desc currently in use:[\t ]+([0-9]+)"):   "max_fd",
		regexp.MustCompile("Number of file desc currently in use:[\t ]+([0-9]+)"): "current_fd",
		regexp.MustCompile("Available number of file descriptors:[\t ]+([0-9]+)"): "avail_fd",
		regexp.MustCompile("Reserved number of file descriptors:[\t ]+([0-9]+)"):  "reserved_fd",
		regexp.MustCompile("Store Disk files open:[\t ]+([0-9]+)"):                "open_files",
		regexp.MustCompile("Files queued for open:[\t ]+([0-9]+)"):                "queued_files",
		regexp.MustCompile("memPoolAlloc calls:[\t ]+([0-9]+)"):                   "memory_poll_alloc",
		regexp.MustCompile("memPoolFree calls:[\t ]+([0-9]+)"):                    "memory_poll_free",
	}

	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)

		for rexp, key := range regexpmap {
			match := rexp.FindStringSubmatch(s)
			if match == nil {
				continue
			}

			v, err := strconv.ParseFloat(match[1], 64)
			if err != nil {
				return nil, err
			}
			stat[key] = v

			break
		}
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (m SquidPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "3128", "Port")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var squid SquidPlugin
	squid.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	helper := mp.NewMackerelPlugin(squid)
	helper.Tempfile = *optTempfile

	helper.Run()
}
