mackerel-plugin-haproxy
=====================

HAProxy custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-haproxy [-host=<host>] [-port=<port>] [-path=<stats-path>] [-scheme=<http|https>] [-tempfile=<tempfile>]
or
mackerel-plugin-haproxy [-uri=<uri>] [-tempfile=<tempfile>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.haproxy]
command = "/path/to/mackerel-plugin-haproxy -port=8000"
```
