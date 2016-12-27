package mpconntrack

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestGraphDefinition(t *testing.T) {
	var conntrack ConntrackPlugin

	graphdef := conntrack.GraphDefinition()
	if len(graphdef) != 1 {
		t.Errorf("GetTempfilename: %d should be 1", len(graphdef))
	}
}

func TestCLI_Run(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("mackerel-plugin-conntrack -version", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("ExitStatus=%d, want %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("mackerel-plugin-conntrack version %s", Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("Output=%q, want %q", errStream.String(), expected)
	}

	args = strings.Split("mackerel-plugin-conntrack -unknown", " ")

	status = cli.Run(args)
	if status != ExitCodeParseFlagError {
		t.Errorf("ExitStatus=%d, want %d", status, ExitCodeOK)
	}

}
