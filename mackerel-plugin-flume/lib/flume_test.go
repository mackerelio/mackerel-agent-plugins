package mpflume

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsonDate = `{"SOURCE.source":{"EventReceivedCount":"260969","Type":"SOURCE","AppendBatchAcceptedCount":"8357","EventAcceptedCount":"260969","AppendReceivedCount":"0","StartTime":"1503541553238","OpenConnectionCount":"0","AppendAcceptedCount":"0","AppendBatchReceivedCount":"8357","StopTime":"0"},"CHANNEL.channel":{"ChannelCapacity":"5000000","ChannelFillPercentage":"0.0","Type":"CHANNEL","ChannelSize":"0","EventTakeSuccessCount":"260969","StartTime":"1503541553170","EventTakeAttemptCount":"277651","EventPutAttemptCount":"260969","EventPutSuccessCount":"260969","StopTime":"0"},"SINK.sink":{"ConnectionCreatedCount":"109","Type":"SINK","ConnectionClosedCount":"108","BatchCompleteCount":"1567","BatchEmptyCount":"11067","EventDrainAttemptCount":"260969","StartTime":"1503541553174","EventDrainSuccessCount":"260969","BatchUnderflowCount":"5615","StopTime":"0","ConnectionFailedCount":"0"}}`

/*
{
  "SOURCE.source": {
    "EventReceivedCount": "260969",
    "Type": "SOURCE",
    "AppendBatchAcceptedCount": "8357",
    "EventAcceptedCount": "260969",
    "AppendReceivedCount": "0",
    "StartTime": "1503541553238",
    "OpenConnectionCount": "0",
    "AppendAcceptedCount": "0",
    "AppendBatchReceivedCount": "8357",
    "StopTime": "0"
  },
  "CHANNEL.channel": {
    "ChannelCapacity": "5000000",
    "ChannelFillPercentage": "0.0",
    "Type": "CHANNEL",
    "ChannelSize": "0",
    "EventTakeSuccessCount": "260969",
    "StartTime": "1503541553170",
    "EventTakeAttemptCount": "277651",
    "EventPutAttemptCount": "260969",
    "EventPutSuccessCount": "260969",
    "StopTime": "0"
  },
  "SINK.sink": {
    "ConnectionCreatedCount": "109",
    "Type": "SINK",
    "ConnectionClosedCount": "108",
    "BatchCompleteCount": "1567",
    "BatchEmptyCount": "11067",
    "EventDrainAttemptCount": "260969",
    "StartTime": "1503541553174",
    "EventDrainSuccessCount": "260969",
    "BatchUnderflowCount": "5615",
    "StopTime": "0",
    "ConnectionFailedCount": "0"
  }
}
*/

func getTestData() map[string]interface{} {
	var data map[string]interface{}
	_ = json.Unmarshal([]byte(jsonDate), &data)
	return data
}

func TestParseMetrics(t *testing.T) {
	flume := &FlumePlugin{}
	ret := flume.parseMetrics(getTestData())

	// Channel
	assert.EqualValues(t, ret["channel.capacity.channel.ChannelCapacity"], 5000000)
	assert.EqualValues(t, ret["channel.capacity.channel.ChannelSize"], 0)
	assert.EqualValues(t, ret["channel.use_rate.channel.ChannelFillPercentage"], 0)
	assert.EqualValues(t, ret["channel.event_put_num.channel.EventPutAttemptCount"], 260969)
	assert.EqualValues(t, ret["channel.event_put_num.channel.EventPutSuccessCount"], 260969)
	assert.EqualValues(t, ret["channel.event_take_num.channel.EventTakeAttemptCount"], 277651)
	assert.EqualValues(t, ret["channel.event_take_num.channel.EventTakeSuccessCount"], 260969)
	// Sink
	assert.EqualValues(t, ret["sink.batch_num.sink.BatchCompleteCount"], 1567)
	assert.EqualValues(t, ret["sink.batch_num.sink.BatchEmptyCount"], 11067)
	assert.EqualValues(t, ret["sink.batch_num.sink.BatchUnderflowCount"], 5615)
	assert.EqualValues(t, ret["sink.connection.sink.ConnectionCreatedCount"], 109)
	assert.EqualValues(t, ret["sink.connection.sink.ConnectionClosedCount"], 108)
	assert.EqualValues(t, ret["sink.connection.sink.ConnectionFailedCount"], 0)
	assert.EqualValues(t, ret["sink.event_drain_num.sink.EventDrainAttemptCount"], 260969)
	assert.EqualValues(t, ret["sink.event_drain_num.sink.EventDrainSuccessCount"], 260969)
	// Source
	assert.EqualValues(t, ret["source.append_num.source.AppendAcceptedCount"], 0)
	assert.EqualValues(t, ret["source.append_num.source.AppendReceivedCount"], 0)
	assert.EqualValues(t, ret["source.append_batch_num.source.AppendBatchAcceptedCount"], 8357)
	assert.EqualValues(t, ret["source.append_batch_num.source.AppendBatchReceivedCount"], 8357)
	assert.EqualValues(t, ret["source.event_num.source.EventAcceptedCount"], 260969)
	assert.EqualValues(t, ret["source.event_num.source.EventReceivedCount"], 260969)
	assert.EqualValues(t, ret["source.connection.source.OpenConnectionCount"], 0)
}
