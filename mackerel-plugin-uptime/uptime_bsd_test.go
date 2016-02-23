// +build freebsd netbsd darwin

package main

import (
	"reflect"
	"testing"
)

func TestParseMetrics(t *testing.T) {
	r, err := calcMetrics("{ sec = 1455448176, usec = 0 } Sun Feb 14 20:09:36 2016\n", 1456242880)
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}
	if !reflect.DeepEqual(r, map[string]interface{}{"seconds": float64(794704)}) {
		t.Errorf("something went wrong. failed to parse uptime")
	}
}
