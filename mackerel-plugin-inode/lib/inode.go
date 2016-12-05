package mpinode

import (
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/mackerel-agent/logging"
)

var logger = logging.GetLogger("metrics.plugin.inode")

// InodePlugin plugin
type InodePlugin struct{}

var dfHeaderPattern = regexp.MustCompile(
	`^Filesystem\s+`,
)

var dfColumnsPattern = regexp.MustCompile(
	`^(.+?)\s+(?:(?:\d+)\s+(?:\d+)\s+(?:\d+)\s+(?:\d+%)\s+|(?:\d+)\s+)?(\d+)\s+(\d+)\s+(\d+%)\s+(.+)$`,
)

var devicePattern = regexp.MustCompile(
	`^/dev/(.*)$`,
)

var deviceUnacceptablePattern = regexp.MustCompile(
	`[^A-Za-z0-9_-]`,
)

//  $ df -iP
// Filesystem      Inodes  IUsed   IFree IUse% Mounted on
// /dev/xvda1     1310720 131197 1179523   11% /
//  $ df -i # on Mac OSX (impossible to display only inode information)
// Filesystem 512-blocks      Used Available Capacity  iused    ifree %iused  Mounted on
// /dev/disk1  974737408 176727800 797497608    19% 22154973 99687201   18%   /

// FetchMetrics interface for mackerelplugin
func (p InodePlugin) FetchMetrics() (map[string]interface{}, error) {
	dfOpt := "-i"
	if runtime.GOOS == "linux" {
		dfOpt = "-iP"
	}
	cmd := exec.Command("df", dfOpt)
	cmd.Env = append(os.Environ(), "LANG=C")
	out, err := cmd.Output()
	if err != nil {
		logger.Warningf("'df -i' command exited with a non-zero status: '%s'", err)
		return nil, err
	}
	result := make(map[string]interface{})
	for _, line := range strings.Split(string(out), "\n") {
		if dfHeaderPattern.MatchString(line) {
			continue
		} else if matches := dfColumnsPattern.FindStringSubmatch(line); matches != nil {
			name := matches[1]
			// https://github.com/docker/docker/blob/v1.5.0/daemon/graphdriver/devmapper/deviceset.go#L981
			if regexp.MustCompile(`^/dev/mapper/docker-`).FindStringSubmatch(name) != nil {
				continue
			}
			if nameMatches := devicePattern.FindStringSubmatch(name); nameMatches != nil {
				device := deviceUnacceptablePattern.ReplaceAllString(nameMatches[1], "_")
				iused, err := strconv.ParseInt(matches[2], 0, 64)
				if err != nil {
					logger.Warningf("Failed to parse value: [%s]", matches[2])
					continue
				}
				ifree, err := strconv.ParseInt(matches[3], 0, 64)
				if err != nil {
					logger.Warningf("Failed to parse value: [%s]", matches[3])
					continue
				}
				result["inode.count."+device+".used"] = uint64(iused)
				result["inode.count."+device+".free"] = uint64(ifree)
				result["inode.count."+device+".total"] = uint64(iused + ifree)
				usedPercentage := 100.0 // 100% if both iused and ifree are 0
				if iused+ifree > 0 {
					usedPercentage = float64(iused) * 100 / float64(iused+ifree)
				}
				result["inode.percentage."+device+".used"] = usedPercentage
			}
		}
	}
	return result, nil
}

// GraphDefinition interface for mackerelplugin
func (p InodePlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"inode.count.#": {
			Label: "Inode",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "used", Label: "used"},
				{Name: "free", Label: "free"},
				{Name: "total", Label: "total"},
			},
		},
		"inode.percentage.#": {
			Label: "Inode Percentage",
			Unit:  "percentage",
			Metrics: []mp.Metrics{
				{Name: "used", Label: "used %"},
			},
		},
	}
}

// Do the plugin
func Do() {
	inode := InodePlugin{}
	helper := mp.NewMackerelPlugin(inode)
	helper.Run()
}
