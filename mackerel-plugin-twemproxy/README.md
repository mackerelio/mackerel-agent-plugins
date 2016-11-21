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

## Notes

This plugin does not collect metrics of `fragments` and `server_ejected_at`.
See https://github.com/mackerelio/mackerel-agent-plugins/pull/283 for details.

## References

- https://github.com/twitter/twemproxy#observability
