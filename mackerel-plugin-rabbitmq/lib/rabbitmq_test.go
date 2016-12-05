package mprabbitmq

import (
	"reflect"
	"testing"

	"github.com/michaelklishin/rabbit-hole"
	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var rabbitmq RabbitMQPlugin

	graphdef := rabbitmq.GraphDefinition()
	if len(graphdef) != 2 {
		t.Errorf("GetTempfilename: %d should be 2", len(graphdef))
	}
}

func TestParse(t *testing.T) {
	var rabbitmq RabbitMQPlugin

	var stub rabbithole.Overview
	stub.QueueTotals.Messages = 1
	stub.QueueTotals.MessagesReady = 2
	stub.QueueTotals.MessagesUnacknowledged = 3
	stub.MessageStats.PublishDetails.Rate = 4

	stat, err := rabbitmq.parseStats(stub)

	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stat["messages"]).String(), "float64")
	assert.EqualValues(t, stat["messages"], 1)
	assert.EqualValues(t, reflect.TypeOf(stat["publish"]).String(), "float64")
	assert.EqualValues(t, stat["publish"], 4)
}
