mackerel-plugin-windows-process-stats
=======================================

Windows Processes custom metrics plugin for mackerel-agent.

## Usage

```shell
mackerel-plugin-windows-process-stats -process=<process name> [-metric-key-prefix=<metric key prefix>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.windows-process-stats]
command = "/path/to/mackerel-plugin-windows-process-stats -process=<process name> -metric-key-prefix=<metric key prefix>"
```
