package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var plugin NVidiaSMIPlugin

	graphdef := plugin.GraphDefinition()
	if len(graphdef) != 5 {
		t.Errorf("GraphDef's size: %d should be 5", len(graphdef))
	}
}

func TestParse(t *testing.T) {
	var plugin NVidiaSMIPlugin
	plugin.Prefix = "nvidia.gpu"
	data := `10, 20, 30, 40, 1024, 64, 962
11, 21, 31, 41, 1024, 65, 961
`

	stats, err := plugin.parseStats(data)

	assert.Nil(t, err)

	assert.EqualValues(t, 10, stats["gpu.util.gpu0"])
	assert.EqualValues(t, 20, stats["memory.util.gpu0"])
	assert.EqualValues(t, 30, stats["temperature.gpu0"])
	assert.EqualValues(t, 40, stats["fanspeed.gpu0"])
	assert.EqualValues(t, 1024, stats["memory.usage.gpu0.total"])
	assert.EqualValues(t, 64, stats["memory.usage.gpu0.used"])
	assert.EqualValues(t, 962, stats["memory.usage.gpu0.free"])

	assert.EqualValues(t, 11, stats["gpu.util.gpu1"])
	assert.EqualValues(t, 21, stats["memory.util.gpu1"])
	assert.EqualValues(t, 31, stats["temperature.gpu1"])
	assert.EqualValues(t, 41, stats["fanspeed.gpu1"])
	assert.EqualValues(t, 1024, stats["memory.usage.gpu1.total"])
	assert.EqualValues(t, 65, stats["memory.usage.gpu1.used"])
	assert.EqualValues(t, 961, stats["memory.usage.gpu1.free"])
}
