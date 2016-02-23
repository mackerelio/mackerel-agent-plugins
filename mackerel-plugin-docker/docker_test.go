package main

import (
	"testing"
)

func TestNormalizeMetricName(t *testing.T) {
	testSets := [][]string{
		[]string{"foo/bar", "foo_bar"},
		[]string{"foo:bar", "foo_bar"},
	}

	for _, testSet := range testSets {
		if normalizeMetricName(testSet[0]) != testSet[1] {
			t.Errorf("normalizeMetricName: '%s' should be normalized to '%s', but '%s'", testSet[0], testSet[1], normalizeMetricName(testSet[0]))
		}
	}
}

func TestGraphDefinition(t *testing.T) {
	var docker DockerPlugin

	graphdef := docker.GraphDefinition()
	if len(graphdef) != 5 {
		t.Errorf("GetTempfilename: %d should be 5", len(graphdef))
	}
}

func TestParseStats121(t *testing.T) {
	var docker DockerPlugin
	stub := `
 {
     "read" : "2015-01-08T22:57:31.547920715Z",
     "networks": {
             "eth0": {
                 "rx_bytes": 5338,
                 "rx_dropped": 0,
                 "rx_errors": 0,
                 "rx_packets": 36,
                 "tx_bytes": 648,
                 "tx_dropped": 0,
                 "tx_errors": 0,
                 "tx_packets": 8
             },
             "eth5": {
                 "rx_bytes": 4641,
                 "rx_dropped": 0,
                 "rx_errors": 0,
                 "rx_packets": 26,
                 "tx_bytes": 690,
                 "tx_dropped": 0,
                 "tx_errors": 0,
                 "tx_packets": 9
             }
     },
     "memory_stats" : {
        "stats" : {
           "total_pgmajfault" : 0,
           "cache" : 0,
           "mapped_file" : 0,
           "total_inactive_file" : 0,
           "pgpgout" : 414,
           "rss" : 6537216,
           "total_mapped_file" : 0,
           "writeback" : 0,
           "unevictable" : 0,
           "pgpgin" : 477,
           "total_unevictable" : 0,
           "pgmajfault" : 0,
           "total_rss" : 6537216,
           "total_rss_huge" : 6291456,
           "total_writeback" : 0,
           "total_inactive_anon" : 0,
           "rss_huge" : 6291456,
           "hierarchical_memory_limit" : 67108864,
           "total_pgfault" : 964,
           "total_active_file" : 0,
           "active_anon" : 6537216,
           "total_active_anon" : 6537216,
           "total_pgpgout" : 414,
           "total_cache" : 0,
           "inactive_anon" : 0,
           "active_file" : 0,
           "pgfault" : 964,
           "inactive_file" : 0,
           "total_pgpgin" : 477
        },
        "max_usage" : 6651904,
        "usage" : 6537216,
        "failcnt" : 0,
        "limit" : 67108864
     },
     "blkio_stats" : {},
     "cpu_stats" : {
        "cpu_usage" : {
           "percpu_usage" : [
              16970827,
              1839451,
              7107380,
              10571290
           ],
           "usage_in_usermode" : 10000000,
           "total_usage" : 36488948,
           "usage_in_kernelmode" : 20000000
        },
        "system_cpu_usage" : 20091722000000000,
        "throttling_data" : {}
     }
  }`

	dockerStats := bytes.NewBufferString(stub)

	res := map[string]interface{}{}
	stat, err := docker.parseStats(&res, dockerStats)
	fmt.Println(stat)
	assert.Nil(t, err)
	// Docker Stats
	assert.EqualValues(t, stat["cpu_usage"], 36488948)

}

func TestParseStats117(t *testing.T) {
	var docker DockerPlugin
	stub := `
{
	"read" : "2015-01-08T22:57:31.547920715Z",
	"network" : {
		"rx_dropped" : 0,
		"rx_bytes" : 648,
		"rx_errors" : 0,
		"tx_packets" : 8,
		"tx_dropped" : 0,
		"rx_packets" : 8,
		"tx_errors" : 0,
		"tx_bytes" : 648
	},
	"memory_stats" : {
		"stats" : {
			"total_pgmajfault" : 0,
			"cache" : 0,
			"mapped_file" : 0,
			"total_inactive_file" : 0,
			"pgpgout" : 414,
			"rss" : 6537216,
			"total_mapped_file" : 0,
			"writeback" : 0,
			"unevictable" : 0,
			"pgpgin" : 477,
			"total_unevictable" : 0,
			"pgmajfault" : 0,
			"total_rss" : 6537216,
			"total_rss_huge" : 6291456,
			"total_writeback" : 0,
			"total_inactive_anon" : 0,
			"rss_huge" : 6291456,
			"hierarchical_memory_limit" : 67108864,
			"total_pgfault" : 964,
			"total_active_file" : 0,
			"active_anon" : 6537216,
			"total_active_anon" : 6537216,
			"total_pgpgout" : 414,
			"total_cache" : 0,
			"inactive_anon" : 0,
			"active_file" : 0,
			"pgfault" : 964,
			"inactive_file" : 0,
			"total_pgpgin" : 477
		},
		"max_usage" : 6651904,
		"usage" : 6537216,
		"failcnt" : 0,
		"limit" : 67108864
	},
	"blkio_stats" : {},
	"cpu_stats" : {
		"cpu_usage" : {
			"percpu_usage" : [
				16970827,
				1839451,
				7107380,
				10571290
			],
			"usage_in_usermode" : 10000000,
			"total_usage" : 36488948,
			"usage_in_kernelmode" : 20000000
		},
		"system_cpu_usage" : 20091722000000000,
		"throttling_data" : {}
	}
}`
	dockerStats := bytes.NewBufferString(stub)

	res := map[string]interface{}{}
	stat, err := docker.parseStats(&res, dockerStats)
	fmt.Println(stat)
	assert.Nil(t, err)
	// Docker Stats
	assert.EqualValues(t, stat["cpu_usage"], 36488948)

}
