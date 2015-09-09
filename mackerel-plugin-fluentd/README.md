mackerel-plugin-fluentd
=========================

Fluentd (http://www.fluentd.org/) custom metrics plugin for mackerel.io agent.

The source code of this plugin is based on https://github.com/y-matsuwitter/mackerel-fluentd .

## Synopsis

```shell
mackerel-plugin-fluentd [-host=<host>] [-port=<port>] [-tempfile=<tempfile>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.memcached]
command = "/path/to/mackerel-plugin-fluentd"
```

## Enable monitor_agent for fluentd

This plugin needs to enable monitor_agent at the target fluentd process.
Add following configuraion to your fluentd.conf.

```
<source>
type monitor_agent
bind 0.0.0.0
port 24220
</source>
```

See http://docs.fluentd.org/articles/monitoring in details.
