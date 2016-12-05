package mpredis

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/soh335/go-test-redisserver"
)

var metrics = []string{
	"instantaneous_ops_per_sec", "total_connections_received", "rejected_connections", "connected_clients",
	"blocked_clients", "connected_slaves", "keys", "expired", "keyspace_hits", "keyspace_misses", "used_memory",
	"used_memory_rss", "used_memory_peak", "used_memory_lua",
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
	_, err = conn.Do("SET", "TEST_KEY", 1)
	if err != nil {
		t.Errorf("Failed to send a SET command. %s", err)
		return
	}
	_, err = conn.Do("SETEX", "TEST_EXPIRED_KEY", 1, 2)
	if err != nil {
		t.Errorf("Failed to send a SETEX command. %s", err)
		return
	}

	redis := RedisPlugin{
		Timeout: 5,
		Prefix:  "redis",
		Socket:  s.Config["unixsocket"],
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
		if v == "keys" && value != 2.0 {
			t.Errorf("metric of key should be 2, but %v", value)
		}
		if v == "expired" && value != 1.0 {
			t.Errorf("metric of expired should be 1, but %v", value)
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
		Timeout: 5,
		Prefix:  "redis",
		Socket:  s.Config["unixsocket"],
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
		Timeout: 5,
		Prefix:  "redis",
		Port:    portStr,
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
