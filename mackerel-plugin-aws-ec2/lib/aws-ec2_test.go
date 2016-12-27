package mpawsec2

import (
	"testing"
)

func TestGraphDefinition(t *testing.T) {
	var ec2 EC2Plugin
	graphdef := ec2.GraphDefinition()

	expected := 6
	if actual := len(graphdef); actual != expected {
		t.Errorf("GraphDefinition(): %d should be %d", actual, expected)
	}
}
