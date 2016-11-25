package mpelasticsearch

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	json, err := ioutil.ReadFile("./stat.json")
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(json))
})

func TestGraphDefinition(t *testing.T) {
	elasticsearch := ElasticsearchPlugin{
		Prefix:      "elasticsearch",
		LabelPrefix: "Elasticsearch",
	}
	graphdef := elasticsearch.GraphDefinition()

	assert.EqualValues(t, "Elasticsearch HTTP", graphdef["elasticsearch.http"].Label)
	assert.EqualValues(t, "Elasticsearch Thread-Pool Threads", graphdef["elasticsearch.thread_pool.threads"].Label)
	assert.EqualValues(t, "threads_fetch_shard_started", graphdef["elasticsearch.thread_pool.threads"].Metrics[16].Name)
	assert.EqualValues(t, "threads_fetch_shard_store", graphdef["elasticsearch.thread_pool.threads"].Metrics[17].Name)
	assert.EqualValues(t, "threads_listener", graphdef["elasticsearch.thread_pool.threads"].Metrics[18].Name)
}

func TestFetchMetrics(t *testing.T) {
	ts := httptest.NewServer(testHandler)
	defer ts.Close()

	var elasticsearch ElasticsearchPlugin
	elasticsearch.URI = ts.URL
	stat, err := elasticsearch.FetchMetrics()
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 6991, stat["http_opened"])
	assert.EqualValues(t, 1, stat["threads_generic"])
	assert.EqualValues(t, 0, stat["threads_merge"])
	assert.EqualValues(t, 2, stat["threads_fetch_shard_started"])
	assert.EqualValues(t, 3, stat["threads_fetch_shard_store"])
	assert.EqualValues(t, 1, stat["threads_listener"])
}
