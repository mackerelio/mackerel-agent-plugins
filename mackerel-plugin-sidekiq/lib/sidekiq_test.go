package mpsidekiq

import (
	"testing"
)

func TestGraphDefinition(t *testing.T) {
	var sp SidekiqPlugin

	graphdef := sp.GraphDefinition()

	expect := 3

	if len(graphdef) != expect {
		t.Errorf("GraphDefinition(): %d should be %d", len(graphdef), expect)
	}
}
