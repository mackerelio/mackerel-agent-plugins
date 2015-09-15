package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var fluentd FluentdMetrics

	graphdef := fluentd.GraphDefinition()
	if len(graphdef) != 3 {
		t.Errorf("GetTempfilename: %d should be 3", len(graphdef))
	}
}

func TestNormalizePluginID(t *testing.T) {
	testSets := [][]string{
		[]string{"foo/bar", "foo_bar"},
		[]string{"foo:bar", "foo_bar"},
	}

	for _, testSet := range testSets {
		if normalizePluginID(testSet[0]) != testSet[1] {
			t.Errorf("normalizeMetricName: '%s' should be normalized to '%s', but '%s'", testSet[0], testSet[1], normalizePluginID(testSet[0]))
		}
	}
}

func TestParse(t *testing.T) {
	var fluentd FluentdMetrics
	stub := `{"plugins":[{"plugin_id":"object:3feb368cfad0","plugin_category":"output","type":"mackerel","config":{"type":"mackerel","api_key":"aaa","service":"foo","metrics_name":"${[1]}-bar.${out_key}","remove_prefix":"","out_keys":"Latency"},"output_plugin":true,"buffer_queue_length":0,"buffer_total_queued_size":53,"retry_count":0}]}`

	fluentdStats := []byte(stub)

	stat, err := fluentd.ParseStats(fluentdStats)
	//fmt.Println(stat)
	assert.Nil(t, err)
	// Fluentd Stats
	assert.EqualValues(t, reflect.TypeOf(stat["fluentd.buffer_total_queued_size.object_3feb368cfad0"]).String(), "float64")
	assert.EqualValues(t, stat["fluentd.buffer_total_queued_size.object_3feb368cfad0"].(float64), 53)
}
