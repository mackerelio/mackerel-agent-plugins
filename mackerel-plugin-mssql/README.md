mackerel-plugin-mssql
=====================

MSSQL custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-mssql [-prefix=<PREFIX>] [-instance=<SQLSERVER|SQLEXPRESS>] [-tempfile=<tempfile>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.mssql]
command = "/path/to/mackerel-plugin-mssql"
```

