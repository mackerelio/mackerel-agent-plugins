package main

import (
	"github.com/codegangsta/cli"
)

var Flags = []cli.Flag{
	cliTemperaturePath,
	cliTempFile,
}

var cliTemperaturePath = cli.StringFlag{
	Name:   "temperature_path, s",
	Value:  "/sys/class/thermal/thermal_zone0/temp",
	Usage:  "Set Raspberry PI CPU Temperature path.",
	EnvVar: "ENVVAR_TEMPERATURE_PATH",
}

var cliTempFile = cli.StringFlag{
	Name:   "tempfile, t",
	Value:  "/tmp/mackerel-plugin-rpi",
	Usage:  "Set temporary file path.",
	EnvVar: "ENVVAR_TEMPFILE",
}
