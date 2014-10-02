package main

import (
	"github.com/codegangsta/cli"
)

var Flags = []cli.Flag{
	cliHttpHost,
	cliHttpPort,
	cliStatusPage,
	cliTempFile,
}

var cliHttpHost = cli.StringFlag{
	Name:   "http_host, o",
	Value:  "127.0.0.1",
	Usage:  "Set apache2 listeing ip.",
	EnvVar: "ENVVAR_HTTP_HOST",
}

var cliHttpPort = cli.IntFlag{
	Name:   "http_port, p",
	Value:  80,
	Usage:  "Set apache2 listeing port.",
	EnvVar: "ENVVAR_HTTP_PORT",
}

var cliStatusPage = cli.StringFlag{
	Name:   "status_page, s",
	Value:  "/server-status?auto",
	Usage:  "Set apache2 mod_status page address.",
	EnvVar: "ENVVAR_STATUS_PAGE",
}

var cliTempFile = cli.StringFlag{
	Name:   "tempfile, t",
	Value:  "/tmp/mackerel-plugin-apache2",
	Usage:  "Set temporary file path.",
	EnvVar: "ENVVAR_TEMPFILE",
}
