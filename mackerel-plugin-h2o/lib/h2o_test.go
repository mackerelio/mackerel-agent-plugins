package mph2o

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var h2o H2OPlugin

	graphdef := h2o.GraphDefinition()
	if len(graphdef) != 17 {
		t.Errorf("GetTempfilename: %d should be 17", len(graphdef))
	}
}

func TestParse(t *testing.T) {
	var h2o H2OPlugin
	stub := `{
 "server-version": "2.3.0-DEV",
 "openssl-version": "LibreSSL 2.4.5",
 "current-time": "01/Dec/2017:08:18:16 +0000",
 "restart-time": "01/Dec/2017:08:13:28 +0000",
 "uptime": 288,
 "generation": null,
 "connections": 1,
 "max-connections": 1024,
 "listeners": 4,
 "worker-threads": 2,
 "num-sessions": 14,
 "requests": [
  {"host": "10.0.2.2", "user": null, "at": "20171201T081816.181999+0000", "method": "GET", "path": "/server-status/json", "query": "", "protocol": "HTTP/2", "referer": null, "user-agent": "curl/7.54.0", "connect-time": "0.065440", "request-header-time": "0", "request-body-time": "0", "request-total-time": "0", "process-time": null, "response-time": null, "connection-id": "13", "ssl.protocol-version": "TLSv1.2", "ssl.session-reused": "0", "ssl.cipher": "ECDHE-RSA-CHACHA20-POLY1305-OLD", "ssl.cipher-bits": "256", "ssl.session-ticket": null, "http1.request-index": null, "http2.stream-id": "1", "http2.priority.received.exclusive": "0", "http2.priority.received.parent": "0", "http2.priority.received.weight": "16", "http2.priority.actual.parent": "0", "http2.priority.actual.weight": "16", "authority": "localhost:8443"}
 ],
 "status-errors.400": 0,
 "status-errors.403": 0,
 "status-errors.404": 2,
 "status-errors.405": 0,
 "status-errors.416": 0,
 "status-errors.417": 0,
 "status-errors.500": 0,
 "status-errors.502": 0,
 "status-errors.503": 0,
 "http2-errors.protocol": 0, 
 "http2-errors.internal": 0, 
 "http2-errors.flow-control": 0, 
 "http2-errors.settings-timeout": 0, 
 "http2-errors.stream-closed": 0, 
 "http2-errors.frame-size": 0, 
 "http2-errors.refused-stream": 0, 
 "http2-errors.cancel": 0, 
 "http2-errors.compression": 0, 
 "http2-errors.connect": 0, 
 "http2-errors.enhance-your-calm": 0, 
 "http2-errors.inadequate-security": 0, 
 "http2.read-closed": 3, 
 "http2.write-closed": 0
,
 "connect-time-0": 0,
 "connect-time-25": 0,
 "connect-time-50": 0,
 "connect-time-75": 0,
 "connect-time-99": 0
, "header-time-0": 0,
 "header-time-25": 0,
 "header-time-50": 0,
 "header-time-75": 0,
 "header-time-99": 0
, "body-time-0": 0,
 "body-time-25": 0,
 "body-time-50": 0,
 "body-time-75": 0,
 "body-time-99": 0
, "request-total-time-0": 0,
 "request-total-time-25": 0,
 "request-total-time-50": 0,
 "request-total-time-75": 0,
 "request-total-time-99": 0
, "process-time-0": 0,
 "process-time-25": 0,
 "process-time-50": 0,
 "process-time-75": 0,
 "process-time-99": 0
, "response-time-0": 0,
 "response-time-25": 0,
 "response-time-50": 0,
 "response-time-75": 0,
 "response-time-99": 0
, "duration-0": 0,
 "duration-25": 0,
 "duration-50": 0,
 "duration-75": 0,
 "duration-99": 0
,
 "requests": [
  {"host": "10.0.2.2", "user": null, "at": "20171201T081816.181999+0000", "method": "GET", "path": "/server-status/json", "query": "", "protocol": "HTTP/2", "referer": null, "user-agent": "curl/7.54.0", "connect-time": "0.065440", "request-header-time": "0", "request-body-time": "0", "request-total-time": "0", "process-time": null, "response-time": null, "connection-id": "13", "ssl.protocol-version": "TLSv1.2", "ssl.session-reused": "0", "ssl.cipher": "ECDHE-RSA-CHACHA20-POLY1305-OLD", "ssl.cipher-bits": "256", "ssl.session-ticket": null, "http1.request-index": null, "http2.stream-id": "1", "http2.priority.received.exclusive": "0", "http2.priority.received.parent": "0", "http2.priority.received.weight": "16", "http2.priority.actual.parent": "0", "http2.priority.actual.weight": "16", "authority": "localhost:8443"}
 ],
 "status-errors.400": 0,
 "status-errors.403": 0,
 "status-errors.404": 2,
 "status-errors.405": 0,
 "status-errors.416": 0,
 "status-errors.417": 0,
 "status-errors.500": 0,
 "status-errors.502": 0,
 "status-errors.503": 0,
 "http2-errors.protocol": 0, 
 "http2-errors.internal": 0, 
 "http2-errors.flow-control": 0, 
 "http2-errors.settings-timeout": 0, 
 "http2-errors.stream-closed": 0, 
 "http2-errors.frame-size": 0, 
 "http2-errors.refused-stream": 0, 
 "http2-errors.cancel": 0, 
 "http2-errors.compression": 0, 
 "http2-errors.connect": 0, 
 "http2-errors.enhance-your-calm": 0, 
 "http2-errors.inadequate-security": 0, 
 "http2.read-closed": 3, 
 "http2.write-closed": 0
,
 "connect-time-0": 0,
 "connect-time-25": 0,
 "connect-time-50": 0,
 "connect-time-75": 0,
 "connect-time-99": 0
, "header-time-0": 0,
 "header-time-25": 0,
 "header-time-50": 0,
 "header-time-75": 0,
 "header-time-99": 0
, "body-time-0": 0,
 "body-time-25": 0,
 "body-time-50": 0,
 "body-time-75": 0,
 "body-time-99": 0
, "request-total-time-0": 0,
 "request-total-time-25": 0,
 "request-total-time-50": 0,
 "request-total-time-75": 0,
 "request-total-time-99": 0
, "process-time-0": 0,
 "process-time-25": 0,
 "process-time-50": 0,
 "process-time-75": 0,
 "process-time-99": 0
, "response-time-0": 0,
 "response-time-25": 0,
 "response-time-50": 0,
 "response-time-75": 0,
 "response-time-99": 0
, "duration-0": 0,
 "duration-25": 0,
 "duration-50": 0,
 "duration-75": 0,
 "duration-99": 0
}`

	h2oStats := strings.NewReader(stub)

	stat, err := h2o.parseStats(h2oStats)
	fmt.Println(stat)
	assert.Nil(t, err)
	assert.EqualValues(t, 288, stat["uptime"])
	assert.EqualValues(t, 1, stat["requests"])
	assert.EqualValues(t, 1, stat["connections"])
	assert.EqualValues(t, 2, stat["status-errors_404"])
	assert.EqualValues(t, 3, stat["http2_read-closed"])
	assert.EqualValues(t, 0, stat["connect-time-25"])
}
