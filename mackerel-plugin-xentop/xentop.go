package main

import (
	"flag"
	"fmt"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"log"
	"os"
	"os/exec"
)

// これは多分必要なくなる
var graphdef map[string](mp.Graphs) = map[string](mp.Graphs){
	"xentop.cpu": mp.Graphs{
		Label:   "Xentop CPU",
		Unit:    "float",
		Metrics: [](mp.Metrics){},
	},
	"xentop.memory": mp.Graphs{
		Label:   "Xentop Memory",
		Unit:    "float",
		Metrics: [](mp.Metrics){},
	},
	"xentop.network": mp.Graphs{
		Label:   "Xentop Network",
		Unit:    "float",
		Metrics: [](mp.Metrics){},
	},
	"xentop.io": mp.Graphs{},
}

type XentopMetrics struct {
	HostName string
	Metrics  mp.Metrics
}

// ここに何を入れようか
type XentopPlugin struct {
	GraphName          string
	GraphUnit          string
	XentopMetricsSlice []XentopMetrics
}

func (m XentopPlugin) FetchMetrics() (map[string]float64, error) {
	// ここの中で何とかしてmetricsの配列に値を入れないとならない
	// statに値を入れながら，graphdefにキーを追加していく
	stat := make(map[string]float64)

	cmd := exec.Command("/bin/sh", "-c", "sudo xentop --batch -i 1 -f")
	//TODO 出力を標準出力に出さないようにする
	s, err := cmd.Run()
	if err != nil {
		return nil, err
	}

	//TODO 正規表現で必要な情報を抜き出す

}

// ここでグラフを定義する
func (m XentopPlugin) GraphDefinition() map[string](mp.Graphs) {
	metrics := []mp.Metrics{}
}

func main() {
	// flagの取得

	var xentop XentopPlugin

	helper := mp.NewMackerelPlugin(xentop)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
