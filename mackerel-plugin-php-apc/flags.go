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
	Usage:  "Set httpd listeing ip.",
	EnvVar: "ENVVAR_HTTP_HOST",
}

var cliHttpPort = cli.IntFlag{
	Name:   "http_port, p",
	Value:  80,
	Usage:  "Set httpd listeing port.",
	EnvVar: "ENVVAR_HTTP_PORT",
}

var cliStatusPage = cli.StringFlag{
	Name:   "status_page, s",
	Value:  "/mackerel/php-apc.php",
	Usage:  "Set httpd mod_status page address.",
	EnvVar: "ENVVAR_STATUS_PAGE",
}

var cliTempFile = cli.StringFlag{
	Name:   "tempfile, t",
	Value:  "/tmp/mackerel-plugin-php-apc",
	Usage:  "Set temporary file path.",
	EnvVar: "ENVVAR_TEMPFILE",
}
