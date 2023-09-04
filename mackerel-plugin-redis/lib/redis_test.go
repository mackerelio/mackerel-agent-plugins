//go:build linux

package mpredis

import (
	"strings"
	"testing"

	"github.com/go-redis/redismock/v9"
)

var metrics = []string{
	"instantaneous_ops_per_sec", "total_connections_received", "rejected_connections", "connected_clients",
	"blocked_clients", "connected_slaves", "keys", "expires", "expired", "evicted_keys", "keyspace_hits", "keyspace_misses", "used_memory",
	"used_memory_rss", "used_memory_peak", "used_memory_lua", "uptime_in_seconds",
}

func TestFetchMetrics(t *testing.T) {
	db, mock := redismock.NewClientMock()

	mock.ExpectInfo().SetVal(strings.Join([]string{
		"db0:keys=4,expires=3,avg_ttl=0",
		"expired_keys:2",

		"instantaneous_ops_per_sec:0",
		"total_connections_received:0",
		"rejected_connections:0",
		"connected_clients:1",
		"blocked_clients:1",
		"connected_slaves:0",
		"evicted_keys:0",
		"keyspace_hits:0",
		"keyspace_misses:0",
		"used_memory:0",
		"used_memory_rss:0",
		"used_memory_peak:0",
		"used_memory_lua:0",
		"uptime_in_seconds:1",
	},
		"\r\n"),
	)

	redis := RedisPlugin{
		rdb:           db,
		Timeout:       5,
		Prefix:        "redis",
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
	db, mock := redismock.NewClientMock()

	val := map[string]string{"maxmemory": "0.0"}
	mock.ExpectConfigGet("maxmemory").SetVal(val)
	rp := RedisPlugin{
		rdb:           db,
		Timeout:       5,
		Prefix:        "redis",
		ConfigCommand: "config",
	}

	stat1 := make(map[string]interface{})
	err := rp.fetchPercentageOfMemory(stat1)
	if err != nil {
		t.Errorf("something went wrong")
	}

	if value, ok := stat1["percentage_of_memory"]; !ok {
		t.Errorf("metric of 'percentage_of_memory' cannnot be fetched")
	} else if value != 0.0 {
		t.Errorf("metric of 'percentage_of_memory' should be 0.0, but %v", value)
	}
}

func TestFetchMetricsPercentageOfMemory_100percent(t *testing.T) {
	db, mock := redismock.NewClientMock()

	val := map[string]string{"maxmemory": "1048576.0"}
	mock.ExpectConfigGet("maxmemory").SetVal(val)
	rp := RedisPlugin{
		rdb:           db,
		Timeout:       5,
		Prefix:        "redis",
		ConfigCommand: "config",
	}

	stat1 := make(map[string]interface{})
	stat1["used_memory"] = float64(1048576)
	err := rp.fetchPercentageOfMemory(stat1)
	if err != nil {
		t.Errorf("something went wrong")
	}

	if value, ok := stat1["percentage_of_memory"]; !ok {
		t.Errorf("metric of 'percentage_of_memory' cannnot be fetched")
	} else if value == 0.0 {
		t.Errorf("metric of 'percentage_of_memory' should not be 0.0, but %v", value)
	}
}
