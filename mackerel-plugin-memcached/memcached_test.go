package main

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var memcached MemcachedPlugin

	graphdef := memcached.GraphDefinition()
	if len(graphdef) != 7 {
		t.Errorf("GetTempfilename: %d should be 3", len(graphdef))
	}
}

func TestParse(t *testing.T) {
	var memcached MemcachedPlugin
	stub := `STAT pid 1994
STAT uptime 92066123
STAT time 1436890963
STAT version 1.4.0
STAT pointer_size 64
STAT rusage_user 1393.803107
STAT rusage_system 2947.180187
STAT curr_connections 1003
STAT total_connections 965032539
STAT connection_structures 16388
STAT cmd_get 4306259844
STAT cmd_set 2423543841
STAT cmd_flush 0
STAT get_hits 2769383483
STAT get_misses 1536876361
STAT delete_misses 244469885
STAT delete_hits 14456835
STAT incr_misses 0
STAT incr_hits 0
STAT decr_misses 0
STAT decr_hits 0
STAT cas_misses 0
STAT cas_hits 0
STAT cas_badval 0
STAT bytes_read 8328670869009
STAT bytes_written 9151962263382
STAT limit_maxbytes 2147483648
STAT accepting_conns 1
STAT listen_disabled_num 0
STAT threads 5
STAT conn_yields 1487476
STAT bytes 621371972
STAT curr_items 955652
STAT total_items 2423543841
STAT evictions 236677775
END
`

	memcachedStats := bytes.NewBufferString(stub)

	stat, err := memcached.ParseStats(memcachedStats)
	fmt.Println(stat)
	assert.Nil(t, err)
	// Memcached Stats
	assert.EqualValues(t, reflect.TypeOf(stat["get_hits"]).String(), "string")
	assert.EqualValues(t, stat["get_hits"].(string), "2769383483")
}
