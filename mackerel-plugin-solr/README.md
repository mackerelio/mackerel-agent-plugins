mackerel-plugin-solr
=====================

Solr custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-solr [-host=<hostname>] [-port=<port>]
```

## Example of mackerel-agent.conf

### Default

```
[plugin.metrics.solr]
command = "/path/to/mackerel-plugin-solr"
```

### Custom

```
[plugin.metrics.solr]
command = "/path/to/mackerel-plugin-solr -host=192.168.33.10 -port=8984"
```

## Dependency Solr URL

- http://{host}:{port}/solr/admin/cores
- http://{host}:{port}/solr/{core name}/admin/mbeans

## Munin Solr plugin

- https://github.com/munin-monitoring/contrib/tree/master/plugins/solr

## Apache Solr official website

- http://lucene.apache.org/solr/
