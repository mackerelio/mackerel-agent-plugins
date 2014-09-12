package main

import (
	"github.com/codegangsta/cli"
)

var Flags = []cli.Flag{
	cliTempFile,
}

var cliTempFile = cli.StringFlag{
	Name:   "tempfile, t",
	Value:  "/tmp/mackerel-plugin-apache2",
	Usage:  "Set temporary file path.",
	EnvVar: "ENVVAR_TEMPFILE",
}
