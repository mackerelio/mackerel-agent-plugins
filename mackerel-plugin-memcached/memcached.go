package main

import (
	"bufio"
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"memcached.connections": mp.Graphs{
		Label: "Memcached Connections",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "curr_connections", Label: "Connections", Diff: false},
		},
	},
	"memcached.cmd": mp.Graphs{
		Label: "Memcached Command",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "cmd_get", Label: "Get", Diff: true},
			mp.Metrics{Name: "cmd_set", Label: "Set", Diff: true},
			mp.Metrics{Name: "cmd_flush", Label: "Flush", Diff: true},
			mp.Metrics{Name: "cmd_touch", Label: "Touch", Diff: true},
		},
	},
	"memcached.hitmiss": mp.Graphs{
		Label: "Memcached Hits/Misses",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "get_hits", Label: "Get Hits", Diff: true},
			mp.Metrics{Name: "get_misses", Label: "Get Misses", Diff: true},
			mp.Metrics{Name: "delete_hits", Label: "Delete Hits", Diff: true},
			mp.Metrics{Name: "delete_misses", Label: "Delete Misses", Diff: true},
			mp.Metrics{Name: "incr_hits", Label: "Incr Hits", Diff: true},
			mp.Metrics{Name: "incr_misses", Label: "Incr Misses", Diff: true},
			mp.Metrics{Name: "cas_hits", Label: "Cas Hits", Diff: true},
			mp.Metrics{Name: "cas_misses", Label: "Cas Misses", Diff: true},
			mp.Metrics{Name: "touch_hits", Label: "Touch Hits", Diff: true},
			mp.Metrics{Name: "touch_misses", Label: "Touch Misses", Diff: true},
		},
	},
	"memcached.evictions": mp.Graphs{
		Label: "Memcached Evictions",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "evictions", Label: "Evictions", Diff: true},
		},
	},
	"memcached.unfetched": mp.Graphs{
		Label: "Memcached Unfetched",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "expired_unfetched", Label: "Expired unfetched", Diff: true},
			mp.Metrics{Name: "evicted_unfetched", Label: "Evicted unfetched", Diff: true},
		},
	},
	"memcached.rusage": mp.Graphs{
		Label: "Memcached Resouce Usage",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "rusage_user", Label: "User", Diff: true},
			mp.Metrics{Name: "rusage_system", Label: "System", Diff: true},
		},
	},
	"memcached.bytes": mp.Graphs{
		Label: "Memcached Traffics",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "bytes_read", Label: "Read", Diff: true},
			mp.Metrics{Name: "bytes_written", Label: "Write", Diff: true},
		},
	},
}

type MemcachedPlugin struct {
	Target   string
	Tempfile string
}

func (m MemcachedPlugin) FetchMetrics() (map[string]float64, error) {
	conn, err := net.Dial("tcp", m.Target)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(conn, "stats")
	scanner := bufio.NewScanner(conn)
	stat := make(map[string]float64)

	for scanner.Scan() {
		line := scanner.Text()
		s := string(line)
		if s == "END" {
			return stat, nil
		}

		res := strings.Split(s, " ")
		if res[0] == "STAT" {
			stat[res[1]], err = strconv.ParseFloat(res[2], 64)
			if err != nil {
				log.Println("FetchMetrics:", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return stat, err
	}
	return nil, nil
}

func (m MemcachedPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "11211", "Port")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var memcached MemcachedPlugin
	helper := mp.NewMackerelPlugin(memcached)

	memcached.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-memcached-%s-%s", *optHost, *optPort)
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
