mackerel-plugin-multicore
=====================

Get multicore CPU metrics for linux.

- CPU usage by cores
- loadavg5 per cores

## Synopsis

```shell
mackerel-plugin-multicore
```

## Example of mackerel-agent.conf

```
[plugin.metrics.multicore]
command = "/path/to/mackerel-plugin-multicore"
```
