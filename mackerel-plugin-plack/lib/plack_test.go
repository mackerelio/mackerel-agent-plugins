package mpplack

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var plack PlackPlugin

	graphdef := plack.GraphDefinition()
	if len(graphdef) != 3 {
		t.Errorf("GetTempfilename: %d should be 3", len(graphdef))
	}
}

func TestParse(t *testing.T) {
	var plack PlackPlugin
	stub := `{
  "Uptime": "1410520211",
  "TotalAccesses": "2",
  "IdleWorkers": "2",
  "TotalKbytes": "5",
  "BusyWorkers": "1",
  "stats": [
    {
      "pid": 11062,
      "method": "GET",
      "ss": 51,
      "remote_addr": "127.0.0.1",
      "host": "localhost:8000",
      "protocol": "HTTP/1.1",
      "status": "_",
      "uri": "/server-status?json"
    },
    {
      "ss": 41,
      "remote_addr": "127.0.0.1",
      "host": "localhost:8000",
      "protocol": "HTTP/1.1",
      "pid": 11063,
      "method": "GET",
      "status": "_",
      "uri": "/server-status?json"
    },
    {
      "ss": 0,
      "remote_addr": "127.0.0.1",
      "host": "localhost:8000",
      "protocol": "HTTP/1.1",
      "pid": 11064,
      "method": "GET",
      "status": "A",
      "uri": "/server-status?json"
    }
  ]
}
`

	plackStats := bytes.NewBufferString(stub)

	stat, err := plack.parseStats(plackStats)
	fmt.Println(stat)
	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stat["requests"]).String(), "uint64")
	assert.EqualValues(t, stat["requests"], 2)
	assert.EqualValues(t, reflect.TypeOf(stat["bytes_sent"]).String(), "uint64")
	assert.EqualValues(t, stat["bytes_sent"], 5)

	stubWithIntUptime := `
{"TotalKbytes":"36","IdleWorkers":"0","BusyWorkers":"0","TotalAccesses":"670","stats":[],"Uptime":1474047568}
`

	plackStatsWithIntUptime := bytes.NewBufferString(stubWithIntUptime)

	statWithIntUptime, err := plack.parseStats(plackStatsWithIntUptime)
	fmt.Println(statWithIntUptime)
	assert.Nil(t, err)
}
