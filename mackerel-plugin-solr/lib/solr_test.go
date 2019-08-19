package mpsolr

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	SolrVersions = []string{"5.5.4", "6.4.2", "7.7.2", "8.1.1"}
	solrVersion  string
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/solr/admin/info/system":
		fmt.Fprintf(w, fetchJSON("system"))
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
	path := fmt.Sprintf("./stats/%s/%s.json", solrVersion, name)
	json, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(json)
}

func fetchJSONForMbeans(key string) string {
	switch key {
	case "QUERYHANDLER", "QUERY":
		return fetchJSON("query")
	case "UPDATEHANDLER", "UPDATE":
		return fetchJSON("update")
	case "REPLICATION":
		return fetchJSON("replication")
	case "CACHE":
		return fetchJSON("cache")
	default:
		return ""
	}
}

func setupSolr(mockURL string, version string) (solr SolrPlugin) {
	solrVersion = version
	solr = SolrPlugin{
		BaseURL: mockURL + "/solr",
		Cores:   []string{"testcore"},
		Prefix:  "solr",
	}
	solr.loadVersion()
	return
}

func TestGraphDefinition(t *testing.T) {
	ts := httptest.NewServer(testHandler)
	defer ts.Close()

	for _, version := range SolrVersions {
		solr := setupSolr(ts.URL, version)
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
		if version == "7.7.2" || version == "8.1.1" {
			assert.EqualValues(t, "testcore ClientErrors", graphdef["solr.testcore.clientErrors"].Label)
			assert.EqualValues(t, "testcore ServerErrors", graphdef["solr.testcore.serverErrors"].Label)
			assert.EqualValues(t, "testcore RequestTimes", graphdef["solr.testcore.requestTimes"].Label)
		} else {
			assert.EqualValues(t, "testcore AvgRequestsPerSecond", graphdef["solr.testcore.avgRequestsPerSecond"].Label)
			assert.EqualValues(t, "testcore 5minRateRequestsPerSecond", graphdef["solr.testcore.5minRateRequestsPerSecond"].Label)
			assert.EqualValues(t, "testcore 15minRateRequestsPerSecond", graphdef["solr.testcore.15minRateRequestsPerSecond"].Label)
			assert.EqualValues(t, "testcore AvgTimePerRequest", graphdef["solr.testcore.avgTimePerRequest"].Label)
			assert.EqualValues(t, "testcore MedianRequestTime", graphdef["solr.testcore.medianRequestTime"].Label)
			assert.EqualValues(t, "testcore 75thPcRequestTime", graphdef["solr.testcore.75thPcRequestTime"].Label)
			assert.EqualValues(t, "testcore 95thPcRequestTime", graphdef["solr.testcore.95thPcRequestTime"].Label)
			assert.EqualValues(t, "testcore 99thPcRequestTime", graphdef["solr.testcore.99thPcRequestTime"].Label)
			assert.EqualValues(t, "testcore 999thPcRequestTime", graphdef["solr.testcore.999thPcRequestTime"].Label)
		}

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
}

func TestFetchMetrics(t *testing.T) {
	ts := httptest.NewServer(testHandler)
	defer ts.Close()

	for _, version := range SolrVersions {
		solr := setupSolr(ts.URL, version)
		solr.loadStats()
		stat, err := solr.FetchMetrics()
		if err != nil {
			t.Fatal(err)
		}

		msgPrefix := "Solr Version: " + solrVersion + ", Stat Key: "

		assert.EqualValues(t, 12345, stat["testcore_numDocs"], msgPrefix+"testcore_numDocs")
		assert.EqualValues(t, 54321, stat["testcore_deletedDocs"], msgPrefix+"testcore_deletedDocs")

		assert.EqualValues(t, 128, stat["testcore_indexHeapUsageBytes"], msgPrefix+"testcore_indexHeapUsageBytes")
		assert.EqualValues(t, 256, stat["testcore_segmentCount"], msgPrefix+"testcore_segmentCount")
		assert.EqualValues(t, 512, stat["testcore_sizeInBytes"], msgPrefix+"testcore_sizeInBytes")

		assert.EqualValues(t, 111.0, stat["testcore_requests_select"], msgPrefix+"testcore_requests_select")
		assert.EqualValues(t, 222.0, stat["testcore_errors_select"], msgPrefix+"testcore_errors_select")
		assert.EqualValues(t, 333.0, stat["testcore_timeouts_select"], msgPrefix+"testcore_timeouts_select")
		if version == "7.7.2" || version == "8.1.1" {
			assert.EqualValues(t, 432.0, stat["testcore_clientErrors_select"], msgPrefix+"testcore_clientErrors_select")
			assert.EqualValues(t, 234.0, stat["testcore_serverErrors_select"], msgPrefix+"testcore_serverErrors_select")
			assert.EqualValues(t, 0.123, stat["testcore_requestTimes_select"], msgPrefix+"testcore_requestTimes_select")
		}
		if version == "5.5.4" || version == "6.4.2" {
			assert.EqualValues(t, 444.0, stat["testcore_avgRequestsPerSecond_select"], msgPrefix+"testcore_avgRequestsPerSecond_select")
		}
		switch version {
		case "5.5.4":
			assert.EqualValues(t, 555.0, stat["testcore_5minRateReqsPerSecond_select"], msgPrefix+"testcore_5minRateRequestsPerSecond_select")
			assert.EqualValues(t, 666.0, stat["testcore_15minRateReqsPerSecond_select"], msgPrefix+"testcore_15minRateRequestsPerSecond_select")
		case "6.4.2":
			assert.EqualValues(t, 555.0, stat["testcore_5minRateRequestsPerSecond_select"], msgPrefix+"testcore_5minRateRequestsPerSecond_select")
			assert.EqualValues(t, 666.0, stat["testcore_15minRateRequestsPerSecond_select"], msgPrefix+"testcore_15minRateRequestsPerSecond_select")
		default:
		}
		if version == "5.5.4" || version == "6.4.2" {
			assert.EqualValues(t, 777.0, stat["testcore_avgTimePerRequest_select"], msgPrefix+"testcore_avgTimePerRequest_select")
			assert.EqualValues(t, 888.0, stat["testcore_medianRequestTime_select"], msgPrefix+"testcore_medianRequestTime_select")
			assert.EqualValues(t, 999.0, stat["testcore_75thPcRequestTime_select"], msgPrefix+"testcore_75thPcRequestTime_select")
			assert.EqualValues(t, 100.1, stat["testcore_95thPcRequestTime_select"], msgPrefix+"testcore_95thPcRequestTime_select")
			assert.EqualValues(t, 100.2, stat["testcore_99thPcRequestTime_select"], msgPrefix+"testcore_99thPcRequestTime_select")
			assert.EqualValues(t, 100.3, stat["testcore_999thPcRequestTime_select"], msgPrefix+"testcore_999thPcRequestTime_select")
		}

		assert.EqualValues(t, 1111.0, stat["testcore_requests_update"], msgPrefix+"testcore_requests_update")
		assert.EqualValues(t, 2222.0, stat["testcore_errors_update"], msgPrefix+"testcore_errors_update")
		assert.EqualValues(t, 3333.0, stat["testcore_timeouts_update"], msgPrefix+"testcore_timeouts_update")
		if version == "7.7.2" || version == "8.1.1" {
			assert.EqualValues(t, 4321.0, stat["testcore_clientErrors_update"], msgPrefix+"testcore_clientErrors_update")
			assert.EqualValues(t, 1234.0, stat["testcore_serverErrors_update"], msgPrefix+"testcore_serverErrors_update")
			assert.EqualValues(t, 0.1234, stat["testcore_requestTimes_update"], msgPrefix+"testcore_requestTimes_update")
		}
		if version == "5.5.4" || version == "6.4.2" {
			assert.EqualValues(t, 4444.0, stat["testcore_avgRequestsPerSecond_update"], msgPrefix+"testcore_avgRequestsPerSecond_update")
		}
		switch version {
		case "5.5.4":
			assert.EqualValues(t, 5555.0, stat["testcore_5minRateReqsPerSecond_update"], msgPrefix+"testcore_5minRateRequestsPerSecond_update")
			assert.EqualValues(t, 6666.0, stat["testcore_15minRateReqsPerSecond_update"], msgPrefix+"testcore_15minRateRequestsPerSecond_update")
		case "6.4.2":
			assert.EqualValues(t, 5555.0, stat["testcore_5minRateRequestsPerSecond_update"], msgPrefix+"testcore_5minRateRequestsPerSecond_update")
			assert.EqualValues(t, 6666.0, stat["testcore_15minRateRequestsPerSecond_update"], msgPrefix+"testcore_15minRateRequestsPerSecond_update")
		default:
		}
		if version == "5.5.4" || version == "6.4.2" {
			assert.EqualValues(t, 7777.0, stat["testcore_avgTimePerRequest_update"], msgPrefix+"testcore_avgTimePerRequest_update")
			assert.EqualValues(t, 8888.0, stat["testcore_medianRequestTime_update"], msgPrefix+"testcore_medianRequestTime_update")
			assert.EqualValues(t, 9999.0, stat["testcore_75thPcRequestTime_update"], msgPrefix+"testcore_75thPcRequestTime_update")
			assert.EqualValues(t, 1000.1, stat["testcore_95thPcRequestTime_update"], msgPrefix+"testcore_95thPcRequestTime_update")
			assert.EqualValues(t, 1000.2, stat["testcore_99thPcRequestTime_update"], msgPrefix+"testcore_99thPcRequestTime_update")
			assert.EqualValues(t, 1000.3, stat["testcore_999thPcRequestTime_update"], msgPrefix+"testcore_999thPcRequestTime_update")
		}

		assert.EqualValues(t, 11111.0, stat["testcore_requests_replication"], msgPrefix+"testcore_requests_replication")
		assert.EqualValues(t, 22222.0, stat["testcore_errors_replication"], msgPrefix+"testcore_errors_replication")
		assert.EqualValues(t, 33333.0, stat["testcore_timeouts_replication"], msgPrefix+"testcore_timeouts_replication")
		if version == "7.7.2" || version == "8.1.1" {
			assert.EqualValues(t, 54321.0, stat["testcore_clientErrors_replication"], msgPrefix+"testcore_clientErrors_replication")
			assert.EqualValues(t, 12345.0, stat["testcore_serverErrors_replication"], msgPrefix+"testcore_serverErrors_replication")
			assert.EqualValues(t, 0.12345, stat["testcore_requestTimes_replication"], msgPrefix+"testcore_requestTimes_replication")
		}
		if version == "5.5.4" || version == "6.4.2" {
			assert.EqualValues(t, 44444.0, stat["testcore_avgRequestsPerSecond_replication"], msgPrefix+"testcore_avgRequestsPerSecond_replication")
		}
		switch version {
		case "5.5.4":
			assert.EqualValues(t, 55555.0, stat["testcore_5minRateReqsPerSecond_replication"], msgPrefix+"testcore_5minRateRequestsPerSecond_replication")
			assert.EqualValues(t, 66666.0, stat["testcore_15minRateReqsPerSecond_replication"], msgPrefix+"testcore_15minRateRequestsPerSecond_replication")
		case "6.4.2":
			assert.EqualValues(t, 55555.0, stat["testcore_5minRateRequestsPerSecond_replication"], msgPrefix+"testcore_5minRateRequestsPerSecond_replication")
			assert.EqualValues(t, 66666.0, stat["testcore_15minRateRequestsPerSecond_replication"], msgPrefix+"testcore_15minRateRequestsPerSecond_replication")
		default:
		}
		if version == "5.5.4" || version == "6.4.2" {
			assert.EqualValues(t, 77777.0, stat["testcore_avgTimePerRequest_replication"], msgPrefix+"testcore_avgTimePerRequest_replication")
			assert.EqualValues(t, 88888.0, stat["testcore_medianRequestTime_replication"], msgPrefix+"testcore_medianRequestTime_replication")
			assert.EqualValues(t, 99999.0, stat["testcore_75thPcRequestTime_replication"], msgPrefix+"testcore_75thPcRequestTime_replication")
			assert.EqualValues(t, 10000.1, stat["testcore_95thPcRequestTime_replication"], msgPrefix+"testcore_95thPcRequestTime_replication")
			assert.EqualValues(t, 10000.2, stat["testcore_99thPcRequestTime_replication"], msgPrefix+"testcore_99thPcRequestTime_replication")
			assert.EqualValues(t, 10000.3, stat["testcore_999thPcRequestTime_replication"], msgPrefix+"testcore_999thPcRequestTime_replication")
		}

		assert.EqualValues(t, 135, stat["testcore_lookups_filterCache"], msgPrefix+"testcore_lookups_filterCache")
		assert.EqualValues(t, 357, stat["testcore_lookups_perSegFilter"], msgPrefix+"testcore_lookups_perSegFilter")
		assert.EqualValues(t, 579, stat["testcore_lookups_queryResultCache"], msgPrefix+"testcore_lookups_queryResultCache")
		assert.EqualValues(t, 246, stat["testcore_lookups_documentCache"], msgPrefix+"testcore_lookups_documentCache")
		assert.EqualValues(t, 468, stat["testcore_lookups_fieldValueCache"], msgPrefix+"testcore_lookups_fieldValueCache")
	}
}
