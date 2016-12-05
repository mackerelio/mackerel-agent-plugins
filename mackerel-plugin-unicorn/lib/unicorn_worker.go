package mpunicorn

import (
	"fmt"
	"time"

	"strconv"
	"strings"
)

func idleWorkerCount(pids []string) (int, error) {
	var beforeCPU []string
	var afterCPU []string
	idles := 0

	for _, pid := range pids {
		cputime, err := cpuTime(pid)
		if err != nil {
			return idles, err
		}
		beforeCPU = append(beforeCPU, cputime)
	}
	time.Sleep(1 * time.Second)
	for _, pid := range pids {
		cputime, err := cpuTime(pid)
		if err != nil {
			return idles, err
		}
		afterCPU = append(afterCPU, cputime)
	}
	for i := range pids {
		b, _ := strconv.Atoi(beforeCPU[i])
		a, _ := strconv.Atoi(afterCPU[i])
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
