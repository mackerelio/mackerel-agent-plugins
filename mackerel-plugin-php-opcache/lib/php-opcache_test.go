package mpphpopcache

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPhpOpcacheStatus_1(t *testing.T) {
	stub := `used_memory:10941328
free_memory:123276400
wasted_memory:0
current_wasted_percentage:0
num_cached_scripts:1
num_cached_keys:2
max_cached_keys:7963
hits:123
oom_restarts:0
hash_restarts:0
manual_restarts:0
misses:53
blacklist_misses:12
blacklist_miss_ratio:10
opcache_hit_rate:15`

	stat := make(map[string]float64)

	err := parsePhpOpcacheStatus(stub, &stat)
	assert.Nil(t, err)
	assert.EqualValues(t, stat["used_memory"], 10941328)
	assert.EqualValues(t, stat["free_memory"], 123276400)
	assert.EqualValues(t, stat["wasted_memory"], 0)
	assert.EqualValues(t, stat["current_wasted_percentage"], 0)
	assert.EqualValues(t, stat["num_cached_scripts"], 1)
	assert.EqualValues(t, stat["num_cached_keys"], 2)
	assert.EqualValues(t, stat["max_cached_keys"], 7963)
	assert.EqualValues(t, stat["hits"], 123)
	assert.EqualValues(t, stat["oom_restarts"], 0)
	assert.EqualValues(t, stat["hash_restarts"], 0)
	assert.EqualValues(t, stat["manual_restarts"], 0)
	assert.EqualValues(t, stat["misses"], 53)
	assert.EqualValues(t, stat["blacklist_misses"], 12)
	assert.EqualValues(t, stat["blacklist_miss_ratio"], 10)
	assert.EqualValues(t, stat["opcache_hit_rate"], 15)
}

func TestGetPhpOpcacheMetrics_1(t *testing.T) {
	stub := `used_memory:10941328
free_memory:123276400
wasted_memory:0
current_wasted_percentage:0
num_cached_scripts:1
num_cached_keys:2
max_cached_keys:7963
hits:0
oom_restarts:0
hash_restarts:0
manual_restarts:0
misses:1
blacklist_misses:0
blacklist_miss_ratio:0
opcache_hit_rate:0`

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

	ret, err := getPhpOpcacheMetrics(host, uint16(port), path)
	assert.Nil(t, err)
	assert.NotNil(t, ret)
	assert.NotEmpty(t, ret)
	assert.Contains(t, ret, "used_memory")
	assert.Contains(t, ret, "free_memory")
	assert.Contains(t, ret, "wasted_memory")
	assert.Contains(t, ret, "current_wasted_percentage")
	assert.Contains(t, ret, "num_cached_scripts")
	assert.Contains(t, ret, "num_cached_keys")
	assert.Contains(t, ret, "max_cached_keys")
	assert.Contains(t, ret, "hits")
	assert.Contains(t, ret, "oom_restarts")
	assert.Contains(t, ret, "hash_restarts")
	assert.Contains(t, ret, "manual_restarts")
	assert.Contains(t, ret, "misses")
	assert.Contains(t, ret, "blacklist_misses")
	assert.Contains(t, ret, "blacklist_miss_ratio")
	assert.Contains(t, ret, "opcache_hit_rate")
}
