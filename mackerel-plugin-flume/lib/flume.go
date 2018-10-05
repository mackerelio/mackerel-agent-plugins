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
	// Channel is Flume Channel Type
	Channel = "CHANNEL"
	// Sink is Flume Sink Type
	Sink = "SINK"
	// Source is Flume Source Type
	Source = "SOURCE"
)

// FlumePlugin mackerel plugin
type FlumePlugin struct {
	URI    string
	Prefix string
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p *FlumePlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "flume"
	}
	return p.Prefix
}

// GraphDefinition interface for mackerelplugin
func (p *FlumePlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(p.Prefix)
	return map[string]mp.Graphs{
		"channel.capacity.#": {
			Label: labelPrefix + " Channel Capacity",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "ChannelCapacity", Label: "capacity"},
				{Name: "ChannelSize", Label: "size"},
			},
		},
		"channel.use_rate.#": {
			Label: labelPrefix + " Channel Use Rate",
			Unit:  "percentage",
			Metrics: []mp.Metrics{
				{Name: "ChannelFillPercentage", Label: "fill percentage"},
			},
		},
		"channel.event_put_num.#": {
			Label: labelPrefix + " Channel Event Put Num",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "EventPutAttemptCount", Label: "attempt", Diff: true},
				{Name: "EventPutSuccessCount", Label: "success", Diff: true},
			},
		},
		"channel.event_take_num.#": {
			Label: labelPrefix + " Channel Event Take Num",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "EventTakeAttemptCount", Label: "attempt", Diff: true},
				{Name: "EventTakeSuccessCount", Label: "success", Diff: true},
			},
		},
		"sink.batch_num.#": {
			Label: labelPrefix + " Sink Batch Num",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "BatchCompleteCount", Label: "complete", Diff: true},
				{Name: "BatchEmptyCount", Label: "empty", Diff: true},
				{Name: "BatchUnderflowCount", Label: "underflow", Diff: true},
			},
		},
		"sink.connection.#": {
			Label: labelPrefix + " Sink Connection",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "ConnectionCreatedCount", Label: "created", Diff: true},
				{Name: "ConnectionClosedCount", Label: "closed", Diff: true},
				{Name: "ConnectionFailedCount", Label: "failed", Diff: true},
			},
		},
		"sink.event_drain_num.#": {
			Label: labelPrefix + " Sink Event Drain Num",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "EventDrainAttemptCount", Label: "attempt", Diff: true},
				{Name: "EventDrainSuccessCount", Label: "success", Diff: true},
			},
		},
		"source.append_num.#": {
			Label: labelPrefix + " Source Append Num",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "AppendAcceptedCount", Label: "accepted", Diff: true},
				{Name: "AppendReceivedCount", Label: "received", Diff: true},
			},
		},
		"source.append_batch_num.#": {
			Label: labelPrefix + " Source Append Batch Num",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "AppendBatchAcceptedCount", Label: "accepted", Diff: true},
				{Name: "AppendBatchReceivedCount", Label: "received", Diff: true},
			},
		},
		"source.event_num.#": {
			Label: labelPrefix + " Source Event Num",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "EventAcceptedCount", Label: "accepted", Diff: true},
				{Name: "EventReceivedCount", Label: "received", Diff: true},
			},
		},
		"source.connection.#": {
			Label: labelPrefix + " Source Connection",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "OpenConnectionCount", Label: "open"},
			},
		},
	}
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
	req, err := http.NewRequest(http.MethodGet, p.URI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "mackerel-plugin-flume")

	res, err := http.DefaultClient.Do(req)
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
		case Sink:
			p.parseSink(ret, componentName, v.(map[string]interface{}))
		case Source:
			p.parseSource(ret, componentName, v.(map[string]interface{}))
		}
	}

	return ret
}

func (p *FlumePlugin) convertFloat64(value string) float64 {
	f, _ := strconv.ParseFloat(value, 64)
	return f
}

func (p *FlumePlugin) parseChannel(ret map[string]float64, componentName string, value map[string]interface{}) {
	ret["channel.capacity."+componentName+".ChannelCapacity"] = p.convertFloat64(value["ChannelCapacity"].(string))
	ret["channel.capacity."+componentName+".ChannelSize"] = p.convertFloat64(value["ChannelSize"].(string))
	ret["channel.use_rate."+componentName+".ChannelFillPercentage"] = p.convertFloat64(value["ChannelFillPercentage"].(string))
	ret["channel.event_put_num."+componentName+".EventPutAttemptCount"] = p.convertFloat64(value["EventPutAttemptCount"].(string))
	ret["channel.event_put_num."+componentName+".EventPutSuccessCount"] = p.convertFloat64(value["EventPutSuccessCount"].(string))
	ret["channel.event_take_num."+componentName+".EventTakeAttemptCount"] = p.convertFloat64(value["EventTakeAttemptCount"].(string))
	ret["channel.event_take_num."+componentName+".EventTakeSuccessCount"] = p.convertFloat64(value["EventTakeSuccessCount"].(string))
}

func (p *FlumePlugin) parseSink(ret map[string]float64, componentName string, value map[string]interface{}) {
	ret["sink.batch_num."+componentName+".BatchCompleteCount"] = p.convertFloat64(value["BatchCompleteCount"].(string))
	ret["sink.batch_num."+componentName+".BatchEmptyCount"] = p.convertFloat64(value["BatchEmptyCount"].(string))
	ret["sink.batch_num."+componentName+".BatchUnderflowCount"] = p.convertFloat64(value["BatchUnderflowCount"].(string))
	ret["sink.connection."+componentName+".ConnectionCreatedCount"] = p.convertFloat64(value["ConnectionCreatedCount"].(string))
	ret["sink.connection."+componentName+".ConnectionClosedCount"] = p.convertFloat64(value["ConnectionClosedCount"].(string))
	ret["sink.connection."+componentName+".ConnectionFailedCount"] = p.convertFloat64(value["ConnectionFailedCount"].(string))
	ret["sink.event_drain_num."+componentName+".EventDrainAttemptCount"] = p.convertFloat64(value["EventDrainAttemptCount"].(string))
	ret["sink.event_drain_num."+componentName+".EventDrainSuccessCount"] = p.convertFloat64(value["EventDrainSuccessCount"].(string))
}

func (p *FlumePlugin) parseSource(ret map[string]float64, componentName string, value map[string]interface{}) {
	ret["source.append_num."+componentName+".AppendAcceptedCount"] = p.convertFloat64(value["AppendAcceptedCount"].(string))
	ret["source.append_num."+componentName+".AppendReceivedCount"] = p.convertFloat64(value["AppendReceivedCount"].(string))
	ret["source.append_batch_num."+componentName+".AppendBatchAcceptedCount"] = p.convertFloat64(value["AppendBatchAcceptedCount"].(string))
	ret["source.append_batch_num."+componentName+".AppendBatchReceivedCount"] = p.convertFloat64(value["AppendBatchReceivedCount"].(string))
	ret["source.event_num."+componentName+".EventAcceptedCount"] = p.convertFloat64(value["EventAcceptedCount"].(string))
	ret["source.event_num."+componentName+".EventReceivedCount"] = p.convertFloat64(value["EventReceivedCount"].(string))
	ret["source.connection."+componentName+".OpenConnectionCount"] = p.convertFloat64(value["OpenConnectionCount"].(string))
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Host Name")
	optPort := flag.String("port", "41414", "Port")
	optPrefix := flag.String("metric-key-prefix", "", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()
	plugin := mp.NewMackerelPlugin(&FlumePlugin{
		URI:    fmt.Sprintf("http://%s:%s/metrics", *optHost, *optPort),
		Prefix: *optPrefix,
	})
	plugin.Tempfile = *optTempfile
	plugin.Run()
}
