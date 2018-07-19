mackerel-plugin-mssql
=====================

MSSQL custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-mssql [-metric-key-prefix=<PREFIX>] [-instance=<SQLSERVER|SQLEXPRESS>] [-tempfile=<tempfile>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.mssql]
command = "/path/to/mackerel-plugin-mssql"
```

## Supported Version

This plugin support following version (or later) of SQL servers.

* Microsoft SQL Server 2017
* Microsoft SQL Express 2017

## Development

To update lib/wmi.go, you need to install [mattn/wmi2struct](https://github.com/mattn/wmi2struct) on Windows.
