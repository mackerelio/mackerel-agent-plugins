package mptwemproxy

import (
	"flag"
	"log"
	"net"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/lestrrat/go-tcptest"
)

var (
	statsServer *tcptest.TCPTest
	stats       string
)

func TestMain(m *testing.M) {
	os.Exit(mainTest(m))
}

func mainTest(m *testing.M) int {
	flag.Parse()

	// run a test server returning a dummy stats json
	startStatsServer := func(port int) {
		err := startTCPServer(port)
		if err != nil {
			panic("Failed to run a test server: " + err.Error())
		}
	}
	server, err := tcptest.Start(startStatsServer, 10*time.Second)
	if err != nil {
		panic("Failed to start a stats server: " + err.Error())
	}
	statsServer = server

	log.Printf("Started a stats server. port=%d", statsServer.Port())

	return m.Run()
}

func startTCPServer(port int) error {
	server, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}

	ch := make(chan net.Conn)
	go func() {
		for {
			client, err := server.Accept()
			if err != nil {
				log.Printf("Failed to accept: err=%v", err)
			}
			ch <- client
		}
	}()

	for {
		go func(client net.Conn) {
			_, err = client.Write([]byte(stats))
			if err != nil {
				log.Printf("Failed to write %v: err=%v", stats, err)
			}
		}(<-ch)
	}
}

var jsonStr = `{
    "service": "nutcracker",
    "source": "ip-10-0-1-10",
    "version": "0.4.1",
    "uptime": 74635,
    "timestamp": 1477054533,
    "total_connections": 3895,
    "curr_connections": 272,
    "redis-index": {
        "index1.cache:6379": {
            "out_queue_bytes": 80,
            "out_queue": 1,
            "in_queue_bytes": 50,
            "in_queue": 2,
            "requests": 3908,
            "request_bytes": 170558,
            "responses": 3908,
            "response_bytes": 176918,
            "server_connections": 4,
            "server_eof": 2,
            "server_timedout": 3,
            "server_ejected_at": 0,
            "server_err": 5
        },
        "client_eof": 1716,
        "client_err": 10,
        "client_connections": 121,
        "server_ejects": 20,
        "forward_error": 1,
        "fragments": 0
    },
    "redis/budget": {
        "budget1.cache:6379": {
            "out_queue_bytes": 81,
            "out_queue": 2,
            "in_queue_bytes": 51,
            "in_queue": 3,
            "requests": 3909,
            "request_bytes": 170559,
            "responses": 3909,
            "response_bytes": 176919,
            "server_connections": 5,
            "server_eof": 3,
            "server_timedout": 4,
            "server_ejected_at": 0,
            "server_err": 3
        },
        "budget2.cache:6379": {
            "out_queue_bytes": 83,
            "out_queue": 4,
            "in_queue_bytes": 53,
            "in_queue": 5,
            "requests": 3911,
            "request_bytes": 170561,
            "responses": 3911,
            "response_bytes": 176921,
            "server_connections": 7,
            "server_eof": 5,
            "server_timedout": 6,
            "server_ejected_at": 2,
            "server_err": 5
        },
        "client_eof": 2716,
        "client_err": 30,
        "client_connections": 221,
        "server_ejects": 40,
        "forward_error": 2,
        "fragments": 0
    }
}`

func TestFetchMetrics(t *testing.T) {
	// response a valid stats json
	stats = jsonStr

	// get metrics
	p := TwemproxyPlugin{
		Address: "localhost:" + strconv.Itoa(statsServer.Port()),
		Prefix:  "twemproxy",
		Timeout: 5,
	}
	metrics, err := p.FetchMetrics()
	if err != nil {
		t.Errorf("Failed to FetchMetrics: %s", err)
		return
	}

	// check the metrics
	expected := map[string]uint64{
		"total_connections":                                                         3895,
		"curr_connections":                                                          272,
		"pool_error.redis-index.client_err":                                         10,
		"pool_error.redis-index.server_ejects":                                      20,
		"pool_error.redis-index.forward_error":                                      1,
		"pool_client_connections.redis-index.client_eof":                            1716,
		"pool_client_connections.redis-index.client_connections":                    121,
		"server_error.redis-index_index1_cache_6379.server_err":                     5,
		"server_error.redis-index_index1_cache_6379.server_timedout":                3,
		"server_connections.redis-index_index1_cache_6379.server_eof":               2,
		"server_connections.redis-index_index1_cache_6379.server_connections":       4,
		"server_queue.redis-index_index1_cache_6379.out_queue":                      1,
		"server_queue.redis-index_index1_cache_6379.in_queue":                       2,
		"server_queue_bytes.redis-index_index1_cache_6379.out_queue_bytes":          80,
		"server_queue_bytes.redis-index_index1_cache_6379.in_queue_bytes":           50,
		"server_communications.redis-index_index1_cache_6379.requests":              3908,
		"server_communications.redis-index_index1_cache_6379.responses":             3908,
		"server_communication_bytes.redis-index_index1_cache_6379.request_bytes":    170558,
		"server_communication_bytes.redis-index_index1_cache_6379.response_bytes":   176918,
		"pool_error.redis_budget.client_err":                                        30,
		"pool_error.redis_budget.server_ejects":                                     40,
		"pool_error.redis_budget.forward_error":                                     2,
		"pool_client_connections.redis_budget.client_eof":                           2716,
		"pool_client_connections.redis_budget.client_connections":                   221,
		"server_error.redis_budget_budget1_cache_6379.server_err":                   3,
		"server_error.redis_budget_budget1_cache_6379.server_timedout":              4,
		"server_connections.redis_budget_budget1_cache_6379.server_eof":             3,
		"server_connections.redis_budget_budget1_cache_6379.server_connections":     5,
		"server_queue.redis_budget_budget1_cache_6379.out_queue":                    2,
		"server_queue.redis_budget_budget1_cache_6379.in_queue":                     3,
		"server_queue_bytes.redis_budget_budget1_cache_6379.out_queue_bytes":        81,
		"server_queue_bytes.redis_budget_budget1_cache_6379.in_queue_bytes":         51,
		"server_communications.redis_budget_budget1_cache_6379.requests":            3909,
		"server_communications.redis_budget_budget1_cache_6379.responses":           3909,
		"server_communication_bytes.redis_budget_budget1_cache_6379.request_bytes":  170559,
		"server_communication_bytes.redis_budget_budget1_cache_6379.response_bytes": 176919,
		"server_error.redis_budget_budget2_cache_6379.server_err":                   5,
		"server_error.redis_budget_budget2_cache_6379.server_timedout":              6,
		"server_connections.redis_budget_budget2_cache_6379.server_eof":             5,
		"server_connections.redis_budget_budget2_cache_6379.server_connections":     7,
		"server_queue.redis_budget_budget2_cache_6379.out_queue":                    4,
		"server_queue.redis_budget_budget2_cache_6379.in_queue":                     5,
		"server_queue_bytes.redis_budget_budget2_cache_6379.out_queue_bytes":        83,
		"server_queue_bytes.redis_budget_budget2_cache_6379.in_queue_bytes":         53,
		"server_communications.redis_budget_budget2_cache_6379.requests":            3911,
		"server_communications.redis_budget_budget2_cache_6379.responses":           3911,
		"server_communication_bytes.redis_budget_budget2_cache_6379.request_bytes":  170561,
		"server_communication_bytes.redis_budget_budget2_cache_6379.response_bytes": 176921,
	}

	for k, v := range expected {
		value, ok := metrics[k]
		if !ok {
			t.Errorf("metric of %s cannot be fetched", k)
			continue
		}
		if v != value {
			t.Errorf("metric of %s should be %v, but %v", k, v, value)
		}
	}
}

func TestFetchMetricsFail(t *testing.T) {
	assertPanic := func(t *testing.T, f func() (map[string]interface{}, error)) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("FetchMetrics should be panic: stats=%v", stats)
			}
		}()
		f()
	}

	p := TwemproxyPlugin{
		Address: "localhost:" + strconv.Itoa(statsServer.Port()),
		Prefix:  "twemproxy",
		Timeout: 5,
	}

	// panic against a lacking stats json
	noTotalConnectionsJSONStr := strings.Replace(
		jsonStr, "\"total_connections\": 3895,\n", "", 1)
	stats = noTotalConnectionsJSONStr
	assertPanic(t, p.FetchMetrics)

	noClientErrJSONStr := strings.Replace(
		jsonStr, "\"client_err\": 10,\n", "", 1)
	stats = noClientErrJSONStr
	assertPanic(t, p.FetchMetrics)

	noOutQueueJSONStr := strings.Replace(
		jsonStr, "\"out_queue\": 1,\n", "", 1)
	stats = noOutQueueJSONStr
	assertPanic(t, p.FetchMetrics)

	// return error against a stats json with addition
	addInvalidParamJSONStr := strings.Replace(
		jsonStr, "3895,", "3895, \"hoge\": 1,", 1)
	stats = addInvalidParamJSONStr
	_, err := p.FetchMetrics()
	if err == nil {
		t.Errorf("FetchMetrics should return error: stats=%v", stats)
	}
}

func TestGraphDefinition(t *testing.T) {
	p := TwemproxyPlugin{
		Address: "",
		Prefix:  "twemproxy",
		Timeout: 5,
	}
	graphdef := p.GraphDefinition()

	expectedNames := map[string]([]string){
		"connections": {
			"total_connections",
			"curr_connections",
		},
		"pool_error.#": {
			"client_err",
			"server_ejects",
			"forward_error",
		},
		"pool_client_connections.#": {
			"client_connections",
			"client_eof",
		},
		"server_error.#": {
			"server_err",
			"server_timedout",
		},
		"server_connections.#": {
			"server_connections",
			"server_eof",
		},
		"server_queue.#": {
			"out_queue",
			"in_queue",
		},
		"server_queue_bytes.#": {
			"out_queue_bytes",
			"in_queue_bytes",
		},
		"server_communications.#": {
			"requests",
			"responses",
		},
		"server_communication_bytes.#": {
			"request_bytes",
			"response_bytes",
		},
	}

	expectedLabels := map[string]([]string){
		"connections": {
			"New Connections",
			"Current Connections",
		},
		"pool_error.#": {
			"Client Error",
			"Server Ejects",
			"Forward Error",
		},
		"pool_client_connections.#": {
			"Client Connections",
			"Client EOF",
		},
		"server_error.#": {
			"Server Error",
			"Server Timedout",
		},
		"server_connections.#": {
			"Server Connections",
			"Server EOF",
		},
		"server_queue.#": {
			"Out Queue",
			"In Queue",
		},
		"server_queue_bytes.#": {
			"Out Queue Bytes",
			"In Queue Bytes",
		},
		"server_communications.#": {
			"Requests",
			"Responses",
		},
		"server_communication_bytes.#": {
			"Request Bytes",
			"Response Bytes",
		},
	}

	for k, names := range expectedNames {
		value, ok := graphdef[k]
		if !ok {
			t.Errorf("graphdef of %s cannot be fetched", k)
			continue
		}

		var metricNames []string
		var metricLabels []string
		for _, metric := range value.Metrics {
			metricNames = append(metricNames, metric.Name)
			metricLabels = append(metricLabels, metric.Label)
		}

		sort.Strings(names)
		sort.Strings(metricNames)
		if !reflect.DeepEqual(names, metricNames) {
			t.Errorf("graphdef of %s should contain names %v, but %v",
				k, names, metricNames)
		}

		labels := expectedLabels[k]
		sort.Strings(labels)
		sort.Strings(metricLabels)
		if !reflect.DeepEqual(labels, metricLabels) {
			t.Errorf("graphdef of %s should contain labels %v, but %v",
				k, labels, metricLabels)
		}
	}
}
