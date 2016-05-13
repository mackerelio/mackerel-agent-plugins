mackerel-plugin-uptime
=====================

uptime custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-uptime [-metric-key-prefix=uptime]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.uptime]
command = "/path/to/mackerel-plugin-uptime"
```
