mackerel-plugin-mysql
=====================

MySQL custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-mysql [-host=<host>] [-port=<port>] [-username=<username>] [-password=<password>] [-tempfile=<tempfile>] [-disable_innodb=true] [-enable_extended=true]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.mysql]
command = "/path/to/mackerel-plugin-mysql"
```

