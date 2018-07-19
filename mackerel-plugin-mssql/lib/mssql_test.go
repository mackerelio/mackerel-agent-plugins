// +build windows

package mpmssql

import (
	"os/exec"
	"strings"
	"testing"
)

func instance() (string, error) {
	b, err := exec.Command("wmi2struct", "-l").CombinedOutput()
	if err != nil {
		return "", nil
	}
	for _, line := range strings.Split(string(b), "\n") {
		if strings.Contains(line, "SQLEXPRESS") {
			return "SQLEXPRESS", nil
		}
		if strings.Contains(line, "SQLServer") {
			return "SQLSERVER", nil
		}
	}
	return "", nil
}

func TestGraphDefinition(t *testing.T) {
	var mssql MSSQLPlugin

	name, err := instance()
	if err != nil {
		t.Fatal(err)
	}
	if name == "" {
		t.Skip("SQLServer/SQLEXPRESS is not installed")
	}
	graphdef := mssql.GraphDefinition()
	if len(graphdef) != 3 {
		t.Errorf("GraphDefinition: %d definitions should be exists", len(graphdef))
	}
	values, err := mssql.FetchMetrics()
	if err != nil {
		t.Errorf("FetchMetrics: %v", err)
	}
	if len(values) != 10 {
		t.Errorf("FetchMetrics: %d metrics must be exists", len(values))
	}
}
