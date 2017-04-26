mackerel-plugin-uwsgi-vassal
=====================

uWSGI (vassal) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-uwsgi-vassal [-socket=<'http://host:port'|'unix:///var/run/uwsgi/stats.sock'>]
```

## Requirements

This plugin requires that uWSGI use `--stats` options in vassal section.
- [The uWSGI Stats Server — uWSGI 2.0 documentation](http://uwsgi-docs.readthedocs.io/en/latest/StatsServer.html)

## Example of mackerel-agent.conf

```
[plugin.metrics.uwsgi-vassal]
command = "/path/to/mackerel-plugin-uwsgi-vassal --socket unix:///var/run/uwsgi-vassal-stats.sock"
```
