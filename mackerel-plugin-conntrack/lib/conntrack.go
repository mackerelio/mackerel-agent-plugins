package mpconntrack

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// ConntrackCountPaths is paths to conntrack_count files.
var ConntrackCountPaths = []string{
	"/proc/sys/net/netfilter/nf_conntrack_count",
	"/proc/sys/net/ipv4/netfilter/ip_conntrack_count",
}

// ConntrackMaxPaths is paths to conntrack_max files.
var ConntrackMaxPaths = []string{
	"/proc/sys/net/nf_conntrack_max",
	"/proc/sys/net/netfilter/nf_conntrack_max",
	"/proc/sys/net/ipv4/ip_conntrack_max",
	"/proc/sys/net/ipv4/netfilter/ip_conntrack_max",
}

// Exists returns whether file exists or not.
func Exists(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

// FindFile returns a first matching file path to search multiple paths.
func FindFile(paths []string) (f string, err error) {
	for _, f := range paths {
		if Exists(f) {
			return f, nil
		}
	}

	return "", fmt.Errorf("Cannot find files %s", paths)
}

// CurrentValue returns a value from a file.
func CurrentValue(paths []string) (n uint64, err error) {
	// Get a path.
	path, err := FindFile(paths)
	if err != nil {
		return 0, err
	}

	// Check whether a file is open.
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Read a value from file.
	cnt := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cnt = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	// Convert to uint64.
	n, err = strconv.ParseUint(cnt, 10, 64)
	return n, nil
}
