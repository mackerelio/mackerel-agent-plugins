package mpredis

import (
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/golib/logging"
)

var logger = logging.GetLogger("metrics.plugin.redis")

// RedisPlugin mackerel plugin for Redis
type RedisPlugin struct {
	Host          string
	Port          string
	Password      string
	Socket        string
	Prefix        string
	Timeout       int
	Tempfile      string
	ConfigCommand string
}

func authenticateByPassword(c redis.Conn, password string) error {
	if _, err := c.Do("AUTH", password); err != nil {
		logger.Errorf("Failed to authenticate. %s", err)
		return err
	}
	return nil
}

func (m RedisPlugin) fetchPercentageOfMemory(c redis.Conn, stat map[string]interface{}) error {
	res, err := redis.StringMap(c.Do(m.ConfigCommand, "GET", "maxmemory"))
	if err != nil {
		logger.Errorf("Failed to run `%s GET maxmemory` command. %s", m.ConfigCommand, err)
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
		stat["percentage_of_memory"] = 100.0 * stat["used_memory"].(float64) / maxsize
	}

	return nil
}

func (m RedisPlugin) fetchPercentageOfClients(c redis.Conn, stat map[string]interface{}) error {
	res, err := redis.StringMap(c.Do(m.ConfigCommand, "GET", "maxclients"))
	if err != nil {
		logger.Errorf("Failed to run `%s GET maxclients` command. %s", m.ConfigCommand, err)
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

func (m RedisPlugin) calculateCapacity(c redis.Conn, stat map[string]interface{}) error {
	if err := m.fetchPercentageOfMemory(c, stat); err != nil {
		return err
	}
	return m.fetchPercentageOfClients(c, stat)
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
	address := net.JoinHostPort(m.Host, m.Port)
	if m.Socket != "" {
		network = "unix"
		address = m.Socket
	}
	c, err := redis.Dial(network, address, redis.DialConnectTimeout(time.Duration(m.Timeout)*time.Second))
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

	str, err := redis.String(c.Do("info"))
	if err != nil {
		logger.Errorf("Failed to run info command. %s", err)
		return nil, err
	}

	stat := make(map[string]interface{})

	keysStat := 0.0
	expiresStat := 0.0
	var slaves []string

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

		if re, _ := regexp.MatchString("^slave\\d+", key); re {
			slaves = append(slaves, key)
			kv := strings.Split(value, ",")
			var offset, lag string
			if len(kv) == 5 {
				_, _, _, offset, lag = kv[0], kv[1], kv[2], kv[3], kv[4]
				lagKv := strings.SplitN(lag, "=", 2)
				lagFv, err := strconv.ParseFloat(lagKv[1], 64)
				if err != nil {
					logger.Warningf("Failed to parse slaves. %s", err)
				} else {
					stat[fmt.Sprintf("%s_lag", key)] = lagFv
				}
			} else {
				_, _, _, offset = kv[0], kv[1], kv[2], kv[3]
			}
			offsetKv := strings.SplitN(offset, "=", 2)
			offsetFv, err := strconv.ParseFloat(offsetKv[1], 64)
			if err != nil {
				logger.Warningf("Failed to parse slaves. %s", err)
				continue
			}
			stat[fmt.Sprintf("%s_offset_delay", key)] = offsetFv
			continue
		}

		if re, _ := regexp.MatchString("^db", key); re {
			kv := strings.SplitN(value, ",", 3)
			keys, expires := kv[0], kv[1]

			keysKv := strings.SplitN(keys, "=", 2)
			keysFv, err := strconv.ParseFloat(keysKv[1], 64)
			if err != nil {
				logger.Warningf("Failed to parse db keys. %s", err)
			} else {
				keysStat += keysFv
			}

			expiresKv := strings.SplitN(expires, "=", 2)
			expiresFv, err := strconv.ParseFloat(expiresKv[1], 64)
			if err != nil {
				logger.Warningf("Failed to parse db expires. %s", err)
			} else {
				expiresStat += expiresFv
			}

			continue
		}

		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
		stat[key] = v
	}

	stat["keys"] = keysStat
	stat["expires"] = expiresStat

	if _, ok := stat["expired_keys"]; ok {
		stat["expired"] = stat["expired_keys"]
	} else {
		stat["expired"] = 0.0
	}

	if m.ConfigCommand != "" {
		if err := m.calculateCapacity(c, stat); err != nil {
			logger.Infof("Failed to calculate capacity. (The cause may be that AWS Elasticache Redis has no `%s` command.) Skip these metrics. %s", m.ConfigCommand, err)
		}
	}

	for _, slave := range slaves {
		stat[fmt.Sprintf("%s_offset_delay", slave)] = stat["master_repl_offset"].(float64) - stat[fmt.Sprintf("%s_offset_delay", slave)].(float64)
	}

	return stat, nil
}

// GraphDefinition interface for mackerelplugin
func (m RedisPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(m.Prefix)

	var graphdef = map[string]mp.Graphs{
		"queries": {
			Label: (labelPrefix + " Queries"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "total_commands_processed", Label: "Queries", Diff: true},
			},
		},
		"connections": {
			Label: (labelPrefix + " Connections"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "total_connections_received", Label: "Connections", Diff: true, Stacked: true},
				{Name: "rejected_connections", Label: "Rejected Connections", Diff: true, Stacked: true},
			},
		},
		"clients": {
			Label: (labelPrefix + " Clients"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "connected_clients", Label: "Connected Clients", Diff: false, Stacked: true},
				{Name: "blocked_clients", Label: "Blocked Clients", Diff: false, Stacked: true},
				{Name: "connected_slaves", Label: "Connected Slaves", Diff: false, Stacked: true},
			},
		},
		"keys": {
			Label: (labelPrefix + " Keys"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "keys", Label: "Keys", Diff: false},
				{Name: "expires", Label: "Keys with expiration", Diff: false},
				{Name: "expired", Label: "Expired Keys", Diff: true},
				{Name: "evicted_keys", Label: "Evicted Keys", Diff: true},
			},
		},
		"keyspace": {
			Label: (labelPrefix + " Keyspace"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "keyspace_hits", Label: "Keyspace Hits", Diff: true},
				{Name: "keyspace_misses", Label: "Keyspace Missed", Diff: true},
			},
		},
		"memory": {
			Label: (labelPrefix + " Memory"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "used_memory", Label: "Used Memory", Diff: false},
				{Name: "used_memory_rss", Label: "Used Memory RSS", Diff: false},
				{Name: "used_memory_peak", Label: "Used Memory Peak", Diff: false},
				{Name: "used_memory_lua", Label: "Used Memory Lua engine", Diff: false},
			},
		},
		"capacity": {
			Label: (labelPrefix + " Capacity"),
			Unit:  "percentage",
			Metrics: []mp.Metrics{
				{Name: "percentage_of_memory", Label: "Percentage of memory", Diff: false},
				{Name: "percentage_of_clients", Label: "Percentage of clients", Diff: false},
			},
		},
		"uptime": {
			Label: (labelPrefix + " Uptime"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "uptime_in_seconds", Label: "Uptime In Seconds", Diff: false},
			},
		},
	}

	network := "tcp"
	address := net.JoinHostPort(m.Host, m.Port)
	if m.Socket != "" {
		network = "unix"
		address = m.Socket
	}

	c, err := redis.Dial(network, address, redis.DialConnectTimeout(time.Duration(m.Timeout)*time.Second))
	if err != nil {
		logger.Errorf("Failed to connect redis. %s", err)
		return nil
	}
	defer c.Close()

	if m.Password != "" {
		if err = authenticateByPassword(c, m.Password); err != nil {
			return nil
		}
	}

	str, err := redis.String(c.Do("info"))
	if err != nil {
		logger.Errorf("Failed to run info command. %s", err)
		return nil
	}

	var metricsLag []mp.Metrics
	var metricsOffsetDelay []mp.Metrics
	for _, line := range strings.Split(str, "\r\n") {
		if line == "" {
			continue
		}

		record := strings.SplitN(line, ":", 2)
		if len(record) < 2 {
			continue
		}
		key, _ := record[0], record[1]

		if re, _ := regexp.MatchString("^slave\\d+", key); re {
			metricsLag = append(metricsLag, mp.Metrics{Name: fmt.Sprintf("%s_lag", key), Label: fmt.Sprintf("Replication lag to %s", key), Diff: false})
			metricsOffsetDelay = append(metricsOffsetDelay, mp.Metrics{Name: fmt.Sprintf("%s_offset_delay", key), Label: fmt.Sprintf("Offset delay to %s", key), Diff: false})
		}
	}

	if len(metricsLag) > 0 {
		graphdef["lag"] = mp.Graphs{
			Label:   (labelPrefix + " Slave Lag"),
			Unit:    "seconds",
			Metrics: metricsLag,
		}
	}
	if len(metricsOffsetDelay) > 0 {
		graphdef["offset_delay"] = mp.Graphs{
			Label:   (labelPrefix + " Slave Offset Delay"),
			Unit:    "count",
			Metrics: metricsOffsetDelay,
		}
	}

	return graphdef
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "6379", "Port")
	optPassword := flag.String("password", os.Getenv("REDIS_PASSWORD"), "Password")
	optSocket := flag.String("socket", "", "Server socket (overrides host and port)")
	optPrefix := flag.String("metric-key-prefix", "redis", "Metric key prefix")
	optTimeout := flag.Int("timeout", 5, "Timeout")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optConfigCommand := flag.String("config-command", "CONFIG", "Custom CONFIG command. Disable CONFIG command when passed \"\".")

	flag.Parse()

	redis := RedisPlugin{
		Timeout:       *optTimeout,
		Prefix:        *optPrefix,
		ConfigCommand: *optConfigCommand,
	}
	if *optSocket != "" {
		redis.Socket = *optSocket
	} else {
		redis.Host = *optHost
		redis.Port = *optPort
		redis.Password = *optPassword
	}
	helper := mp.NewMackerelPlugin(redis)
	helper.Tempfile = *optTempfile

	helper.Run()
}
