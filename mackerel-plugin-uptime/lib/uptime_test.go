package mpuptime

import "testing"

func TestFetchMetrics(t *testing.T) {
	var uptime UptimePlugin

	stat, err := uptime.FetchMetrics()
	if err != nil {
		t.Fatal(err)
	}
	seconds := stat["seconds"]
	if seconds <= 0 {
		t.Errorf("invalid seconds value: %f", seconds)
	}
}

func TestGraphDefinition(t *testing.T) {
	var uptime UptimePlugin

	graphdef := uptime.GraphDefinition()
	if len(graphdef) != 1 {
		t.Errorf("GetTempfilename: %d should be 1", len(graphdef))
	}
}
