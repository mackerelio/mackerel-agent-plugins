package main

import (
	"flag"
	"os"
	"strconv"

	up "github.com/ariarijp/uptimerobot-go/api"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

type UptimeRobotPlugin struct {
	ApiKey    string
	MonitorID int
	Name      string
	Tempfile  string
}

func (m UptimeRobotPlugin) FetchMetrics() (map[string]float64, error) {
	stat := make(map[string]float64)

	cli, err := up.NewClient(m.ApiKey)
	if err != nil {
		return nil, err
	}

	monitors := cli.Monitors()

	var req = up.GetMonitorsRequest{
		MonitorId: m.MonitorID,
	}

	resp, err := monitors.Get(req)
	if err != nil {
		return nil, err
	}

	stat["response_time"] = float64(resp.Monitors[0].ResponseTimes[0].Value)

	return stat, nil
}

func (m UptimeRobotPlugin) GraphDefinition() map[string](mp.Graphs) {
	label := "Response Time"
	if m.Name != "" {
		label = m.Name
	}

	return map[string](mp.Graphs){
		"uptimerobot.ResponseTime": mp.Graphs{
			Label: "Uptime Robot Response Time",
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "response_time", Label: label},
			},
		},
	}
}

func main() {
	optApiKey := flag.String("api-key", "", "API Key")
	optMonitorID := flag.String("monitor-id", "", "Monitor ID")
	optName := flag.String("name", "", "name")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var plugin UptimeRobotPlugin
	plugin.ApiKey = *optApiKey
	plugin.MonitorID, _ = strconv.Atoi(*optMonitorID)

	if *optName != "" {
		plugin.Name = *optName
	}

	helper := mp.NewMackerelPlugin(plugin)

	if *optTempfile != "" {
		helper.Tempfile = *optTempfile
	} else {
		helper.Tempfile = "/tmp/mackerel-plugin-uptimerobot"
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
