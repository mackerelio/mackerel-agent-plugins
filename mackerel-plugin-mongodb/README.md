mackerel-plugin-mongodb
=====================

MongoDB custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-mongodb [-host=<host>] [-port=<port>] [-username=<username>] [-password=<password>] [-tempfile=<tempfile>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.mongodb]
command = "/path/to/mackerel-plugin-mongodb"
```
