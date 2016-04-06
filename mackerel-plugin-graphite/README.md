mackerel-plugin-graphite
===================

Graphite carbon custom metrics plugin for mackerel.io agent.  

## Synopsis

```shell
mackerel-plugin-graphite -host=<host name> -webhost=<graphite-web host name> -webport=<graphite-web host port> -type=(cache or relay) (-instance=<instance name> -metric-label-prefix=<metric label prefix>)
```

## Example of mackerel-agent.conf

```
[plugin.metrics.graphite-carbon]
command = "/path/to/mackerel-plugin-graphite -host=hostname -webhost=hostname -port=8000 -type=cache"
```
