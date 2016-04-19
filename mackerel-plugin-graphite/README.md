mackerel-plugin-graphite
===================

Graphite carbon custom metrics plugin for mackerel.io agent.  

## Synopsis

This plugin posts metrics that carbon-cache and carbon-relay collect themselves.  
carbon-aggregator is not implemented.  
It posts all metrics which were collected by 15 minutes because of reflection delay.  

```shell
mackerel-plugin-graphite -host=<host name> -webhost=<graphite-web host name> -webport=<graphite-web host port> -type=(cache or relay) (-instance=<instance name> -metric-label-prefix=<metric label prefix>)
```

## Example of mackerel-agent.conf

### carbon-cache

```
[plugin.metrics.graphite-carbon]
command = "/path/to/mackerel-plugin-graphite -host=127.0.0.1 -webhost=hostname -port=8000 -type=cache"
```

You don't need specify instance option.  
This plugin automatically attaches `*` to instance name which means all instances of existing.


### carbon-relay

```
[plugin.metrics.graphite-carbon]
command = "/path/to/mackerel-plugin-graphite -host=127.0.0.1 -webhost=hostname -port=8000 -type=relay -instance=a"
```

If you use carbon-relay, you must specify instance option.  

