package mpprocfd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mattn/go-pipeline"
)

// OpenFd interface
type OpenFd interface {
	getNumOpenFileDesc() (map[string]uint64, error)
}

var openFd OpenFd

// RealOpenFd struct
type RealOpenFd struct {
	process string
}

func (o RealOpenFd) getNumOpenFileDesc() (map[string]uint64, error) {
	fds := make(map[string]uint64)

	// Fetch all pids which contain specified process name
	out, err := pipeline.Output(
		[]string{"ps", "aux"},
		[]string{"grep", o.process},
		[]string{"grep", "-v", "grep"},
		[]string{"grep", "-v", "mackerel-plugin-proc-fd"},
		[]string{"awk", "{print $2}"},
	)
	if err != nil {
		// No matching with p.Process invokes this case
		return nil, err
	}

	// List the number of all open files beloging to each pid
	for _, pid := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		out, err = pipeline.Output(
			[]string{"ls", "-l", fmt.Sprintf("/proc/%s/fd", pid)},
			[]string{"grep", "-v", "total"},
			[]string{"wc", "-l"},
		)
		if err != nil {
			// The process with pid terminates"
			return nil, err
		}

		num, err := strconv.ParseUint(strings.TrimSpace(string(out)), 10, 32)
		if err != nil {
			return nil, err
		}
		fds[pid] = num
	}

	return fds, nil
}
