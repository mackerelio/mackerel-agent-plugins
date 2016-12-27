package mpunicorn

import (
	"fmt"

	"strings"
)

func workersMemory() (string, error) {
	out, err := pipedCommands.Output(
		[]string{"ps", "auxw"},
		[]string{"grep", "[u]nicorn"},
		[]string{"grep", "-v", "master"},
		[]string{"awk", "{m+=$6*1024} END{print m;}"},
	)
	if err != nil {
		return "", fmt.Errorf("Cannot get unicorn workers memory: %s", err)
	}
	return strings.Trim(string(out), "\n"), nil
}

func masterMemory() (string, error) {
	out, err := pipedCommands.Output(
		[]string{"ps", "auxw"},
		[]string{"grep", "[u]nicorn"},
		[]string{"grep", "master"},
		[]string{"awk", "{m+=$6*1024} END{print m;}"},
	)
	if err != nil {
		return "", fmt.Errorf("Cannot get unicorn master memory: %s", err)
	}
	return strings.Trim(string(out), "\n"), nil
}

func workersMemoryAvg() (string, error) {
	out, err := pipedCommands.Output(
		[]string{"ps", "auxw"},
		[]string{"grep", "[u]nicorn"},
		[]string{"grep", "-v", "master"},
		[]string{"awk", "{mem=$6*1024+mem; proc++} END{print mem/proc}"},
	)
	if err != nil {
		return "", fmt.Errorf("Cannot get unicorn memory average: %s", err)
	}
	return strings.Trim(string(out), "\n"), nil
}
