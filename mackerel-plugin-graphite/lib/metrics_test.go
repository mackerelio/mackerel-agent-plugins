package mpgraphite

import (
	"testing"
)

func TestGetInstanceName(t *testing.T) {
	cases := []struct{ target, expected string }{
		{"carbon.agents.t_e_s_t-1.avgUpdateTime", "1"},
		{"carbon.agents.t_e_s_t-2.cache.size", "2"},
	}

	for _, tc := range cases {
		m := metrics{}
		m.Target = tc.target
		if actual := m.getInstanceName(); actual != tc.expected {
			t.Errorf("getInstanceName(): %s should be '%s'", actual, tc.expected)
		}
	}
}

func TestGetDestinationName(t *testing.T) {
	cases := []struct{ target, expected string }{
		{"carbon.relays.t_e_s_t-a.avgUpdateTime", ""},
		{"carbon.relays.t_e_s_t-a.destinations.127_0_0_1:3004:a.sent", "127_0_0_1:3004:a"},
	}

	for _, tc := range cases {
		m := metrics{}
		m.Target = tc.target
		if actual := m.getDestinationName(); actual != tc.expected {
			t.Errorf("getDestinationName(): %s should be '%s'", actual, tc.expected)
		}
	}
}

func TestGetMetricName(t *testing.T) {
	cases := []struct{ target, expected string }{
		{"carbon.agents.t_e_s_t-a.avgUpdateTime", "avgUpdateTime"},
		{"carbon.agents.t_e_s_t-a.cache.size", "cache_size"},
		{"carbon.relays.t_e_s_t-a.cpuUsage", "cpuUsage"},
		{"carbon.relays.t_e_s_t-a.destinations.127_0_0_1:3004:a.sent", "destinations_sent"},
	}

	for _, tc := range cases {
		m := metrics{}
		m.Target = tc.target
		if actual := m.getMetricName(); actual != tc.expected {
			t.Errorf("getMetricName(): %s should be '%s'", actual, tc.expected)
		}
	}
}

func TestGetUnitType(t *testing.T) {
	cases := []struct{ target, expected string }{
		{"carbon.agents.t_e_s_t-a.avgUpdateTime", "float"},
		{"carbon.agents.t_e_s_t-a.cache.size", "integer"},
		{"carbon.relays.t_e_s_t-a.cpuUsage", "float"},
		{"carbon.relays.t_e_s_t-a.destinations.127_0_0_1:3004:a.sent", "integer"},
	}

	for _, tc := range cases {
		m := metrics{}
		m.Target = tc.target
		if actual := m.getUnitType(); actual != tc.expected {
			t.Errorf("getUnitType(): %s should be '%s'", actual, tc.expected)
		}
	}
}

func TestCheckAllNil(t *testing.T) {
	cases := [](struct {
		target     string
		datapoints [][]interface{}
		expected   bool
	}){
		{"carbon.agents.t_e_s_t-a.avgUpdateTime", [][]interface{}{{nil, 1}, {nil, 2}}, true},
		{"carbon.agents.t_e_s_t-1.avgUpdateTime", [][]interface{}{{1, 1}, {nil, 2}}, false},
	}

	for _, tc := range cases {
		m := metrics{}
		m.Target = tc.target
		m.Datapoints = tc.datapoints
		if actual := m.isDataAllNil(); actual != tc.expected {
			t.Errorf("isDataAllNil(): %v should be '%v'", actual, tc.expected)
		}
	}
}

func TestGetMetricKey(t *testing.T) {
	cases := []struct{ target, expected string }{
		{"carbon.agents.t_e_s_t-a.avgUpdateTime", "graphite-carbon.cache.avgUpdateTime.a"},
		{"carbon.agents.t_e_s_t-a.cache.size", "graphite-carbon.cache.cache_size.a"},
		{"carbon.relays.t_e_s_t-a.cpuUsage", "graphite-carbon.relay.cpuUsage.cpuUsage"},
		{"carbon.relays.t_e_s_t-a.destinations.127_0_0_1:3004:a.sent", "graphite-carbon.relay.destinations_sent.127_0_0_1-3004-a"},
	}

	for _, tc := range cases {
		m := metrics{}
		m.Target = tc.target
		if actual := m.getMetricKey(); actual != tc.expected {
			t.Errorf("getMetricKey(): %s should be '%s'", actual, tc.expected)
		}
	}
}
