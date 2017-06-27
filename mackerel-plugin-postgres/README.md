mackerel-plugin-postgres
========================

PostgreSQL custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-postgres -user=<username> -password=<password> [-database=<databasename>] [-sslmode=<sslmode>] [-metric-key-prefix=<prefix>] [-connect_timeout=<timeout>]
```
`-database` is optional.

## Example of mackerel-agent.conf

```
[plugin.metrics.postgres]
command = "/path/to/mackerel-plugin-postgres -user=test -password=secret -database=databasename"
```

## References

- [PostgreSQL Documentation (27.2. The Statistics Collector)](http://www.postgresql.org/docs/9.3/static/monitoring-stats.html)
