mackerel-plugin-twemproxy
=====================

twemproxy stats custom metrics plugin for [mackerel-agent](https://github.com/mackerelio/mackerel-agent).

## Synopsis

```shell
mackerel-plugin-twemproxy [-metric-key-prefix=twemproxy] [-timeout=5] [-address=localhost:22222]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.twemproxy]
command = "/path/to/mackerel-plugin-twemproxy"
```

## References

- https://github.com/twitter/twemproxy#observability
