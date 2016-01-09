mackerel-plugin-murmur
======================

Murmur custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-murmur [-host=<host>] [-port=<port>] [-tempfile=<tempfile>] [-timeout=<timeout_ms>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.murmur]
command = "/path/to/mackerel-plugin-murmur -host=localhost -port=64738 -timeout 250"
```