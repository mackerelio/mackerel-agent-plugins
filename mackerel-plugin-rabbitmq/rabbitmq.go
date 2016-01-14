package main

import (
	"flag"
	"os"

	"github.com/michaelklishin/rabbit-hole"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string](mp.Graphs){
	"rabbitmq.queue": mp.Graphs{
		Label: "RabbitMQ Queue",
		Unit: "integer",
		Metrics:[](mp.Metrics){
			mp.Metrics{Name: "messages", Label: "Total", Diff: false},
			mp.Metrics{Name: "ready", Label: "Ready", Diff: false},
			mp.Metrics{Name: "unacknowledged", Label: "Unacknowledged", Diff:false},
		},
	},
	"rabbitmq.message": mp.Graphs{
		Label: "RabbitMQ Message",
		Unit: "integer",
		Metrics:[](mp.Metrics){
			mp.Metrics{Name: "publish", Label: "Publish", Diff: false},
		},
	},
}

type RabbitMQPlugin struct {
	Url string
	User string
	Password string
	TempFile string
}

func (m RabbitMQPlugin) FetchMetrics() (map[string]interface{},error){
	rmqc, err:= rabbithole.NewClient(m.Url, m.User, m.Password)
	if err != nil {
		return nil, err
	}
	res, err := rmqc.Overview()
	if err != nil {
		return nil, err
	}

	return m.parseStats(*res)
}

func (m RabbitMQPlugin) parseStats(res rabbithole.Overview) (map[string]interface{},error){
	stat := make(map[string]interface{})

	stat["messages"] = float64(res.QueueTotals.Messages)
	stat["ready"] = float64(res.QueueTotals.MessagesReady)
	stat["unacknowledged"] = float64(res.QueueTotals.MessagesUnacknowledged)
	stat["publish"] = float64(res.MessageStats.PublishDetails.Rate)

	return stat, nil

}

func (m RabbitMQPlugin) GraphDefinition() map[string](mp.Graphs){
	return graphdef
}

func main(){
	optURI :=  flag.String("uri", "http://localhost:15672", "URI")
	optUser := flag.String("user", "guest", "User")
	optPass := flag.String("password", "guest", "Password")
	flag.Parse()

	var rabbitmq RabbitMQPlugin

	rabbitmq.Url = *optURI
	rabbitmq.User = *optUser
	rabbitmq.Password = *optPass

	helper := mp.NewMackerelPlugin(rabbitmq)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != ""{
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
