mackerel-plugin-proxysql
=====================

[ProxySQL](http://www.proxysql.com/) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-proxysql [-host=<host>] [-port=<port>] [-socket=<socket>] [-username=<username>] [-password=<password>] [-tempfile=<tempfile>]
```

## Example of mackerel-agent.conf

```toml
[plugin.metrics.proxysql]
command = "/path/to/mackerel-plugin-proxysql -host=0.0.0.0 -port=6032 -username=proxysql_remoteadmin -password=proxysql_remoteadmin"
```

```toml
[plugin.metrics.proxysql]
command = "/path/to/mackerel-plugin-proxysql -socket=/var/lib/proxysql/proxysql.sock -username=proxysql_remoteadmin -password=proxysql_remoteadmin"
```
