package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fzzy/radix/redis"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metrics.plugin.redis")

// RedisPlugin mackerel plugin for Redis
type RedisPlugin struct {
	Host     string
	Port     string
	Password string
	Socket   string
	Prefix   string
	Timeout  int
	Tempfile string
}

func authenticateByPassword(c *redis.Client, password string) error {
	if r := c.Cmd("AUTH", password); r.Err != nil {
		logger.Errorf("Faild to authenticate. %s", r.Err)
		return r.Err
	}
	return nil
}

func fetchPercentageOfMemory(c *redis.Client, stat map[string]interface{}) error {
	r := c.Cmd("CONFIG", "GET", "maxmemory")
	if r.Err != nil {
		logger.Errorf("Failed to run `CONFIG GET maxmemory` command. %s", r.Err)
		return r.Err
	}

	res, err := r.Hash()
	if err != nil {
		logger.Errorf("Failed to fetch maxmemory. %s", err)
		return err
	}

	maxsize, err := strconv.ParseFloat(res["maxmemory"], 64)
	if err != nil {
		logger.Errorf("Failed to parse maxmemory. %s", err)
		return err
	}

	if maxsize == 0.0 {
		stat["percentage_of_memory"] = 0.0
	} else {
		stat["percentage_of_memory"] = 100.0 * stat["used_memoty"].(float64) / maxsize
	}

	return nil
}

func fetchPercentageOfClients(c *redis.Client, stat map[string]interface{}) error {
	r := c.Cmd("CONFIG", "GET", "maxclients")
	if r.Err != nil {
		logger.Errorf("Failed to run `CONFIG GET maxclients` command. %s", r.Err)
		return r.Err
	}

	res, err := r.Hash()
	if err != nil {
		logger.Errorf("Failed to fetch maxclients. %s", err)
		return err
	}

	maxsize, err := strconv.ParseFloat(res["maxclients"], 64)
	if err != nil {
		logger.Errorf("Failed to parse maxclients. %s", err)
		return err
	}

	stat["percentage_of_clients"] = 100.0 * stat["connected_clients"].(float64) / maxsize

	return nil
}

func calculateCapacity(c *redis.Client, stat map[string]interface{}) error {
	if err := fetchPercentageOfMemory(c, stat); err != nil {
		return err
	}
	if err := fetchPercentageOfClients(c, stat); err != nil {
		return err
	}
	return nil
}

// MetricKeyPrefix interface for PluginWithPrefix
func (m RedisPlugin) MetricKeyPrefix() string {
	if m.Prefix == "" {
		m.Prefix = "redis"
	}
	return m.Prefix
}

// FetchMetrics interface for mackerelplugin
func (m RedisPlugin) FetchMetrics() (map[string]interface{}, error) {
	network := "tcp"
	target := fmt.Sprintf("%s:%s", m.Host, m.Port)
	if m.Socket != "" {
		target = m.Socket
		network = "unix"
	}
	c, err := redis.DialTimeout(network, target, time.Duration(m.Timeout)*time.Second)
	if err != nil {
		logger.Errorf("Failed to connect redis. %s", err)
		return nil, err
	}
	defer c.Close()

	if m.Password != "" {
		if err = authenticateByPassword(c, m.Password); err != nil {
			return nil, err
		}
	}

	r := c.Cmd("info")
	if r.Err != nil {
		logger.Errorf("Failed to run info command. %s", r.Err)
		return nil, r.Err
	}
	str, err := r.Str()
	if err != nil {
		logger.Errorf("Failed to fetch information. %s", err)
		return nil, err
	}

	stat := make(map[string]interface{})

	keysStat := 0.0
	expiredStat := 0.0

	for _, line := range strings.Split(str, "\r\n") {
		if line == "" {
			continue
		}
		if re, _ := regexp.MatchString("^#", line); re {
			continue
		}

		record := strings.SplitN(line, ":", 2)
		if len(record) < 2 {
			continue
		}
		key, value := record[0], record[1]

		if re, _ := regexp.MatchString("^db", key); re {
			kv := strings.SplitN(value, ",", 3)
			keys, expired := kv[0], kv[1]

			keysKv := strings.SplitN(keys, "=", 2)
			keysFv, err := strconv.ParseFloat(keysKv[1], 64)
			if err != nil {
				logger.Warningf("Failed to parse db keys. %s", err)
			}
			keysStat += keysFv

			expiredKv := strings.SplitN(expired, "=", 2)
			expiredFv, err := strconv.ParseFloat(expiredKv[1], 64)
			if err != nil {
				logger.Warningf("Failed to parse db expired. %s", err)
			}
			expiredStat += expiredFv

			continue
		}

		stat[key], err = strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
	}

	stat["keys"] = keysStat
	stat["expired"] = expiredStat

	if _, ok := stat["keys"]; !ok {
		stat["keys"] = 0
	}
	if _, ok := stat["expired"]; !ok {
		stat["expired"] = 0
	}

	if err := calculateCapacity(c, stat); err != nil {
		logger.Infof("Failed to calculate capacity. (The cause may be that AWS Elasticache Redis has no `CONFIG` command.) Skip these metrics. %s", err)
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (m RedisPlugin) GraphDefinition() map[string](mp.Graphs) {
	labelPrefix := strings.Title(m.Prefix)

	var graphdef = map[string](mp.Graphs){
		"queries": mp.Graphs{
			Label: (labelPrefix + " Queries"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "instantaneous_ops_per_sec", Label: "Queries", Diff: false},
			},
		},
		"connections": mp.Graphs{
			Label: (labelPrefix + " Connections"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "total_connections_received", Label: "Connections", Diff: true, Stacked: true},
				mp.Metrics{Name: "rejected_connections", Label: "Rejected Connections", Diff: true, Stacked: true},
			},
		},
		"clients": mp.Graphs{
			Label: (labelPrefix + " Clients"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "connected_clients", Label: "Connected Clients", Diff: false, Stacked: true},
				mp.Metrics{Name: "blocked_clients", Label: "Blocked Clients", Diff: false, Stacked: true},
				mp.Metrics{Name: "connected_slaves", Label: "Blocked Clients", Diff: false, Stacked: true},
			},
		},
		"keys": mp.Graphs{
			Label: (labelPrefix + " Keys"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "keys", Label: "Keys", Diff: false},
				mp.Metrics{Name: "expired", Label: "Expired Keys", Diff: false},
			},
		},
		"keyspace": mp.Graphs{
			Label: (labelPrefix + " Keyspace"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "keyspace_hits", Label: "Keyspace Hits", Diff: true},
				mp.Metrics{Name: "keyspace_misses", Label: "Keyspace Missed", Diff: true},
			},
		},
		"memory": mp.Graphs{
			Label: (labelPrefix + " Memory"),
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "used_memory", Label: "Used Memory", Diff: false},
				mp.Metrics{Name: "used_memory_rss", Label: "Used Memory RSS", Diff: false},
				mp.Metrics{Name: "used_memory_peak", Label: "Used Memory Peak", Diff: false},
				mp.Metrics{Name: "used_memory_lua", Label: "Used Memory Lua engine", Diff: false},
			},
		},
		"capacity": mp.Graphs{
			Label: (labelPrefix + " Capacity"),
			Unit:  "percentage",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "percentage_of_memory", Label: "Percentage of memory", Diff: false},
				mp.Metrics{Name: "percentage_of_clients", Label: "Percentage of clients", Diff: false},
			},
		},
	}

	return graphdef
}

func main() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "6379", "Port")
	optPassowrd := flag.String("password", "", "Password")
	optSocket := flag.String("socket", "", "Server socket (overrides host and port)")
	optPrefix := flag.String("metric-key-prefix", "redis", "Metric key prefix")
	optTimeout := flag.Int("timeout", 5, "Timeout")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	redis := RedisPlugin{
		Timeout: *optTimeout,
		Prefix:  *optPrefix,
	}
	if *optSocket != "" {
		redis.Socket = *optSocket
	} else {
		redis.Host = *optHost
		redis.Port = *optPort
		redis.Password = *optPassowrd
	}
	helper := mp.NewMackerelPlugin(redis)
	helper.Tempfile = *optTempfile

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
