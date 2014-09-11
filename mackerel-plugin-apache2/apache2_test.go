package main

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestParseApache2Scoreboard( t *testing.T ){
    stub := "Scoreboard: W._SRWKDCLGI...."
    stat := make(map[string]float64)

    err := parseApache2Scoreboard( stub, &stat )
    assert.Nil( t, err )
    assert.Equal( t, stat[ "score-_" ], 1 )
    assert.Equal( t, stat[ "score-S" ], 1 )
    assert.Equal( t, stat[ "score-R" ], 1 )
    assert.Equal( t, stat[ "score-W" ], 2 )
    assert.Equal( t, stat[ "score-K" ], 1 )
    assert.Equal( t, stat[ "score-D" ], 1 )
    assert.Equal( t, stat[ "score-C" ], 1 )
    assert.Equal( t, stat[ "score-L" ], 1 )
    assert.Equal( t, stat[ "score-G" ], 1 )
    assert.Equal( t, stat[ "score-I" ], 1 )
    assert.Equal( t, stat[ "score-." ], 5 )
}

func TestParseApache2Status( t *testing.T ){
    stub := `Total Accesses: 358
Total kBytes: 20
CPULoad: .00117358
Uptime: 102251
ReqPerSec: .00350119
BytesPerSec: .200291
BytesPerReq: 57.2067
BusyWorkers: 1
IdleWorkers: 4
`
    stat := make(map[string]float64)

    err := parseApache2Status( stub, &stat )
    assert.Nil( t, err )
    assert.Equal( t, stat[ "requests" ], 358 )
    assert.Equal( t, stat[ "bytes_sent" ], 20 )
    assert.Equal( t, stat[ "cpu_load" ], 0.00117358 )
    assert.Equal( t, stat[ "busy_workers" ], 1 )
    assert.Equal( t, stat[ "idle_workers" ], 4 )
}

func TestGetApache2Metrics_1( t *testing.T ){
    ret, err := getApache2Metrics( "127.0.0.1", 1080, "/server-status?auto" )
    assert.Nil( t, err, "Please start-up your httpd (127.0.0.1:1080) or unable to connect httpd." )
    assert.NotNil( t, ret, )
    assert.NotEmpty( t, ret )
    assert.Contains( t, ret, "Total Accesses" )
    assert.Contains( t, ret, "Total kBytes" )
    assert.Contains( t, ret, "Uptime" )
    assert.Contains( t, ret, "BusyWorkers" )
    assert.Contains( t, ret, "IdleWorkers" )
    assert.Contains( t, ret, "Scoreboard" )
}
