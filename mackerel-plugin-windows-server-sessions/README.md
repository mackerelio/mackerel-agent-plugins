mackerel-plugin-windows-server-sessions
=======================================

Windows Server Sessions custom metrics plugin for mackerel-agent.

## Usage

```shell
mackerel-plugin-windows-server-sessions [-legacymetricname]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.td-table-count]
command = "/path/to/mackerel-plugin-windows-server-sessions"
```
