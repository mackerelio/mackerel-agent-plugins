package mpsquid

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var squid SquidPlugin

	graphdef := squid.GraphDefinition()
	if len(graphdef) != 6 {
		t.Errorf("GetTempfilename: %d should be 6", len(graphdef))
	}
}

func TestParse(t *testing.T) {
	var squid SquidPlugin
	stubMgrCacheInfo := `Squid Object Cache: Version 3.5.19
Build Info: Debian linux
Service Name: squid
Start Time:	Wed, 28 Nov 2018 09:33:14 GMT
Current Time:	Thu, 29 Nov 2018 09:20:19 GMT
Connection information for squid:
	Number of clients accessing cache:	2
	Number of HTTP requests received:	1663
	Number of ICP messages received:	0
	Number of ICP messages sent:	0
	Number of queued ICP replies:	0
	Number of HTCP messages received:	0
	Number of HTCP messages sent:	0
	Request failure ratio:	 0.00
	Average HTTP requests per minute since start:	1.2
	Average ICP messages per minute since start:	0.0
	Select loop called: 8571054 times, 9.990 ms avg
Cache information for squid:
	Hits as % of all requests:	5min: 0.0%, 60min: 0.0%
	Hits as % of bytes sent:	5min: 100.0%, 60min: 100.0%
	Memory hits as % of hit requests:	5min: 0.0%, 60min: 0.0%
	Disk hits as % of hit requests:	5min: 0.0%, 60min: 0.0%
	Storage Swap size:	104 KB
	Storage Swap capacity:	 0.0% used, 100.0% free
	Storage Mem size:	324 KB
	Storage Mem capacity:	 0.0% used, 100.0% free
	Mean Object Size:	104.00 KB
	Requests given to unlinkd:	0
Median Service Times (seconds)  5 min    60 min:
	HTTP Requests (All):   0.00000  0.00000
	Cache Misses:          0.00000  0.00000
	Cache Hits:            0.00000  0.00000
	Near Hits:             0.00000  0.00000
	Not-Modified Replies:  0.00000  0.00000
	DNS Lookups:           0.00000  0.00000
	ICP Queries:           0.00000  0.00000
Resource usage for squid:
	UP Time:	85625.732 seconds
	CPU Time:	340.460 seconds
	CPU Usage:	0.40%
	CPU Usage, 5 minute avg:	0.41%
	CPU Usage, 60 minute avg:	0.41%
	Maximum Resident Size: 94096 KB
	Page faults with physical i/o: 0
Memory accounted for:
	Total accounted:          779 KB
	memPoolAlloc calls:    387958
	memPoolFree calls:     388241
File descriptor usage for squid:
	Maximum number of file descriptors:   4096
	Largest file desc currently in use:     27
	Number of file desc currently in use:   16
	Files queued for open:                   0
	Available number of file descriptors: 4080
	Reserved number of file descriptors:   100
	Store Disk files open:                   0
Internal Data Structures:
	    54 StoreEntries
	    54 StoreEntries with MemObjects
	    53 Hot Object Cache Items
	     1 on-disk objects
`

	squidStats := strings.NewReader(stubMgrCacheInfo)

	stat, err := squid.ParseMgrInfo(squidStats)
	assert.Nil(t, err)
	assert.EqualValues(t, 1663, stat["requests"])
	assert.EqualValues(t, 0, stat["request_ratio"])
	assert.EqualValues(t, 100, stat["byte_ratio"])
	assert.EqualValues(t, 0.41, stat["cpu_usage"])
	assert.EqualValues(t, 0, stat["swap_used_ratio"])
	assert.EqualValues(t, 0, stat["memory_used_ratio"])
	assert.EqualValues(t, 4096, stat["total_fd"])
	assert.EqualValues(t, 27, stat["max_fd"])
	assert.EqualValues(t, 16, stat["current_fd"])
	assert.EqualValues(t, 4080, stat["avail_fd"])
	assert.EqualValues(t, 100, stat["reserved_fd"])
	assert.EqualValues(t, 0, stat["open_files"])
	assert.EqualValues(t, 0, stat["queued_files"])
	assert.EqualValues(t, 387958, stat["memory_poll_alloc"])
	assert.EqualValues(t, 388241, stat["memory_poll_free"])
}
