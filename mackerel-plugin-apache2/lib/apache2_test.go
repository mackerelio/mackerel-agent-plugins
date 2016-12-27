package mpapache2

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseApache2Scoreboard(t *testing.T) {
	stub := "Scoreboard: W._SRWKDCLGI...."
	stat := make(map[string]interface{})

	err := parseApache2Scoreboard(stub, &stat)
	assert.Nil(t, err)
	assert.EqualValues(t, stat["score-_"], 1)
	assert.EqualValues(t, stat["score-S"], 1)
	assert.EqualValues(t, stat["score-R"], 1)
	assert.EqualValues(t, stat["score-W"], 2)
	assert.EqualValues(t, stat["score-K"], 1)
	assert.EqualValues(t, stat["score-D"], 1)
	assert.EqualValues(t, stat["score-C"], 1)
	assert.EqualValues(t, stat["score-L"], 1)
	assert.EqualValues(t, stat["score-G"], 1)
	assert.EqualValues(t, stat["score-I"], 1)
	assert.EqualValues(t, stat["score-"], 5)
}

func TestParseApache2Status(t *testing.T) {
	stub := `Total Accesses: 358
Total kBytes: 20
CPULoad: .00117358
Uptime: 102251
ReqPerSec: .00350119
BytesPerSec: .200291
BytesPerReq: 57.2067
BusyWorkers: 1
IdleWorkers: 4
`
	stat := make(map[string]interface{})

	err := parseApache2Status(stub, &stat)
	assert.Nil(t, err)
	assert.EqualValues(t, stat["requests"], 358)
	assert.EqualValues(t, stat["bytes_sent"], 20)
	assert.EqualValues(t, stat["cpu_load"], 0.00117358)
	assert.EqualValues(t, stat["busy_workers"], 1)
	assert.EqualValues(t, stat["idle_workers"], 4)
}

func TestGetApache2Metrics_1(t *testing.T) {
	stub := `Total Accesses: 668
Total kBytes: 2789
CPULoad: .000599374
Uptime: 171846
ReqPerSec: .0038872
BytesPerSec: 16.6192
BytesPerReq: 4275.35
BusyWorkers: 1
IdleWorkers: 3
Scoreboard: W_.__...........................`

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, stub)
			}))
	defer ts.Close()
	re, _ := regexp.Compile("([a-z]+)://([A-Za-z0-9.]+):([0-9]+)(.*)")
	found := re.FindStringSubmatch(ts.URL)
	assert.EqualValues(t, len(found), 5, fmt.Sprintf("Test stub uri format is changed. %s", ts.URL))

	host := found[2]
	port, _ := strconv.Atoi(found[3])
	path := found[4]
	header := []string{fmt.Sprintf("Host: %s", found[2]), "X-Text-Header: test"}

	ret, err := getApache2Metrics(host, uint16(port), path, header)
	assert.Nil(t, err)
	assert.NotNil(t, ret)
	assert.NotEmpty(t, ret)
	assert.Contains(t, ret, "Total Accesses")
	assert.Contains(t, ret, "Total kBytes")
	assert.Contains(t, ret, "Uptime")
	assert.Contains(t, ret, "BusyWorkers")
	assert.Contains(t, ret, "IdleWorkers")
	assert.Contains(t, ret, "Scoreboard")
}
