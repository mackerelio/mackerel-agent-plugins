package mpphpapc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPhpApcStatus_1(t *testing.T) {
	stub := `memory_segments:1
segment_size:134217592
total_memory:134217592
cached_files_count:392
cached_files_size:54266016
cache_hits:606130
cache_misses:392
cache_full_count:0
user_cache_vars_count:770
user_cache_vars_size:45835056
user_cache_hits:8334
user_cache_misses:10997
user_cache_full_count:0`

	stat := make(map[string]float64)

	err := parsePhpApcStatus(stub, &stat)
	assert.Nil(t, err)
	assert.EqualValues(t, stat["memory_segments"], 1)
	assert.EqualValues(t, stat["segment_size"], 134217592)
	assert.EqualValues(t, stat["total_memory"], 134217592)
	assert.EqualValues(t, stat["cached_files_count"], 392)
	assert.EqualValues(t, stat["cached_files_size"], 54266016)
	assert.EqualValues(t, stat["cache_hits"], 606130)
	assert.EqualValues(t, stat["cache_misses"], 392)
	assert.EqualValues(t, stat["cache_full_count"], 0)
	assert.EqualValues(t, stat["user_cache_vars_count"], 770)
	assert.EqualValues(t, stat["user_cache_vars_size"], 45835056)
	assert.EqualValues(t, stat["user_cache_hits"], 8334)
	assert.EqualValues(t, stat["user_cache_misses"], 10997)
	assert.EqualValues(t, stat["user_cache_full_count"], 0)
}

func TestGetPhpApcMetrics_1(t *testing.T) {
	stub := `memory_segments:1
segment_size:134217592
total_memory:134217592
cached_files_count:392
cached_files_size:54266016
cache_hits:606130
cache_misses:392
cache_full_count:0
user_cache_vars_count:770
user_cache_vars_size:45835056
user_cache_hits:8334
user_cache_misses:10997
user_cache_full_count:0`

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

	ret, err := getPhpApcMetrics(host, uint16(port), path)
	assert.Nil(t, err)
	assert.NotNil(t, ret)
	assert.NotEmpty(t, ret)
	assert.Contains(t, ret, "memory_segments")
	assert.Contains(t, ret, "segment_size")
	assert.Contains(t, ret, "total_memory")
	assert.Contains(t, ret, "cached_files_count")
	assert.Contains(t, ret, "cached_files_size")
	assert.Contains(t, ret, "cache_hits")
	assert.Contains(t, ret, "cache_misses")
	assert.Contains(t, ret, "cache_full_count")
	assert.Contains(t, ret, "user_cache_vars_count")
	assert.Contains(t, ret, "user_cache_vars_size")
	assert.Contains(t, ret, "user_cache_hits")
	assert.Contains(t, ret, "user_cache_misses")
	assert.Contains(t, ret, "user_cache_full_count")
}
