package main

import (
	"reflect"
	"testing"
)

func TestCalcMetrics(t *testing.T) {
	r, err := calcMetrics("481453.56 1437723.27\n")
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}
	if !reflect.DeepEqual(r, map[string]interface{}{"seconds": float64(481453.56)}) {
		t.Errorf("something went wrong. failed to parse uptime")
	}
}
