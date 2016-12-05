package mpfluentd

import (
	"reflect"
	"regexp"
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
		{"foo/bar", "foo_bar"},
		{"foo:bar", "foo_bar"},
	}

	for _, testSet := range testSets {
		if normalizePluginID(testSet[0]) != testSet[1] {
			t.Errorf("normalizeMetricName: '%s' should be normalized to '%s', but '%s'", testSet[0], testSet[1], normalizePluginID(testSet[0]))
		}
	}
}

func TestParse(t *testing.T) {
	var fluentd FluentdMetrics
	stub := `{"plugins":[{"plugin_id":"object:3feb368cfad0","plugin_category":"output","type":"mackerel","config":{"type":"mackerel","api_key":"aaa","service":"foo","metrics_name":"${[1]}-bar.${out_key}","remove_prefix":"","out_keys":"Latency","localtime":true},"output_plugin":true,"buffer_queue_length":0,"buffer_total_queued_size":53,"retry_count":0},{"plugin_id":"object:155633c","plugin_category":"input","type":"monitor_agent","config":{"type":"monitor_agent","bind":"0.0.0.0","port":"24220"},"output_plugin":false,"retry_count":null}]}`

	fluentdStats := []byte(stub)

	stat, err := fluentd.parseStats(fluentdStats)
	assert.Nil(t, err)
	// Fluentd Stats
	assert.EqualValues(t, reflect.TypeOf(stat["fluentd.buffer_total_queued_size.object_3feb368cfad0"]).String(), "float64")
	assert.EqualValues(t, stat["fluentd.buffer_total_queued_size.object_3feb368cfad0"].(float64), 53)
	if _, ok := stat["fluentd.buffer_total_queued_size.object_155633c"]; ok {
		t.Errorf("parseStats: stats of other than the output plugin should not exist")
	}
}

func TestPluginTypeOption(t *testing.T) {

	stub := `{"plugins":[{"plugin_id":"out_file","plugin_category":"output","type":"file","config":{"buffer_chunk_limit":"12m","compress":"gzip","type":"file","path":"/path/to/log","time_slice_format":"%Y%m%d%H","time_slice_wait":"1m","buffer_type":"file","buffer_path":"/path/to/buffer"},"output_plugin":true,"buffer_queue_length":0,"buffer_total_queued_size":10940,"retry_count":0},{"plugin_id":"out_mackerel","plugin_category":"output","type":"mackerel","config":{"type":"mackerel","api_key":"aaa","service":"foo","metrics_name":"${[1]}-bar.${out_key}","remove_prefix":"","out_keys":"Latency","localtime":true},"output_plugin":true,"buffer_queue_length":0,"buffer_total_queued_size":53,"retry_count":0}]}`
	fluentdStats := []byte(stub)

	// Specify type option
	var fluentd = FluentdMetrics{pluginType: "mackerel"}
	stat, err := fluentd.parseStats(fluentdStats)

	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stat["fluentd.buffer_total_queued_size.out_mackerel"]).String(), "float64")
	assert.EqualValues(t, stat["fluentd.buffer_total_queued_size.out_mackerel"].(float64), 53)
	if _, ok := stat["fluentd.buffer_total_queued_size.out_file"]; ok {
		t.Errorf("parseStats: stats of plugin that do not match the specified type should not exist")
	}

	// Not specify type option
	fluentd = FluentdMetrics{}
	stat, err = fluentd.parseStats(fluentdStats)

	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stat["fluentd.buffer_total_queued_size.out_mackerel"]).String(), "float64")
	assert.EqualValues(t, stat["fluentd.buffer_total_queued_size.out_mackerel"].(float64), 53)
	assert.EqualValues(t, reflect.TypeOf(stat["fluentd.buffer_total_queued_size.out_file"]).String(), "float64")
	assert.EqualValues(t, stat["fluentd.buffer_total_queued_size.out_file"].(float64), 10940)
}

func TestPluginIDPatternOption(t *testing.T) {

	stub := `{"plugins":[{"plugin_id":"match_plugin_id","plugin_category":"output","type":"file","config":{"buffer_chunk_limit":"12m","compress":"gzip","type":"file","path":"/path/to/log","time_slice_format":"%Y%m%d%H","time_slice_wait":"1m","buffer_type":"file","buffer_path":"/path/to/buffer"},"output_plugin":true,"buffer_queue_length":0,"buffer_total_queued_size":10940,"retry_count":0},{"plugin_id":"do_not_match_plugin_id","plugin_category":"output","type":"mackerel","config":{"type":"mackerel","api_key":"aaa","service":"foo","metrics_name":"${[1]}-bar.${out_key}","remove_prefix":"","out_keys":"Latency","localtime":true},"output_plugin":true,"buffer_queue_length":0,"buffer_total_queued_size":53,"retry_count":0}]}`
	fluentdStats := []byte(stub)

	// Specify type option
	var fluentd = FluentdMetrics{
		pluginIDPattern: regexp.MustCompile("^match"),
	}
	stat, err := fluentd.parseStats(fluentdStats)

	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stat["fluentd.buffer_total_queued_size.match_plugin_id"]).String(), "float64")
	assert.EqualValues(t, stat["fluentd.buffer_total_queued_size.match_plugin_id"].(float64), 10940)
	if _, ok := stat["fluentd.buffer_total_queued_size.do_not_match_plugin_id"]; ok {
		t.Errorf("parseStats: stats of plugin that do not match the specified id pattern should not exist")
	}

	// Not specify type option
	fluentd = FluentdMetrics{}
	stat, err = fluentd.parseStats(fluentdStats)

	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stat["fluentd.buffer_total_queued_size.match_plugin_id"]).String(), "float64")
	assert.EqualValues(t, stat["fluentd.buffer_total_queued_size.match_plugin_id"].(float64), 10940)
	assert.EqualValues(t, reflect.TypeOf(stat["fluentd.buffer_total_queued_size.do_not_match_plugin_id"]).String(), "float64")
	assert.EqualValues(t, stat["fluentd.buffer_total_queued_size.do_not_match_plugin_id"].(float64), 53)
}
