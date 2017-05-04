package mpsolr

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/solr/admin/cores":
		fmt.Fprintf(w, fetchJSON("cores"))
	case "/solr/testcore/admin/mbeans":
		key := r.URL.Query()["cat"][0]
		fmt.Fprintf(w, fetchJSONForMbeans(key))
	default:
		fmt.Fprintf(w, "")
	}
})

func fetchJSON(name string) string {
	json, err := ioutil.ReadFile("./stats/6.4.2/" + name + ".json")
	if err != nil {
		panic(err)
	}
	return string(json)
}

func fetchJSONForMbeans(key string) string {
	switch key {
	case "QUERYHANDLER":
		return fetchJSON("queryhandler")
	case "CACHE":
		return fetchJSON("cache")
	default:
		return ""
	}
}

func TestGraphDefinition(t *testing.T) {
	solr := SolrPlugin{
		Cores:  []string{"testcore"},
		Prefix: "solr",
	}
	graphdef := solr.GraphDefinition()

	assert.EqualValues(t, "testcore DocsCount", graphdef["solr.testcore.docsCount"].Label)
	assert.EqualValues(t, "testcore_numDocs", graphdef["solr.testcore.docsCount"].Metrics[0].Name)
	assert.EqualValues(t, "testcore_deletedDocs", graphdef["solr.testcore.docsCount"].Metrics[1].Name)

	assert.EqualValues(t, "testcore IndexHeapUsageBytes", graphdef["solr.testcore.indexHeapUsageBytes"].Label)
	assert.EqualValues(t, "testcore_indexHeapUsageBytes", graphdef["solr.testcore.indexHeapUsageBytes"].Metrics[0].Name)
	assert.EqualValues(t, "testcore SegmentCount", graphdef["solr.testcore.segmentCount"].Label)
	assert.EqualValues(t, "testcore_segmentCount", graphdef["solr.testcore.segmentCount"].Metrics[0].Name)
	assert.EqualValues(t, "testcore SizeInBytes", graphdef["solr.testcore.sizeInBytes"].Label)
	assert.EqualValues(t, "testcore_sizeInBytes", graphdef["solr.testcore.sizeInBytes"].Metrics[0].Name)

	assert.EqualValues(t, "testcore Requests", graphdef["solr.testcore.requests"].Label)
	assert.EqualValues(t, "testcore_requests_updatejson", graphdef["solr.testcore.requests"].Metrics[0].Name)
	assert.EqualValues(t, "testcore_requests_select", graphdef["solr.testcore.requests"].Metrics[1].Name)
	assert.EqualValues(t, "testcore_requests_updatejsondocs", graphdef["solr.testcore.requests"].Metrics[2].Name)
	assert.EqualValues(t, "testcore_requests_get", graphdef["solr.testcore.requests"].Metrics[3].Name)
	assert.EqualValues(t, "testcore_requests_updatecsv", graphdef["solr.testcore.requests"].Metrics[4].Name)
	assert.EqualValues(t, "testcore_requests_replication", graphdef["solr.testcore.requests"].Metrics[5].Name)
	assert.EqualValues(t, "testcore_requests_update", graphdef["solr.testcore.requests"].Metrics[6].Name)
	assert.EqualValues(t, "testcore_requests_dataimport", graphdef["solr.testcore.requests"].Metrics[7].Name)

	assert.EqualValues(t, "testcore Errors", graphdef["solr.testcore.errors"].Label)
	assert.EqualValues(t, "testcore Timeouts", graphdef["solr.testcore.timeouts"].Label)
	assert.EqualValues(t, "testcore AvgRequestsPerSecond", graphdef["solr.testcore.avgRequestsPerSecond"].Label)
	assert.EqualValues(t, "testcore 5minRateRequestsPerSecond", graphdef["solr.testcore.5minRateRequestsPerSecond"].Label)
	assert.EqualValues(t, "testcore 15minRateRequestsPerSecond", graphdef["solr.testcore.15minRateRequestsPerSecond"].Label)
	assert.EqualValues(t, "testcore AvgTimePerRequest", graphdef["solr.testcore.avgTimePerRequest"].Label)
	assert.EqualValues(t, "testcore MedianRequestTime", graphdef["solr.testcore.medianRequestTime"].Label)
	assert.EqualValues(t, "testcore 75thPcRequestTime", graphdef["solr.testcore.75thPcRequestTime"].Label)
	assert.EqualValues(t, "testcore 95thPcRequestTime", graphdef["solr.testcore.95thPcRequestTime"].Label)
	assert.EqualValues(t, "testcore 99thPcRequestTime", graphdef["solr.testcore.99thPcRequestTime"].Label)
	assert.EqualValues(t, "testcore 999thPcRequestTime", graphdef["solr.testcore.999thPcRequestTime"].Label)

	assert.EqualValues(t, "testcore Lookups", graphdef["solr.testcore.lookups"].Label)
	assert.EqualValues(t, "testcore_lookups_filterCache", graphdef["solr.testcore.lookups"].Metrics[0].Name)
	assert.EqualValues(t, "testcore_lookups_perSegFilter", graphdef["solr.testcore.lookups"].Metrics[1].Name)
	assert.EqualValues(t, "testcore_lookups_queryResultCache", graphdef["solr.testcore.lookups"].Metrics[2].Name)
	assert.EqualValues(t, "testcore_lookups_documentCache", graphdef["solr.testcore.lookups"].Metrics[3].Name)
	assert.EqualValues(t, "testcore_lookups_fieldValueCache", graphdef["solr.testcore.lookups"].Metrics[4].Name)

	assert.EqualValues(t, "testcore Hits", graphdef["solr.testcore.hits"].Label)
	assert.EqualValues(t, "testcore Hitratio", graphdef["solr.testcore.hitratio"].Label)
	assert.EqualValues(t, "testcore Inserts", graphdef["solr.testcore.inserts"].Label)
	assert.EqualValues(t, "testcore Evictions", graphdef["solr.testcore.evictions"].Label)
	assert.EqualValues(t, "testcore Size", graphdef["solr.testcore.size"].Label)
	assert.EqualValues(t, "testcore WarmupTime", graphdef["solr.testcore.warmupTime"].Label)
}

func TestFetchMetrics(t *testing.T) {
	ts := httptest.NewServer(testHandler)
	defer ts.Close()

	solr := SolrPlugin{
		BaseURL: ts.URL + "/solr",
		Cores:   []string{"testcore"},
		Prefix:  "solr",
	}
	solr.loadStats()
	stat, err := solr.FetchMetrics()
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 12345, stat["testcore_numDocs"])
	assert.EqualValues(t, 13, stat["testcore_deletedDocs"])

	assert.EqualValues(t, 64, stat["testcore_indexHeapUsageBytes"])
	assert.EqualValues(t, 7, stat["testcore_segmentCount"])
	assert.EqualValues(t, 71, stat["testcore_sizeInBytes"])

	assert.EqualValues(t, 777.0, stat["testcore_requests_select"])
	assert.EqualValues(t, 14.0, stat["testcore_errors_select"])
	assert.EqualValues(t, 3.0, stat["testcore_timeouts_select"])
	assert.EqualValues(t, 5.0, stat["testcore_avgRequestsPerSecond_select"])
	assert.EqualValues(t, 4.0, stat["testcore_5minRateRequestsPerSecond_select"])
	assert.EqualValues(t, 7.0, stat["testcore_15minRateRequestsPerSecond_select"])
	assert.EqualValues(t, 0.3, stat["testcore_avgTimePerRequest_select"])
	assert.EqualValues(t, 0.4, stat["testcore_medianRequestTime_select"])
	assert.EqualValues(t, 0.5, stat["testcore_75thPcRequestTime_select"])
	assert.EqualValues(t, 0.6, stat["testcore_95thPcRequestTime_select"])
	assert.EqualValues(t, 0.7, stat["testcore_99thPcRequestTime_select"])
	assert.EqualValues(t, 0.8, stat["testcore_999thPcRequestTime_select"])

	assert.EqualValues(t, 1, stat["testcore_lookups_filterCache"])
	assert.EqualValues(t, 2, stat["testcore_lookups_perSegFilter"])
	assert.EqualValues(t, 3, stat["testcore_lookups_queryResultCache"])
	assert.EqualValues(t, 4, stat["testcore_lookups_documentCache"])
	assert.EqualValues(t, 5, stat["testcore_lookups_fieldValueCache"])
}
