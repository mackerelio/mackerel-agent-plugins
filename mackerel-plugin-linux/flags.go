package main

import (
	"github.com/codegangsta/cli"
)

var flags = []cli.Flag{
	cliTempFile,
	cliType,
}

var cliTempFile = cli.StringFlag{
	Name:   "tempfile, t",
	Value:  "/tmp/mackerel-plugin-linux",
	Usage:  "Set temporary file path.",
	EnvVar: "ENVVAR_TEMPFILE",
}

var cliType = cli.StringFlag{
	Name:   "type, p",
	Value:  "all",
	Usage:  "Select metrics: all, swap, netstat, diskstats, proc_stat, users",
	EnvVar: "ENVVAR_TYPE",
}
