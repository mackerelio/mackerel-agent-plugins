package mpgraphite

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

var cacheGoldenMetrics = []metrics{
	{"carbon.agents.host_hoge-a.avgUpdateTime", [][]interface{}{{0.1, 1}, {0.2, 2}}},
	{"carbon.agents.host_hoge-b.avgUpdateTime", [][]interface{}{{0.1, 1}, {0.2, 2}}},
	{"carbon.agents.host_hoge-a.cache.size", [][]interface{}{{1, 1}, {2, 2}}},
	{"carbon.agents.host_hoge-b.cache.size", [][]interface{}{{1, 1}, {2, 2}}},
}

var relayGoldenMetrics = []metrics{
	{"carbon.relays.host_hoge-a.cpuUsage", [][]interface{}{{0.1, 1}, {0.2, 2}}},
	{"carbon.relays.host_hoge-a.destinations.127_0_0_1:3004:a.sent", [][]interface{}{{1, 1}, {2, 2}}},
	{"carbon.relays.host_hoge-a.destinations.127_0_0_1:3104:b.sent", [][]interface{}{{1, 1}, {2, 2}}},
}

var cacheHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	d, err := json.Marshal(cacheGoldenMetrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(d)
})

var relayHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	d, err := json.Marshal(relayGoldenMetrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(d)
})

func TestFetchData(t *testing.T) {
	ts := httptest.NewServer(cacheHandler)
	defer ts.Close()

	plugin := GraphitePlugin{
		Host:        "host_hoge",
		WebHost:     "webhost.hoge",
		WebPort:     "8000",
		Type:        "cache",
		Instance:    "*",
		LabelPrefix: "Carbon",
		URL:         ts.URL,
	}

	d, err := plugin.fetchData()
	if err != nil {
		t.Errorf("fetchData() failed: %v", err)
	}
	if actual := d[0].Target; actual != cacheGoldenMetrics[0].Target {
		t.Errorf("fetchData(): %s should be %s", actual, cacheGoldenMetrics[0].Target)
	}
	if actual := d[0].Datapoints[0][0]; actual != cacheGoldenMetrics[0].Datapoints[0][0] {
		t.Errorf("fetchData(): %s should be %s", actual, cacheGoldenMetrics[0].Datapoints[0][0])
	}
}

func TestCacheGraphDefinition(t *testing.T) {
	ts := httptest.NewServer(cacheHandler)
	defer ts.Close()

	plugin := GraphitePlugin{
		Host:        "host_hoge",
		WebHost:     "webhost.hoge",
		WebPort:     "8000",
		Type:        "cache",
		Instance:    "*",
		LabelPrefix: "Carbon",
		URL:         ts.URL,
	}

	graph := plugin.GraphDefinition()
	if actual := len(graph); actual != len(cacheMeta) {
		t.Errorf("GraphDefinition(): %d should be %d", actual, len(cacheMeta))
	}
	for _, g := range graph {
		if actual := len(g.Metrics); actual != 2 {
			t.Errorf("GraphDefinition(): %d should be 2", actual)
		}
	}
}

func TestRelayGraphDefinition(t *testing.T) {
	ts := httptest.NewServer(relayHandler)
	defer ts.Close()

	plugin := GraphitePlugin{
		Host:        "host_hoge",
		WebHost:     "webhost.hoge",
		WebPort:     "8000",
		Type:        "relay",
		Instance:    "a",
		LabelPrefix: "Carbon",
		URL:         ts.URL,
	}

	graph := plugin.GraphDefinition()
	if actual := len(graph); actual != len(relayMeta) {
		t.Errorf("GraphDefinition(): %d should be %d", actual, len(relayMeta))
	}
	for k, g := range graph {
		matched, _ := regexp.MatchString(`cpuUsage|memUsage|metricsRecieved`, k)
		if matched {
			if actual := len(g.Metrics); actual != 1 {
				t.Errorf("GraphDefinition(): %d should be 1", actual)
			}
		} else {
			if actual := len(g.Metrics); actual != 2 {
				t.Errorf("GraphDefinition(): %d should be 2", actual)
			}
		}
	}
}

func TestOutputValueForCache(t *testing.T) {
	ts := httptest.NewServer(cacheHandler)
	defer ts.Close()

	plugin := GraphitePlugin{
		Host:        "host_hoge",
		WebHost:     "webhost.hoge",
		WebPort:     "8000",
		Type:        "cache",
		Instance:    "*",
		LabelPrefix: "Carbon",
		URL:         ts.URL,
	}

	s := new(bytes.Buffer)
	plugin.outputValues(s)

	expected := `graphite-carbon.cache.avgUpdateTime.a	0.100000	1
graphite-carbon.cache.avgUpdateTime.a	0.200000	2
graphite-carbon.cache.avgUpdateTime.b	0.100000	1
graphite-carbon.cache.avgUpdateTime.b	0.200000	2
graphite-carbon.cache.cache_size.a	1	1
graphite-carbon.cache.cache_size.a	2	2
graphite-carbon.cache.cache_size.b	1	1
graphite-carbon.cache.cache_size.b	2	2
`

	if actual := string(s.Bytes()); actual != expected {
		t.Errorf("outputValues(): %s should be %s", actual, expected)
	}
}

func TestOutputValueForRelay(t *testing.T) {
	ts := httptest.NewServer(relayHandler)
	defer ts.Close()

	plugin := GraphitePlugin{
		Host:        "host_hoge",
		WebHost:     "webhost.hoge",
		WebPort:     "8000",
		Type:        "relay",
		Instance:    "a",
		LabelPrefix: "Carbon",
		URL:         ts.URL,
	}

	s := new(bytes.Buffer)
	plugin.outputValues(s)

	expected := `graphite-carbon.relay.cpuUsage.cpuUsage	0.100000	1
graphite-carbon.relay.cpuUsage.cpuUsage	0.200000	2
graphite-carbon.relay.destinations_sent.127_0_0_1-3004-a	1	1
graphite-carbon.relay.destinations_sent.127_0_0_1-3004-a	2	2
graphite-carbon.relay.destinations_sent.127_0_0_1-3104-b	1	1
graphite-carbon.relay.destinations_sent.127_0_0_1-3104-b	2	2
`

	if actual := string(s.Bytes()); actual != expected {
		t.Errorf("outputValues(): %s should be %s", actual, expected)
	}
}
