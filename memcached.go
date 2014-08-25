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
			mp.Metrics{Key: "curr_connections", Label: "Connections", Diff: false},
		},
	},
	"memcached.cmd": mp.Graphs{
		Label: "Memcached Command",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Key: "cmd_get", Label: "Get", Diff: true},
			mp.Metrics{Key: "cmd_set", Label: "Set", Diff: true},
			mp.Metrics{Key: "cmd_flush", Label: "Flush", Diff: true},
			mp.Metrics{Key: "cmd_touch", Label: "Touch", Diff: true},
		},
	},
	"memcached.hitmiss": mp.Graphs{
		Label: "Memcached Hits/Misses",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Key: "get_hits", Label: "Get Hits", Diff: true},
			mp.Metrics{Key: "get_misses", Label: "Get Misses", Diff: true},
			mp.Metrics{Key: "delete_hits", Label: "Delete Hits", Diff: true},
			mp.Metrics{Key: "delete_misses", Label: "Delete Misses", Diff: true},
			mp.Metrics{Key: "incr_hits", Label: "Incr Hits", Diff: true},
			mp.Metrics{Key: "incr_misses", Label: "Incr Misses", Diff: true},
			mp.Metrics{Key: "cas_hits", Label: "Cas Hits", Diff: true},
			mp.Metrics{Key: "cas_misses", Label: "Cas Misses", Diff: true},
			mp.Metrics{Key: "touch_hits", Label: "Touch Hits", Diff: true},
			mp.Metrics{Key: "touch_misses", Label: "Touch Misses", Diff: true},
		},
	},
	"memcached.evictions": mp.Graphs{
		Label: "Memcached Evictions",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Key: "evictions", Label: "Evictions", Diff: true},
		},
	},
	"memcached.unfetched": mp.Graphs{
		Label: "Memcached Unfetched",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Key: "expired_unfetched", Label: "Expired unfetched", Diff: true},
			mp.Metrics{Key: "evicted_unfetched", Label: "Evicted unfetched", Diff: true},
		},
	},
	"memcached.rusage": mp.Graphs{
		Label: "Memcached Resouce Usage",
		Unit:  "float",
		Metrics: [](mp.Metrics){
			mp.Metrics{Key: "rusage_user", Label: "User", Diff: true},
			mp.Metrics{Key: "rusage_system", Label: "System", Diff: true},
		},
	},
	"memcached.bytes": mp.Graphs{
		Label: "Memcached Traffics",
		Unit:  "bytes",
		Metrics: [](mp.Metrics){
			mp.Metrics{Key: "bytes_read", Label: "Read", Diff: true},
			mp.Metrics{Key: "bytes_written", Label: "Write", Diff: true},
		},
	},
}

type MemcachedPlugin struct {
	Target   string
	Tempfile string
}

func (m MemcachedPlugin) FetchData() (map[string]float64, error) {
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
				log.Println("FetchData:", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return stat, err
	}
	return nil, nil
}

func (m MemcachedPlugin) GetGraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

func (m MemcachedPlugin) GetTempfilename() string {
	return m.Tempfile
}

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "11211", "Port")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var memcached MemcachedPlugin

	memcached.Target = fmt.Sprintf("%s:%s", *optHost, *optPort)
	if *optTempfile != "" {
		memcached.Tempfile = *optTempfile
	} else {
		memcached.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-memcached-%s-%s", *optHost, *optPort)
	}

	helper := mp.MackerelPlugin{memcached}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
