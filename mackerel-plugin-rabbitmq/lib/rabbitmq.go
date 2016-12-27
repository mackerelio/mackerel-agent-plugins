package mprabbitmq

import (
	"flag"
	"os"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/michaelklishin/rabbit-hole"
)

var graphdef = map[string]mp.Graphs{
	"rabbitmq.queue": {
		Label: "RabbitMQ Queue",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "messages", Label: "Total", Diff: false},
			{Name: "ready", Label: "Ready", Diff: false},
			{Name: "unacknowledged", Label: "Unacknowledged", Diff: false},
		},
	},
	"rabbitmq.message": {
		Label: "RabbitMQ Message",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "publish", Label: "Publish", Diff: false},
		},
	},
}

// RabbitMQPlugin metrics
type RabbitMQPlugin struct {
	URI      string
	User     string
	Password string
	TempFile string
}

// FetchMetrics interface for mackerelplugin
func (r RabbitMQPlugin) FetchMetrics() (map[string]interface{}, error) {
	rmqc, err := rabbithole.NewClient(r.URI, r.User, r.Password)
	if err != nil {
		return nil, err
	}
	res, err := rmqc.Overview()
	if err != nil {
		return nil, err
	}

	return r.parseStats(*res)
}

func (r RabbitMQPlugin) parseStats(res rabbithole.Overview) (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	stat["messages"] = float64(res.QueueTotals.Messages)
	stat["ready"] = float64(res.QueueTotals.MessagesReady)
	stat["unacknowledged"] = float64(res.QueueTotals.MessagesUnacknowledged)
	stat["publish"] = float64(res.MessageStats.PublishDetails.Rate)

	return stat, nil

}

// GraphDefinition interface for mackerel plugin
func (r RabbitMQPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// Do the plugin
func Do() {
	optURI := flag.String("uri", "http://localhost:15672", "URI")
	optUser := flag.String("user", "guest", "User")
	optPass := flag.String("password", "guest", "Password")
	flag.Parse()

	var rabbitmq RabbitMQPlugin

	rabbitmq.URI = *optURI
	rabbitmq.User = *optUser
	rabbitmq.Password = *optPass

	helper := mp.NewMackerelPlugin(rabbitmq)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
