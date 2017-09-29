package mpmcrouter

import (
	"testing"
)

func TestFetchMetrics(t *testing.T) {
	expected := map[string]interface{}{
		"cmd_add_count":                    uint64(941424),
		"cmd_cas_count":                    uint64(0),
		"cmd_decr_count":                   uint64(0),
		"cmd_delete_count":                 uint64(0),
		"cmd_get_count":                    uint64(1380598681),
		"cmd_gets_count":                   uint64(0),
		"cmd_incr_count":                   uint64(2571178),
		"cmd_lease_get_count":              uint64(0),
		"cmd_lease_set_count":              uint64(0),
		"cmd_meta_count":                   uint64(0),
		"cmd_other_count":                  uint64(0),
		"cmd_replace_count":                uint64(0),
		"cmd_set_count":                    uint64(2929856),
		"cmd_stats_count":                  uint64(0),
		"result_busy_all_count":            uint64(0),
		"result_busy_count":                uint64(0),
		"result_connect_error_all_count":   uint64(0),
		"result_connect_error_count":       uint64(0),
		"result_connect_timeout_all_count": uint64(0),
		"result_connect_timeout_count":     uint64(0),
		"result_data_timeout_all_count":    uint64(2368),
		"result_data_timeout_count":        uint64(1951),
		"result_error_all_count":           uint64(532481),
		"result_error_count":               uint64(408682),
		"result_local_error_all_count":     uint64(0),
		"result_local_error_count":         uint64(0),
		"result_tko_all_count":             uint64(530113),
		"result_tko_count":                 uint64(406731),
		"duration_us":                      float64(2653.1359895317773),
	}

	p := &McrouterPlugin{
		StatsFile: "testdata/libmcrouter.mcrouter.6000.stats",
	}
	metrics, err := p.FetchMetrics()

	if err != nil {
		t.Errorf("Failed to FetchMetrics: %s", err)
		return
	}

	for key, expectedValue := range expected {
		gotValue, ok := metrics[key]
		if !ok {
			t.Errorf("metric of %s cannot be fetched", key)
			continue
		}
		if gotValue != expectedValue {
			t.Errorf("metric of %s should be %v, but %v", key, gotValue, expectedValue)
		}
	}
}
