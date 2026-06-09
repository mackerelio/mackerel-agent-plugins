package mpnginx

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var nginx NginxPlugin

	graphdef := nginx.GraphDefinition()
	if len(graphdef) != 3 {
		t.Errorf("GetTempfilename: %d should be 3", len(graphdef))
	}
}

var responseStub = `Active connections: 123
server accepts handled requests
 1693613501 1693613501 7996986318
Reading: 66 Writing: 16 Waiting: 41
`

func TestParse(t *testing.T) {
	var nginx NginxPlugin
	stub := responseStub
	nginxStats := bytes.NewBufferString(stub)

	stat, err := nginx.parseStats(nginxStats)
	fmt.Println(stat)
	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stat["writing"]).String(), "float64")
	assert.EqualValues(t, stat["writing"], 16)
	assert.EqualValues(t, reflect.TypeOf(stat["accepts"]).String(), "float64")
	assert.EqualValues(t, stat["accepts"], 1693613501)
}

func TestHTTP(t *testing.T) {
	sv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, responseStub)
		}),
	)
	defer sv.Close()

	var nginx NginxPlugin
	nginx.URI = sv.URL
	stat, err := nginx.FetchMetrics()
	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stat["writing"]).String(), "float64")
	assert.EqualValues(t, stat["writing"], 16)
	assert.EqualValues(t, reflect.TypeOf(stat["accepts"]).String(), "float64")
	assert.EqualValues(t, stat["accepts"], 1693613501)
}

// TestHTTPSInsecure tests that FetchMetrics cannot fetch metrics from HTTPS endpoint when TLS certificate is invalid.
// This is the default behavior of this plugin.
func TestHTTPSWithInvalidCert(t *testing.T) {
	sv := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, responseStub)
		}),
	)
	defer sv.Close()

	var nginx NginxPlugin
	nginx.URI = sv.URL

	stat, err := nginx.FetchMetrics()
	assert.NotNil(t, err)
	assert.Nil(t, stat)
}

// TestHTTPSInsecure tests that FetchMetrics returns metrics when the server has an invalid TLS certificate.
func TestHTTPSInsecure(t *testing.T) {
	sv := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, responseStub)
		}),
	)
	defer sv.Close()

	var nginx NginxPlugin
	nginx.URI = sv.URL
	nginx.TLSSkipVerify = true

	stat, err := nginx.FetchMetrics()
	assert.Nil(t, err)
	assert.NotNil(t, stat)
	assert.EqualValues(t, reflect.TypeOf(stat["writing"]).String(), "float64")
	assert.EqualValues(t, stat["writing"], 16)
	assert.EqualValues(t, reflect.TypeOf(stat["accepts"]).String(), "float64")
	assert.EqualValues(t, stat["accepts"], 1693613501)
}
