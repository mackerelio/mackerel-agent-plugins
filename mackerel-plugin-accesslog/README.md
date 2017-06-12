mackerel-plugin-accesslog
=====================

accesslog custom metrics plugin for mackerel.io agent.

Apache log format (common and combined) and LTSV log format are supported.

## Synopsis

```shell
mackerel-plugin-accesslog /path/to/access.log
```

## Example of mackerel-agent.conf

```
[plugin.metrics.accesslog]
command = "/path/to/mackerel-plugin-accesslog /path/to/access.log"
```
