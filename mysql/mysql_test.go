package main

import (
	"testing"
)

func TestGetTempfilename(t *testing.T) {
	var mysql MySQLPlugin
	tempfile := "localhost"
	mysql.Tempfile = tempfile
	if tempfile != mysql.GetTempfilename() {
		t.Errorf("GetTempfilename: %s is not target %s", mysql.GetTempfilename(), tempfile)
	}
}
