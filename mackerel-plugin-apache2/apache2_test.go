package main

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

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
