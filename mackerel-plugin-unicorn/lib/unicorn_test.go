package mpunicorn

import "testing"

func TestFetchMetrics(t *testing.T) {
	var unicorn UnicornPlugin

	stat, _ := unicorn.FetchMetrics()
	if len(stat) != 5 {
		t.Errorf("GetStat: %d should be 5", len(stat))
	}
}

func TestGraphDefinition(t *testing.T) {
	var unicorn UnicornPlugin

	graphdef := unicorn.GraphDefinition()
	if len(graphdef) != 2 {
		t.Errorf("GetTempfilename: %d should be 2", len(graphdef))
	}
}
