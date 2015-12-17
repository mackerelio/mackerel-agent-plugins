package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var rack RackStatsPlugin

	graphdef := rack.GraphDefinition()
	if len(graphdef) != 1 {
		t.Errorf("GetTempfilename: %d should be 1", len(graphdef))
	}
}

var testPort string
var testSock string

func requestHandlerHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "calling: 10")
	fmt.Fprintln(w, "writing: 20")
	fmt.Fprintln(w, fmt.Sprintf("0.0.0.0:%s active: 31", testPort))
	fmt.Fprintln(w, fmt.Sprintf("0.0.0.0:%s queued: 40", testPort))
	fmt.Fprintln(w, fmt.Sprintf("0.0.0.0:%s0 active: 50", testPort))
	fmt.Fprintln(w, fmt.Sprintf("0.0.0.0:%s0 queued: 60", testPort))
	fmt.Fprintln(w, "/path/to/unicorn.sock active: 70")
	fmt.Fprintln(w, "/path/to/unicorn.sock queued: 80")
	fmt.Fprintln(w, "/path/to/unicorn2.sock active: 90")
	fmt.Fprintln(w, "/path/to/unicorn2.sock queued: 100")
}

func requestHandlerUnix(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "calling: 10")
	fmt.Fprintln(w, "writing: 20")
	fmt.Fprintln(w, "0.0.0.0:8080 active: 50")
	fmt.Fprintln(w, "0.0.0.0:8080 queued: 60")
	fmt.Fprintln(w, fmt.Sprintf("%s active: 71", testSock))
	fmt.Fprintln(w, fmt.Sprintf("%s queued: 80", testSock))
	fmt.Fprintln(w, "/path/to/unicorn2.sock active: 90")
	fmt.Fprintln(w, "/path/to/unicorn2.sock queued: 100")
}

func TestParseHttp(t *testing.T) {
	var rack RackStatsPlugin

	ts := httptest.NewServer(http.HandlerFunc(requestHandlerHTTP))
	defer ts.Close()

	u, _ := url.Parse(ts.URL)

	var host string
	host, testPort, _ = net.SplitHostPort(u.Host)

	rack.Address = fmt.Sprintf("http://%s:%s", host, testPort)
	rack.Path = "/_raindrops"

	stats, err := rack.parseStats()
	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stats["calling"]).String(), "float64")
	assert.EqualValues(t, stats["calling"], 10)
	assert.EqualValues(t, reflect.TypeOf(stats["writing"]).String(), "float64")
	assert.EqualValues(t, stats["writing"], 20)
	assert.EqualValues(t, reflect.TypeOf(stats["active"]).String(), "float64")
	assert.EqualValues(t, stats["active"], 30)
	assert.EqualValues(t, reflect.TypeOf(stats["queued"]).String(), "float64")
	assert.EqualValues(t, stats["queued"], 40)
}

func TestParseUnix(t *testing.T) {
	var rack RackStatsPlugin

	dir, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	testSock = fmt.Sprintf("%s/unicorn.sock", dir)

	l, err := net.Listen("unix", testSock)
	if err != nil {
		t.Error(err)
	}

	rack.Address = fmt.Sprintf("unix:%s", testSock)
	rack.Path = "/_raindrops"

	var mux = http.NewServeMux()
	mux.Handle(rack.Path, http.HandlerFunc(requestHandlerUnix))

	go http.Serve(l, mux)

	stats, err := rack.parseStats()
	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stats["calling"]).String(), "float64")
	assert.EqualValues(t, stats["calling"], 10)
	assert.EqualValues(t, reflect.TypeOf(stats["writing"]).String(), "float64")
	assert.EqualValues(t, stats["writing"], 20)
	assert.EqualValues(t, reflect.TypeOf(stats["active"]).String(), "float64")
	assert.EqualValues(t, stats["active"], 70)
	assert.EqualValues(t, reflect.TypeOf(stats["queued"]).String(), "float64")
	assert.EqualValues(t, stats["queued"], 80)
}
