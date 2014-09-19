mackerel-plugin-squid
=====================

Squid custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-squid [-host=<host>] [-port=<manager_port>] [-tempfile=<tempfile>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.squid]
command = "/path/to/mackerel-plugin-squid -port=6666"
```
