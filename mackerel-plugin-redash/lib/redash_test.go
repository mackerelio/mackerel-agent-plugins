package mpredash

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	statsServer *httptest.Server
	stub        string
)

func TestMain(m *testing.M) {
	os.Exit(mainTest(m))
}

func mainTest(m *testing.M) int {
	flag.Parse()

	// run a test server returning a dummy stats json
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, stub)
			}))
	statsServer = ts

	log.Println("Started a stats server")

	return m.Run()
}

var jsonStr = `{
    "waiting": [
        {
            "username": "Scheduled",
            "retries": 0,
            "scheduled_retries": 0,
            "task_id": "8596cc9a-c42d-4518-a29c-12cb993b2b1e",
            "created_at": 1496649134.345524,
            "updated_at": 1496649148.593866,
            "state": "waiting",
            "query_id": 136,
            "run_time": 14.219444036483765,
            "error": null,
            "scheduled": true,
            "started_at": 1496649134.35855,
            "data_source_id": 2,
            "query_hash": "6e4d43f0dc697d33f632b63ecb55dd0f"
        },
        {
            "username": "FugaFuga",
            "retries": 0,
            "scheduled_retries": 0,
            "task_id": "8596cc9a-c42d-4518-a29c-12cb993b2b1e",
            "created_at": 1496649134.345524,
            "updated_at": 1496649148.593866,
            "state": "unknown",
            "query_id": 136,
            "run_time": 14.219444036483765,
            "error": null,
            "scheduled": false,
            "started_at": 1496649134.35855,
            "data_source_id": 2,
            "query_hash": "6e4d43f0dc697d33f632b63ecb55dd0f"
        }
    ],
    "done": [
        {
            "username": "Scheduled",
            "retries": 0,
            "scheduled_retries": 0,
            "task_id": "8596cc9a-c42d-4518-a29c-12cb993b2b1e",
            "created_at": 1496649134.345524,
            "updated_at": 1496649148.593866,
            "state": "finished",
            "query_id": 136,
            "run_time": 14.219444036483765,
            "error": null,
            "scheduled": true,
            "started_at": 1496649134.35855,
            "data_source_id": 2,
            "query_hash": "6e4d43f0dc697d33f632b63ecb55dd0f"
        },
        {
            "username": "Scheduled",
            "retries": 0,
            "scheduled_retries": 0,
            "task_id": "8596cc9a-c42d-4518-a29c-12cb993b2b1e",
            "created_at": 1496649134.345524,
            "updated_at": 1496649148.593866,
            "state": "finished",
            "query_id": 137,
            "run_time": 14.219444036483765,
            "error": null,
            "scheduled": true,
            "started_at": 1496649134.35855,
            "data_source_id": 3,
            "query_hash": "6e4d43f0dc697d33f632b63ecb55dd0f"
        },
        {
            "username": "Hogehoge",
            "retries": 0,
            "scheduled_retries": 0,
            "task_id": "8596cc9a-c42d-4518-a29c-12cb993b2b1e",
            "created_at": 1496649134.345524,
            "updated_at": 1496649148.593866,
            "state": "failed",
            "query_id": 136,
            "run_time": 14.219444036483765,
            "error": "You have an error in your SQL syntax",
            "scheduled": false,
            "started_at": 1496649134.35855,
            "data_source_id": 2,
            "query_hash": "6e4d43f0dc697d33f632b63ecb55dd0f"
        }
    ],
    "in_progress": []
}`

func TestFetchMetrics(t *testing.T) {
	// response a valid stats json
	stub = jsonStr

	// get metrics
	p := RedashPlugin{
		URI:     statsServer.URL,
		Prefix:  "redash",
		Timeout: 5,
	}
	metrics, err := p.FetchMetrics()
	if err != nil {
		t.Errorf("Failed to FetchMetrics: %s", err)
		return
	}

	// check the metrics
	expected := map[string]uint64{
		"wait":                                       2,
		"done":                                       3,
		"in_progress":                                0,
		"task_scheduled_count.wait.adhoc":            1,
		"task_scheduled_count.done.adhoc":            1,
		"task_scheduled_count.in_progress.adhoc":     0,
		"task_scheduled_count.wait.scheduled":        1,
		"task_scheduled_count.done.scheduled":        2,
		"task_scheduled_count.in_progress.scheduled": 0,
		"waiting":         1,
		"finished":        2,
		"executing_query": 0,
		"failed":          1,
		"processing":      0,
		"checking_alerts": 0,
		"other":           1,
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
	p := RedashPlugin{
		URI:     statsServer.URL,
		Prefix:  "redash",
		Timeout: 5,
	}

	// return error against an invalid stats json
	stub = "{waiting: [],}"
	_, err := p.FetchMetrics()
	if err == nil {
		t.Errorf("FetchMetrics should return error: stub=%v", stub)
	}
}
