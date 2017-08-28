package mpflume

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

const (
	Channel = "CHANNEL"
	Sink    = "SINK"
	Source  = "SOURCE"
)

// FlumePlugin mackerel plugin
type FlumePlugin struct {
	URI    string
	Prefix string
}

var graphdef = map[string]mp.Graphs{}

// MetricKeyPrefix interface for PluginWithPrefix
func (p *FlumePlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "flume"
	}
	return p.Prefix
}

// GraphDefinition interface for mackerelplugin
func (p *FlumePlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

// FetchMetrics interface for mackerelplugin
func (p *FlumePlugin) FetchMetrics() (map[string]float64, error) {
	m, err := p.getMetrics()
	if err != nil {
		return nil, err
	}
	return p.parseMetrics(m), nil
}

func (p *FlumePlugin) getMetrics() (map[string]interface{}, error) {
	res, err := http.Get(p.URI)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var m map[string]interface{}
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func (p *FlumePlugin) parseMetrics(metrics map[string]interface{}) map[string]float64 {
	ret := make(map[string]float64)
	for k, v := range metrics {
		array := strings.Split(k, ".")
		typeName := array[0]
		componentName := array[1]
		switch typeName {
		case Channel:
			p.parseChannel(ret, componentName, v.(map[string]interface{}))
			p.addGraphdefChannel(componentName)
		case Sink:
			p.parseSink(ret, componentName, v.(map[string]interface{}))
			p.addGraphdefSink(componentName)
		case Source:
			p.parseSource(ret, componentName, v.(map[string]interface{}))
			p.addGraphdefSource(componentName)
		}
	}

	return ret
}

func (p *FlumePlugin) convertFloat64(value string) float64 {
	f, _ := strconv.ParseFloat(value, 64)
	return f
}

func (p *FlumePlugin) parseChannel(ret map[string]float64, componentName string, value map[string]interface{}) {
	ret[componentName + ".ChannelCapacity"]       = p.convertFloat64(value["ChannelCapacity"].(string))
	ret[componentName + ".ChannelSize"]           = p.convertFloat64(value["ChannelSize"].(string))
	ret[componentName + ".ChannelFillPercentage"] = p.convertFloat64(value["ChannelFillPercentage"].(string))
	ret[componentName + ".EventPutAttemptCount"]  = p.convertFloat64(value["EventPutAttemptCount"].(string))
	ret[componentName + ".EventPutSuccessCount"]  = p.convertFloat64(value["EventPutSuccessCount"].(string))
	ret[componentName + ".EventTakeAttemptCount"] = p.convertFloat64(value["EventTakeAttemptCount"].(string))
	ret[componentName + ".EventTakeSuccessCount"] = p.convertFloat64(value["EventTakeSuccessCount"].(string))
}

func (p *FlumePlugin) addGraphdefChannel(componentName string) {
	labelPrefix := strings.Title(p.Prefix + " " + componentName)
	graphdefPrefix := "channel."

	graphdef[graphdefPrefix + "capacity"] = mp.Graphs{
		Label: labelPrefix + " Channel Capacity",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: componentName + ".ChannelCapacity", Label: "Channel Capacity"},
			{Name: componentName + ".ChannelSize", Label: "Channel Size"},
		},
	}
	graphdef[graphdefPrefix + "use_rate"] = mp.Graphs{
		Label: labelPrefix + " Channel Use Rate",
		Unit:  "percentage",
		Metrics: []mp.Metrics{
			{Name: componentName + ".ChannelFillPercentage", Label: "Channel Fill Percentage"},
		},
	}
	graphdef[graphdefPrefix + "event_put_num"] = mp.Graphs{
		Label: labelPrefix + " Event Put Num",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: componentName + ".EventPutAttemptCount", Label: "Attempt Count", Diff: true},
			{Name: componentName + ".EventPutSuccessCount", Label: "Success Count", Diff: true},
		},
	}
	graphdef[graphdefPrefix + "event_take_num"] = mp.Graphs{
		Label: labelPrefix + " Event Take Num",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: componentName + ".EventTakeAttemptCount", Label: "Attempt Count", Diff: true},
			{Name: componentName + ".EventTakeSuccessCount", Label: "Success Count", Diff: true},
		},
	}
}

func (p *FlumePlugin) parseSink(ret map[string]float64, componentName string, value map[string]interface{}) {
	ret[componentName + ".BatchCompleteCount"]     = p.convertFloat64(value["BatchCompleteCount"].(string))
	ret[componentName + ".BatchEmptyCount"]        = p.convertFloat64(value["BatchEmptyCount"].(string))
	ret[componentName + ".BatchUnderflowCount"]    = p.convertFloat64(value["BatchUnderflowCount"].(string))
	ret[componentName + ".ConnectionCreatedCount"] = p.convertFloat64(value["ConnectionCreatedCount"].(string))
	ret[componentName + ".ConnectionClosedCount"]  = p.convertFloat64(value["ConnectionClosedCount"].(string))
	ret[componentName + ".ConnectionFailedCount"]  = p.convertFloat64(value["ConnectionFailedCount"].(string))
	ret[componentName + ".EventDrainAttemptCount"] = p.convertFloat64(value["EventDrainAttemptCount"].(string))
	ret[componentName + ".EventDrainSuccessCount"] = p.convertFloat64(value["EventDrainSuccessCount"].(string))
}

func (p *FlumePlugin) addGraphdefSink(componentName string) {
	labelPrefix := strings.Title(p.Prefix + " " + componentName)
	graphdefPrefix := "sink."

	graphdef[graphdefPrefix + "batch_num"] = mp.Graphs{
		Label: labelPrefix + " Batch Num",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: componentName + ".BatchCompleteCount", Label: "Complete Count", Diff: true},
			{Name: componentName + ".BatchEmptyCount", Label: "Empty Count", Diff: true},
			{Name: componentName + ".BatchUnderflowCount", Label: "Underflow Count", Diff: true},
		},
	}
	graphdef[graphdefPrefix + "connection"] = mp.Graphs{
		Label: labelPrefix + " Connection",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: componentName + ".ConnectionCreatedCount", Label: "Created Count", Diff: true},
			{Name: componentName + ".ConnectionClosedCount", Label: "Closed Count", Diff: true},
			{Name: componentName + ".ConnectionFailedCount", Label: "Failed Count", Diff: true},
		},
	}
	graphdef[graphdefPrefix + "event_drain_num"] = mp.Graphs{
		Label: labelPrefix + " Event Drain Num",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: componentName + ".EventDrainAttemptCount", Label: "Attempt Count", Diff: true},
			{Name: componentName + ".EventDrainSuccessCount", Label: "Success Count", Diff: true},
		},
	}
}

func (p *FlumePlugin) parseSource(ret map[string]float64, componentName string, value map[string]interface{}) {
	ret[componentName + ".AppendAcceptedCount"]      = p.convertFloat64(value["AppendAcceptedCount"].(string))
	ret[componentName + ".AppendReceivedCount"]      = p.convertFloat64(value["AppendReceivedCount"].(string))
	ret[componentName + ".AppendBatchAcceptedCount"] = p.convertFloat64(value["AppendBatchAcceptedCount"].(string))
	ret[componentName + ".AppendBatchReceivedCount"] = p.convertFloat64(value["AppendBatchReceivedCount"].(string))
	ret[componentName + ".EventAcceptedCount"]       = p.convertFloat64(value["EventAcceptedCount"].(string))
	ret[componentName + ".EventReceivedCount"]       = p.convertFloat64(value["EventReceivedCount"].(string))
	ret[componentName + ".OpenConnectionCount"]      = p.convertFloat64(value["OpenConnectionCount"].(string))
}

func (p *FlumePlugin) addGraphdefSource(componentName string) {
	labelPrefix := strings.Title(p.Prefix + " " + componentName)
	graphdefPrefix := "source."

	graphdef[graphdefPrefix + "append_num"] = mp.Graphs{
		Label: labelPrefix + " Append Num",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: componentName + ".AppendAcceptedCount", Label: "Accepted Count", Diff: true},
			{Name: componentName + ".AppendReceivedCount", Label: "Received Count", Diff: true},
		},
	}
	graphdef[graphdefPrefix + "append_batch_num"] = mp.Graphs{
		Label: labelPrefix + " Append Batch Num",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: componentName + ".AppendBatchAcceptedCount", Label: "Accepted Count", Diff: true},
			{Name: componentName + ".AppendBatchReceivedCount", Label: "Received Count", Diff: true},
		},
	}
	graphdef[graphdefPrefix + "event_num"] = mp.Graphs{
		Label: labelPrefix + " Event Num",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: componentName + ".EventAcceptedCount", Label: "Accepted Count", Diff: true},
			{Name: componentName + ".EventReceivedCount", Label: "Received Count", Diff: true},
		},
	}
	graphdef[graphdefPrefix + "connection"] = mp.Graphs{
		Label: labelPrefix + " Connection",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: componentName + ".OpenConnectionCount", Label: "Open Count"},
		},
	}
}

// Do the plugin
func Do() {
	optHost   := flag.String("host", "localhost", "Host Name")
	optPort   := flag.String("port", "41414", "Port")
	optPrefix := flag.String("metric-key-prefix", "", "Metric key prefix")
	flag.Parse()
	mp.NewMackerelPlugin(&FlumePlugin{
		URI:    fmt.Sprintf("http://%s:%s/metrics", *optHost, *optPort),
		Prefix: *optPrefix,
	}).Run()
}
