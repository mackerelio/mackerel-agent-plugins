package main

import (
	"github.com/codegangsta/cli"
)

var Flags = []cli.Flag{
	cliTempFile,
	cliType,
}

var cliTempFile = cli.StringFlag{
	Name:   "tempfile, t",
	Value:  "/tmp/mackerel-plugin-apache2",
	Usage:  "Set temporary file path.",
	EnvVar: "ENVVAR_TEMPFILE",
}

var cliType = cli.StringFlag{
	Name:   "type, p",
	Value:  "all",
	Usage:  "Select metrics: all, vmstat, netstat, diskstats, proc_stat, users",
	EnvVar: "ENVVAR_TYPE",
}
