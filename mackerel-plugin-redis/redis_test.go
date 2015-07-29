package main

import (
	"testing"

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
		t.Errorf("something went wrong")
	}
	defer s.Stop()
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
		if _, ok := stat[v]; !ok {
			t.Errorf("metric of %s cannot be fetched", v)
		}
	}
}

func TestFetchMetrics(t *testing.T) {
	// should detect empty port
	portStr := "63331"
	s, err := redistest.NewServer(true, map[string]string{
		"port": portStr,
	})
	if err != nil {
		t.Errorf("something went wrong")
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
