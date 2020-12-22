package mpredis

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/soh335/go-test-redisserver"
)

var metrics = []string{
	"instantaneous_ops_per_sec", "total_connections_received", "rejected_connections", "connected_clients",
	"blocked_clients", "connected_slaves", "keys", "expires", "expired", "evicted_keys", "keyspace_hits", "keyspace_misses", "used_memory",
	"used_memory_rss", "used_memory_peak", "used_memory_lua", "uptime_in_seconds",
}

func TestFetchMetricsUnixSocket(t *testing.T) {
	s, err := redistest.NewServer(true, nil)
	if err != nil {
		t.Errorf("Failed to invoke testserver. %s", err)
		return
	}
	defer s.Stop()

	// set test data
	conn, err := redis.Dial("unix", s.Config["unixsocket"])
	if err != nil {
		t.Errorf("Failed to create a testclient. %s", err)
		return
	}
	_, err = conn.Do("SET", "TEST_KEY0", 1)
	if err != nil {
		t.Errorf("Failed to send a SET command. %s", err)
		return
	}
	conn.Do("SET", "TEST_KEY1", 1, "EX", 1)
	conn.Do("SET", "TEST_KEY2", 1, "EX", 2)
	conn.Do("SET", "TEST_KEY3", 1, "EX", 10)
	conn.Do("SET", "TEST_KEY4", 1, "EX", 20)
	conn.Do("SET", "TEST_KEY5", 1, "EX", 30)
	time.Sleep(3 * time.Second)

	redis := RedisPlugin{
		Timeout:       5,
		Prefix:        "redis",
		Socket:        s.Config["unixsocket"],
		ConfigCommand: "CONFIG",
	}
	stat, err := redis.FetchMetrics()

	if err != nil {
		t.Errorf("something went wrong")
	}

	for _, v := range metrics {
		value, ok := stat[v]
		if !ok {
			t.Errorf("metric of %s cannot be fetched", v)
		}
		if v == "keys" && value != 4.0 {
			t.Errorf("metric of key should be 4, but %v", value)
		}
		if v == "expires" && value != 3.0 {
			t.Errorf("metric of expires should be 3, but %v", value)
		}
		if v == "expired" && value != 2.0 {
			t.Errorf("metric of expired should be 2, but %v", value)
		}
	}
}

func TestFetchMetricsPercentageOfMemory(t *testing.T) {
	s, err := redistest.NewServer(true, nil)
	if err != nil {
		t.Errorf("Failed to invoke testserver. %s", err)
		return
	}
	defer s.Stop()

	rp := RedisPlugin{
		Timeout:       5,
		Prefix:        "redis",
		Socket:        s.Config["unixsocket"],
		ConfigCommand: "CONFIG",
	}

	conn, err := redis.Dial("unix", s.Config["unixsocket"])

	// Without maxmemory
	_, err = conn.Do("CONFIG", "SET", "maxmemory", 0.0)
	if err != nil {
		t.Errorf("Failed to send a CONFIG command. %s", err)
		return
	}

	stat1, err := rp.FetchMetrics()
	if err != nil {
		t.Errorf("something went wrong")
	}

	if value, ok := stat1["percentage_of_memory"]; !ok {
		t.Errorf("metric of 'percentage_of_memory' cannnot be fetched")
	} else if value != 0.0 {
		t.Errorf("metric of 'percentage_of_memory' should be 0.0, but %v", value)
	}

	// With maxmemory
	_, err = conn.Do("CONFIG", "SET", "maxmemory", 1024*1024)
	if err != nil {
		t.Errorf("Failed to send a CONFIG command. %s", err)
		return
	}

	stat2, err := rp.FetchMetrics()
	if err != nil {
		t.Errorf("something went wrong")
	}

	if value, ok := stat2["percentage_of_memory"]; !ok {
		t.Errorf("metric of 'percentage_of_memory' cannnot be fetched")
	} else if value == 0.0 {
		t.Errorf("metric of 'percentage_of_memory' should not be 0.0, but %v", value)
	}
}

func TestFetchMetrics(t *testing.T) {
	// should detect empty port
	portStr := "63331"
	s, err := redistest.NewServer(true, map[string]string{
		"port": portStr,
	})
	if err != nil {
		t.Errorf("Failed to invoke testserver. %s", err)
		return
	}
	defer s.Stop()
	redis := RedisPlugin{
		Timeout:       5,
		Prefix:        "redis",
		Port:          portStr,
		ConfigCommand: "CONFIG",
	}
	stat, err := redis.FetchMetrics()

	if err != nil {
		t.Errorf("something went wrong")
	}

	for _, v := range metrics {
		if _, ok := stat[v]; !ok {
			t.Errorf("metric of %s cannot be fetched", v)
		}
	}
}
