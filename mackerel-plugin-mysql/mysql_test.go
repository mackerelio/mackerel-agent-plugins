package main

import (
	"testing"
)

func TestGraphDefinition(t *testing.T) {
	var mysql MySQLPlugin

	graphdef := mysql.GraphDefinition()
	if len(graphdef) != 10 {
		t.Errorf("GetTempfilename: %d should be 10", len(graphdef))
	}
}
