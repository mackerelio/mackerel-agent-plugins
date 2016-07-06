package main

import (
	"github.com/urfave/cli"
)

var flags = []cli.Flag{
	cliHTTPHost,
	cliHTTPPort,
	cliHeader,
	cliStatusPage,
	cliTempFile,
}

var cliHTTPHost = cli.StringFlag{
	Name:   "http_host, o",
	Value:  "127.0.0.1",
	Usage:  "Set apache2 listeing ip.",
	EnvVar: "ENVVAR_HTTP_HOST",
}

var cliHTTPPort = cli.IntFlag{
	Name:   "http_port, p",
	Value:  80,
	Usage:  "Set apache2 listeing port.",
	EnvVar: "ENVVAR_HTTP_PORT",
}

var cliHeader = cli.StringSliceFlag{
	Name:  "header, H",
	Value: &cli.StringSlice{},
	Usage: "Set http header. (e.g. \"Host: servername\")",
	EnvVar: "ENVVAR_HEADER",
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
