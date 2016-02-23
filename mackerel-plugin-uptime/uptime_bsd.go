// +build freebsd netbsd darwin

package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func fetchMetrics() (map[string]interface{}, error) {
	cmd := exec.Command("sysctl", "-n", "kern.boottime")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("faild to fetch uptime: %s", err)
	}
	// { sec = 1455448176, usec = 0 } Sun Feb 14 20:09:36 2016
	cols := strings.Split(out.String(), " ")
	epoch, err := strconv.ParseUint(strings.Trim(cols[3], ","), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Faild to parse uptime: %s", err)
	}
	return map[string]interface{}{"seconds": float64(uint64(time.Now().Unix()) - epoch)}, nil
}
