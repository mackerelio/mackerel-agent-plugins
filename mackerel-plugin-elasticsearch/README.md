mackerel-plugin-elasticsearch
=====================

Elasticsearch custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-elasticsearch [-scheme=<'http'|'https'>] [-host=<host>] [-port=<manage_port>] [-tempfile=<tempfile>] [-metric-key-prefix=<prefix>] [-metric-label-prefix=<label-prefix>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.elasticsearch]
command = "/path/to/mackerel-plugin-elasticsearch -port=6666"
```
