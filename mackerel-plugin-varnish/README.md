mackerel-plugin-varnish
=====================

Varnish custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-varnish [-varnishstat=<varnishstat-path>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.varnish]
command = "/path/to/mackerel-plugin-varnish"
```
