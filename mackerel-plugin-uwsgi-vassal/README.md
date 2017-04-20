mackerel-plugin-uwsgi-vassal
=====================

uWSGI (vassal) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-uwsgi-vassal [-socket=<http://uri|unix:///tmp/tmp.sock>]
```

## Requirements

This plugin requires that uWSGI use `--stats` options in vassal section.

## Example of mackerel-agent.conf

```
[plugin.metrics.uwsgi-vassal]
command = "/path/to/mackerel-plugin-uwsgi-vassal --socket unix:///var/run/uwsgi-vassal-stats.sock"
```
