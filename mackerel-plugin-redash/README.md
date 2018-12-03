mackerel-plugin-redash
=====================

redash stats custom metrics plugin for [mackerel-agent](https://github.com/mackerelio/mackerel-agent).

## Synopsis

```shell
mackerel-plugin-redash -api-key=YOUR_REDASH_API_KEY [-metric-key-prefix=redash] [-timeout=5] [-uri=http://localhost/api/admin/queries/tasks]
```

`REDASH_API_KEY` environment variable is used as `api_key` query in uri if it is set.

## Example of mackerel-agent.conf

```
[plugin.metrics.redash]
command = "/path/to/mackerel-plugin-redash"
```

## Notes

- This plugin uses an unsafe(unofficial) API.
 - Operation verified with version `0.11.0+b2016`.
- This plugin does not collect all metrics fetched by the uri.

## References

- https://discuss.redash.io/t/monitoring-the-status-of-tasks-using-a-http-request/683
