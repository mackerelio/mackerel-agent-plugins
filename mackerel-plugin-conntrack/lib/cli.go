package mpconntrack

import (
	"flag"
	"fmt"
	"io"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK             int = 0
	ExitCodeParseFlagError int = 1 + iota
)

// ConntrackPlugin mackerel plugin for *_conntrack.
type ConntrackPlugin struct{}

// GraphDefinition interface for mackerelplugin.
func (c ConntrackPlugin) GraphDefinition() map[string]mp.Graphs {
	// graphdef is Graph definition for mackerelplugin.
	var graphdef = map[string]mp.Graphs{
		"conntrack.count": {
			Label: "Conntrack Count",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "*", Label: "%1", Diff: false, Stacked: true, Type: "uint64"},
			},
		},
	}

	return graphdef
}

// FetchMetrics interface for mackerelplugin.
func (c ConntrackPlugin) FetchMetrics() (map[string]interface{}, error) {
	conntrackCount, err := CurrentValue(ConntrackCountPaths)
	if err != nil {
		return nil, err
	}

	conntrackMax, err := CurrentValue(ConntrackMaxPaths)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]interface{})
	stat["conntrack.count.used"] = conntrackCount
	stat["conntrack.count.free"] = (conntrackMax - conntrackCount)

	return stat, nil
}

// CLI is the object for command line interface.
type CLI struct {
	outStream, errStream io.Writer
}

// Run is to parse flags and Run helper (MackerelPlugin) with the given arguments.
func (c *CLI) Run(args []string) int {
	// Flags
	var (
		tempfile string
		version  bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.BoolVar(&version, "version", false, "Print version information and quit.")
	flags.StringVar(&tempfile, "tempfile", "", "Temp file name")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError
	}

	// Show version
	if version {
		fmt.Fprintf(c.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	// Create MackerelPlugin for Conntrack
	var cp ConntrackPlugin
	helper := mp.NewMackerelPlugin(cp)
	helper.Tempfile = tempfile

	helper.Run()

	return ExitCodeOK
}
