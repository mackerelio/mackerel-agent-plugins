package main

import (
	"github.com/urfave/cli"
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

var cliType = cli.StringSliceFlag{
	Name:   "type, p",
	Value:  &cli.StringSlice{},
	Usage:  "Select metrics type(s) to fetch: all, swap, netstat, diskstats, proc_stat, users",
	EnvVar: "ENVVAR_TYPE",
}
