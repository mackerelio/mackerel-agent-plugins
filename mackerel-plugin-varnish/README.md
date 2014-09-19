mackerel-plugin-varnish
=====================

Varnish custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-varnish [-host=<host>] [-port=<management_port>] [-tempfile=<tempfile>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.varnish]
command = "/path/to/mackerel-plugin-varnish -port=6666"
```
