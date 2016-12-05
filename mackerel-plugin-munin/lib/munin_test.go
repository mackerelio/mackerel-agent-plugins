package mpmunin

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var plg = "swap"
var muninms = map[string](*MuninMetric){}
var ss = services{}

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
hoge.label ho GE
hoge.type DERIVE
`
	var title string

	parsePluginConfig(stub, &muninms, &title)

	assert.EqualValues(t, title, "Swap in/out")

	var met *MuninMetric

	assert.NotNil(t, muninms["swap_in"])
	met = muninms["swap_in"]
	assert.EqualValues(t, met.Label, "swapIn")
	assert.EqualValues(t, met.Type, "DERIVE")
	assert.EqualValues(t, met.Draw, "STACK")
	assert.EqualValues(t, met.Value, "")

	assert.NotNil(t, muninms["swap_out"])
	met = muninms["swap_out"]
	assert.EqualValues(t, met.Label, "swapOut")
	assert.EqualValues(t, met.Type, "DERIVE")
	assert.EqualValues(t, met.Draw, "")
	assert.EqualValues(t, met.Value, "")

	assert.NotNil(t, muninms["hoge"])
	met = muninms["hoge"]
	assert.EqualValues(t, met.Label, "ho GE")
	assert.EqualValues(t, met.Type, "DERIVE")
	assert.EqualValues(t, met.Draw, "")
	assert.EqualValues(t, met.Value, "")
}

func TestParsePluginVals(t *testing.T) {
	stub := `swap_out.value 2833950519
swap_in.value 2833950530
`
	parsePluginVals(stub, &muninms)

	var met *MuninMetric

	assert.NotNil(t, muninms["swap_in"])
	met = muninms["swap_in"]
	assert.EqualValues(t, met.Value, "2833950530")

	assert.NotNil(t, muninms["swap_out"])
	met = muninms["swap_out"]
	assert.EqualValues(t, met.Value, "2833950519")
}

func TestRemoveUselessMetrics(t *testing.T) {
	removeUselessMetrics(&muninms)

	assert.Nil(t, muninms["hoge"])
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

	var s serviceEnvs

	s = ss["*"]
	assert.EqualValues(t, s["hoge"], "wwww")
	assert.EqualValues(t, s["poyo"], "1 2 3")
	assert.EqualValues(t, s["sharp"], "# #")

	s = ss["s*"]
	assert.EqualValues(t, s["foo"], "bababa")

	s = ss["swap*"]
	assert.EqualValues(t, s["OutOfService"], "")
	assert.EqualValues(t, s["foo"], "bar")
	assert.EqualValues(t, s["hoge"], "wild")
	assert.EqualValues(t, s["foo2"], "bar2")

	s = ss["swap"]
	assert.EqualValues(t, s["hoge"], "abs")
	assert.EqualValues(t, s["piyo"], "")

	pluginconfstub2 := `
[swap]
env.piyo piYO
`
	getEnvSettingsReader(&ss, plg, bytes.NewBufferString(pluginconfstub2))
	s = ss["swap"]
	assert.EqualValues(t, s["piyo"], "piYO")
}

func TestCompileEnvPairs(t *testing.T) {
	envs := *compileEnvPairs(&ss, plg)

	assert.EqualValues(t, envs["snap"], "")
	assert.EqualValues(t, envs["OutOfService"], "")

	assert.EqualValues(t, envs["poyo"], "1 2 3")
	assert.EqualValues(t, envs["sharp"], "# #")
	assert.EqualValues(t, envs["foo"], "bar")
	assert.EqualValues(t, envs["foo2"], "bar2")
	assert.EqualValues(t, envs["hoge"], "abs")
	assert.EqualValues(t, envs["piyo"], "piYO")
}
