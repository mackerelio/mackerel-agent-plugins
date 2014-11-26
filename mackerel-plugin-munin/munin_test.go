package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var plg string = "swap"
var ss Services = Services{}

func TestParsePluginConfig(t *testing.T) {
	stub := `graph_title Swap in/out
graph_args -l 0 --base 1000
graph_vlabel pages per ${graph_period} in (-) / out (+)
graph_category system
swap_in.label swapIn
swap_in.type DERIVE
swap_in.max 100000
swap_in.min 0
swap_in.graph no
swap_in.draw STACK
swap_in.value must_be_ignored
swap_out.label swapOut
swap_out.type DERIVE
swap_out.max 100000
swap_out.min 0
swap_out.negative swap_in
`
	muninms := make(map[string](*MuninMetric))
	var title string

	parsePluginConfig(stub, &muninms, &title)

	assert.Equal(t, title, "Swap in/out")

	var met *MuninMetric

	assert.NotNil(t, muninms["swap_in"])
	met = muninms["swap_in"]
	assert.Equal(t, met.Label, "swapIn")
	assert.Equal(t, met.Type, "DERIVE")
	assert.Equal(t, met.Draw, "STACK")
	assert.Equal(t, met.Value, "")

	assert.NotNil(t, muninms["swap_out"])
	met = muninms["swap_out"]
	assert.Equal(t, met.Label, "swapOut")
	assert.Equal(t, met.Type, "DERIVE")
	assert.Equal(t, met.Draw, "")
	assert.Equal(t, met.Value, "")
}

func TestParsePluginVals(t *testing.T) {
	stub := `swap_out.value 2833950519
swap_in.value 2833950530
`
	muninms := make(map[string](*MuninMetric))

	parsePluginVals(stub, &muninms)

	var met *MuninMetric

	assert.NotNil(t, muninms["swap_in"])
	met = muninms["swap_in"]
	assert.Equal(t, met.Value, "2833950530")

	assert.NotNil(t, muninms["swap_out"])
	met = muninms["swap_out"]
	assert.Equal(t, met.Value, "2833950519")
}

func TestGetEnvSettingsReader(t *testing.T) {
	pluginconfstub := `
env.OutOfService 1

[swap*]
   env.foo bar
env.hoge wild        	

[s*]
	env.foo bababa

[swap]
env.hoge abs
#env.piyo                 #        aaa

[snap]
env.snap yes
# single comment

[*]
env.hoge wwww
env.poyo 1 2 3 # commented...
env.sharp \# \# # commented

[swap*]
env.foo2 bar2
`
	getEnvSettingsReader(&ss, plg, bytes.NewBufferString(pluginconfstub))

	assert.Nil(t, ss["snap"])

	var s ServiceEnvs

	s = ss["*"]
	assert.Equal(t, s["hoge"], "wwww")
	assert.Equal(t, s["poyo"], "1 2 3")
	assert.Equal(t, s["sharp"], "# #")

	s = ss["s*"]
	assert.Equal(t, s["foo"], "bababa")

	s = ss["swap*"]
	assert.Equal(t, s["OutOfService"], "")
	assert.Equal(t, s["foo"], "bar")
	assert.Equal(t, s["hoge"], "wild")
	assert.Equal(t, s["foo2"], "bar2")

	s = ss["swap"]
	assert.Equal(t, s["hoge"], "abs")
	assert.Equal(t, s["piyo"], "")

	pluginconfstub2 := `
[swap]
env.piyo piYO
`
	getEnvSettingsReader(&ss, plg, bytes.NewBufferString(pluginconfstub2))
	s = ss["swap"]
	assert.Equal(t, s["piyo"], "piYO")
}

func TestCompileEnvPairs(t *testing.T) {
	envs := *compileEnvPairs(&ss, plg)

	assert.Equal(t, envs["snap"], "")
	assert.Equal(t, envs["OutOfService"], "")

	assert.Equal(t, envs["poyo"], "1 2 3")
	assert.Equal(t, envs["sharp"], "# #")
	assert.Equal(t, envs["foo"], "bar")
	assert.Equal(t, envs["foo2"], "bar2")
	assert.Equal(t, envs["hoge"], "abs")
	assert.Equal(t, envs["piyo"], "piYO")
}
