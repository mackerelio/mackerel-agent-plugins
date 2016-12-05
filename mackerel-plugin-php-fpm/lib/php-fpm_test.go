package mpphpfpm

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	jsonStr := `{
    "pool":"www",
    "process manager":"dynamic",
    "start time":1461398921,
    "start since":1624,
    "accepted conn":664,
    "listen queue":1,
    "max listen queue":3,
    "listen queue len":2,
    "idle processes":40,
    "active processes":10,
    "total processes":50,
    "max active processes":100,
    "max children reached":200,
    "slow requests":1000
  }`

	httpmock.RegisterResponder("GET", "http://httpmock/status?json",
		httpmock.NewStringResponder(200, jsonStr))

	p := PhpFpmPlugin{
		URL:     "http://httpmock/status?json",
		Prefix:  "php-fpm",
		Timeout: 5,
	}
	status, _ := getStatus(p)

	assert.EqualValues(t, 50, status.TotalProcesses)
	assert.EqualValues(t, 10, status.ActiveProcesses)
	assert.EqualValues(t, 40, status.IdleProcesses)
	assert.EqualValues(t, 100, status.MaxActiveProcesses)
	assert.EqualValues(t, 200, status.MaxChildrenReached)
	assert.EqualValues(t, 1, status.ListenQueue)
	assert.EqualValues(t, 2, status.ListenQueueLen)
	assert.EqualValues(t, 3, status.MaxListenQueue)
	assert.EqualValues(t, 1000, status.SlowRequests)
}
