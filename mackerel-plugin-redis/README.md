mackerel-plugin-redis
=====================

Redis custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-redis [-host=<hostname>] [-port=<port>] [-password=<password>] [-socket=<unix socket>] [-timeout=<time>] [-metric-key-prefix=<prefix>] [-config-command=<CONFIG command name>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.redis]
command = "/path/to/mackerel-plugin-redis -port=6379 -timeout=5"
```

### Using multiple Redis instances on one server

```
[plugin.metrics.redis6379]
command = "/path/to/mackerel-plugin-redis -port=6379 -timeout=5 -metric-key-prefix=redis6379"

[plugin.metrics.redis6380]
command = "/path/to/mackerel-plugin-redis -port=6380 -timeout=5 -metric-key-prefix=redis6380"
```

## References

- http://redis.io/commands/INFO
