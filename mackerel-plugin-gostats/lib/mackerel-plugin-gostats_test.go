package mpgostats

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseStats(t *testing.T) {
	stat := `{
  "time": 1449124022112358000,
  "go_version": "go1.5.1",
  "go_os": "darwin",
  "go_arch": "amd64",
  "cpu_num": 4,
  "goroutine_num": 6,
  "gomaxprocs": 4,
  "cgo_call_num": 5,
  "memory_alloc": 213360,
  "memory_total_alloc": 213360,
  "memory_sys": 3377400,
  "memory_lookups": 15,
  "memory_mallocs": 1137,
  "memory_frees": 0,
  "memory_stack": 393216,
  "heap_alloc": 213360,
  "heap_sys": 655360,
  "heap_idle": 65536,
  "heap_inuse": 589824,
  "heap_released": 0,
  "heap_objects": 1137,
  "gc_next": 4194304,
  "gc_last": 0,
  "gc_num": 0,
  "gc_per_second": 0,
  "gc_pause_per_second": 0,
  "gc_pause": []
}`

	expected := map[string]float64{
		"goroutine_num":       6.0,
		"cgo_call_num":        5.0,
		"memory_sys":          3377400.0,
		"memory_alloc":        213360.0,
		"memory_stack":        393216.0,
		"memory_lookups":      15.0,
		"memory_frees":        0.0,
		"memory_mallocs":      1137.0,
		"heap_sys":            655360.0,
		"heap_idle":           65536.0,
		"heap_inuse":          589824.0,
		"heap_released":       0.0,
		"gc_num":              0.0,
		"gc_per_second":       0.0,
		"gc_pause_per_second": 0.0,
	}

	m := GostatsPlugin{}
	got, err := m.parseStats(strings.NewReader(stat))
	if err != nil {
		t.Fatalf("error should be nil but got: %v", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("stats differs: %v (expected: %v)", got, expected)
	}
}
