package main

import "testing"

func TestGraphDefinition(t *testing.T) {
	var unicorn UnicornPlugin

	graphdef := unicorn.GraphDefinition()
	if len(graphdef) != 2 {
		t.Errorf("GetTempfilename: %d should be 2", len(graphdef))
	}
}
