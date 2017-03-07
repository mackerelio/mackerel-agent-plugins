package mpgcpcomputeengine

import (
	"testing"
)

func TestGraphDefinition(t *testing.T) {
	var ce ComputeEnginePlugin

	graphdef := ce.GraphDefinition()

	expect := 7

	if len(graphdef) != expect {
		t.Errorf("GraphDefinition(): %d should be %d", len(graphdef), expect)
	}
}
