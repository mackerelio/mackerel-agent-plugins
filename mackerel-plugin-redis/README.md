mackerel-plugin-redis
=====================

Redis custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-redis [-hostname=<hostname>] [-port=<port>] [-timeout=<time>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.redis]
command = "/path/to/mackerel-plugin-redis -port=6379 -timeout=5"
```

## References

- http://redis.io/commands/INFO
