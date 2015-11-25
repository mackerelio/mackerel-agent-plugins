package main

import (
	"fmt"
	"time"

	"strconv"
	"strings"
)

func usedMemory() (string, error) {
	out, err := pipedCommands.Output(
		[]string{"ps", "auxw"},
		[]string{"grep", "[u]nicorn"},
		[]string{"awk", "{m+=$6*1024} END{print m;}"},
	)
	if err != nil {
		return "", fmt.Errorf("Cannot get unicorn used memory: %s", err)
	}
	return strings.Trim(string(out), "\n"), nil
}

func averageMemory() (string, error) {
	out, err := pipedCommands.Output(
		[]string{"ps", "auxw"},
		[]string{"grep", "[u]nicorn"},
		[]string{"grep", "-v", "master"},
		[]string{"awk", "{mem=$6*1024+mem; proc++} END{printf(\"%d\n\", mem/proc)}"},
	)
	if err != nil {
		return "", fmt.Errorf("Cannot get unicorn memory average: %s", err)
	}
	return strings.Trim(string(out), "\n"), nil
}

func idleWorkerCount(pids []string) (int, error) {
	var beforeCpu []string
	var afterCpu []string
	idles := 0

	for _, pid := range pids {
		cputime, err := cpuTime(pid)
		if err != nil {
			return idles, err
		}
		beforeCpu = append(beforeCpu, cputime)
	}
	time.Sleep(1 * time.Second)
	for _, pid := range pids {
		cputime, err := cpuTime(pid)
		if err != nil {
			return idles, err
		}
		afterCpu = append(afterCpu, cputime)
	}
	for i, _ := range pids {
		b, _ := strconv.Atoi(beforeCpu[i])
		a, _ := strconv.Atoi(afterCpu[i])
		if (a - b) == 0 {
			idles++
		}
	}

	return idles, nil
}

func fetchUnicornWorkerPids(m string) ([]string, error) {
	var workerPids []string

	out, err := command.Output("ps", "wh", "--ppid", m)
	if err != nil {
		return workerPids, fmt.Errorf("Failed to ps of command: %s", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, "worker") {
			continue
		}
		words := strings.SplitN(strings.TrimSpace(line), " ", 5)
		if len(words) < 5 {
			continue
		}
		workerPids = append(workerPids, words[0])
	}

	if len(workerPids) > 0 {
		return workerPids, nil
	}

	return workerPids, fmt.Errorf("Cannot get unicorn worker pids")
}

func cpuTime(pid string) (string, error) {
	out, err := pipedCommands.Output(
		[]string{"cat", fmt.Sprintf("/proc/%s/stat", pid)},
		[]string{"awk", "{print $14+$15}"},
	)
	if err != nil {
		return "", fmt.Errorf("Failed to cat /proc/%s/stat: %s", pid, err)
	}

	return strings.Trim(string(out), "\n"), nil
}
