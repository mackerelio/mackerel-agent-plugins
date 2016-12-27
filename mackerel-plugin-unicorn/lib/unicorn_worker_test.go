package mpunicorn

import (
	"reflect"
	"testing"
)

type TestFetchUnicornWorkerPidsCommand struct{}

func (r TestFetchUnicornWorkerPidsCommand) Output(command string, args ...string) ([]byte, error) {
	out := ` PID TTY      STAT   TIME COMMAND
  584 ?        Sl     0:27 unicorn worker[7] -c config/unicorn.rb -E deployment
 1857 ?        Sl     0:24 unicorn worker[6] -c config/unicorn.rb -E deployment
 2258 ?        Sl     0:19 unicorn worker[4] -c config/unicorn.rb -E deployment
 2627 ?        Sl     0:21 unicorn worker[2] -c config/unicorn.rb -E deployment
 2872 ?        Sl     0:22 unicorn worker[5] -c config/unicorn.rb -E deployment
 3085 ?        Sl     0:21 unicorn worker[3] -c config/unicorn.rb -E deployment
 3546 ?        Sl     0:17 unicorn worker[10] -c config/unicorn.rb -E deployment
 4392 ?        Sl     0:16 unicorn worker[8] -c config/unicorn.rb -E deployment
 6049 ?        Sl     0:12 unicorn worker[14] -c config/unicorn.rb -E deployment
 8430 ?        Sl     0:13 unicorn worker[0] -c config/unicorn.rb -E deployment
 8744 ?        Sl     0:10 unicorn worker[13] -c config/unicorn.rb -E deployment
 9293 ?        Sl     0:11 unicorn worker[1] -c config/unicorn.rb -E deployment
10425 ?        Sl     0:08 unicorn worker[11] -c config/unicorn.rb -E deployment
11152 ?        Sl     0:05 unicorn worker[9] -c config/unicorn.rb -E deployment
11576 ?        Sl     0:05 unicorn worker[15] -c config/unicorn.rb -E deployment
11685 ?        Sl     0:04 unicorn worker[12] -c config/unicorn.rb -E deployment`
	return []byte(out), nil
}

func TestFetchUnicornWorkerPids(t *testing.T) {
	command = TestFetchUnicornWorkerPidsCommand{}
	masterPid := "30661"
	expectedPids := []string{"584", "1857", "2258", "2627", "2872", "3085", "3546",
		"4392", "6049", "8430", "8744", "9293", "10425", "11152", "11576", "11685"}
	pids, _ := fetchUnicornWorkerPids(masterPid)
	if !reflect.DeepEqual(pids, expectedPids) {
		t.Errorf("fetchUnicornWorkerPids: expected %s but got %s", expectedPids, pids)
	}
}

type TestCPUTimePipedCommands struct{}

func (r TestCPUTimePipedCommands) Output(commands ...[]string) ([]byte, error) {
	return []byte("418\n"), nil
}

func TestIdleWorkerCount(t *testing.T) {
	pipedCommands = TestCPUTimePipedCommands{}
	pids := []string{"584", "1857", "2258", "2627", "2872", "3085", "3546"}
	expectedCount := len(pids)
	c, _ := idleWorkerCount(pids)
	if c != 7 {
		t.Errorf("idleWorkerCount: expected %d but got %d", expectedCount, c)
	}
}

func TestCPUTime(t *testing.T) {
	pipedCommands = TestCPUTimePipedCommands{}
	pid := "3061"
	expectedCPUTime := "418"
	c, _ := cpuTime(pid)
	if c != expectedCPUTime {
		t.Errorf("cpuTime: expected %s but got %s", expectedCPUTime, c)
	}
}
