mackerel-plugin-postgres
========================

PostgreSQL custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-postgres -user=<username> -password=<password>
```

## Example of mackerel-agent.conf

```
[plugin.metrics.postgres]
command = "/path/to/mackerel-plugin-postgres -user=test -password=secret"
```

## References

- [PostgreSQL Documentation (27.2. The Statistics Collector)](http://www.postgresql.org/docs/9.3/static/monitoring-stats.html)
