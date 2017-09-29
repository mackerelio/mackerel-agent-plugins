mackerel-plugin-mcrouter
=====================

Mcrouter stats custom metrics plugin for [mackerel-agent](https://github.com/mackerelio/mackerel-agent).

## Synopsis

```shell
mackerel-plugin-mcrouter -stats-file /path/to/mcrouter.stats [-metric-key-prefix=mcrouter]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.mcrouter]
command = "/path/to/mackerel-plugin-mcrouter -stats-file /path/to/mcrouter.stats"
```

## Graphs and Metrics

### mcrouter.cmd_count

- mcrouter.cmd_count.cmd_add_count
- mcrouter.cmd_count.cmd_cas_count
- mcrouter.cmd_count.cmd_decr_count
- mcrouter.cmd_count.cmd_delete_count
- mcrouter.cmd_count.cmd_get_count
- mcrouter.cmd_count.cmd_gets_count
- mcrouter.cmd_count.cmd_incr_count
- mcrouter.cmd_count.cmd_lease_get_count
- mcrouter.cmd_count.cmd_lease_set_count
- mcrouter.cmd_count.cmd_meta_count
- mcrouter.cmd_count.cmd_other_count
- mcrouter.cmd_count.cmd_replace_count
- mcrouter.cmd_count.cmd_set_count
- mcrouter.cmd_count.cmd_stats_count

### mcrouter.result_count

- mcrouter.result_count.result_busy_all_count
- mcrouter.result_count.result_busy_count
- mcrouter.result_count.result_connect_error_all_count
- mcrouter.result_count.result_connect_error_count
- mcrouter.result_count.result_connect_timeout_all_count
- mcrouter.result_count.result_connect_timeout_count
- mcrouter.result_count.result_data_timeout_all_count
- mcrouter.result_count.result_data_timeout_count
- mcrouter.result_count.result_error_all_count
- mcrouter.result_count.result_error_count
- mcrouter.result_count.result_local_error_all_count
- mcrouter.result_count.result_local_error_count
- mcrouter.result_count.result_tko_all_count
- mcrouter.result_count.result_tko_count

### mcrouter.request_processing_time

- mcrouter.request_processing_time.duration_us

## References

- https://github.com/facebook/mcrouter/wiki/Stats-list

