mackerel-plugin-fluentd
=========================

Fluentd (http://www.fluentd.org/) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-fluentd [-host=<host>] [-port=<port>] [-tempfile=<tempfile>] [-plugin-type=<plugin-type>] [-plugin-id-pattern=<plugin-id-pattern>] [-workers=<workers>] [-extended_metrics=<metric-names>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.fluentd]
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

See https://docs.fluentd.org/input/monitor_agent in details.

If you have specified the `workers` parameter in fluentd's `<system>` directive, you can use the `-workers` option.

```
<system>
  workers 3
</system>

<source>
  @type monitor_agent
  port 24230 # worker0: 24230, worker1: 24231, worker2: 24232
</source>
```

See https://docs.fluentd.org/input/monitor_agent#multi-process-environment in details.

## License

Released under the MIT license
http://opensource.org/licenses/mit-license.php

Original version of the plugin https://github.com/y-matsuwitter/mackerel-fluentd
Copyright (c) 2015 Yuki Matsumoto

Current version is forked from the original version under the MIT license.
Copyright (c) 2015 Shinji Tanaka

