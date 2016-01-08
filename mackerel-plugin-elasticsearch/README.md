mackerel-plugin-elasticsearch
=====================

Elasticsearch custom metrics plugin for mackerel.io agent.

NOTICE: This plugin only support Elasticsearch versions 1.6.0 and higher.

## Synopsis

```shell
mackerel-plugin-elasticsearch [-host=<host>] [-port=<manage_port>] [-tempfile=<tempfile>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.elasticsearch]
command = "/path/to/mackerel-plugin-elasticsearch -port=6666"
```
