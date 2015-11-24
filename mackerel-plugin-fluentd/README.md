mackerel-plugin-fluentd
=========================

Fluentd (http://www.fluentd.org/) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-fluentd [-host=<host>] [-port=<port>] [-tempfile=<tempfile>]
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

See http://docs.fluentd.org/articles/monitoring in details.

## License

Released under the MIT license
http://opensource.org/licenses/mit-license.php

Original version of the plugin https://github.com/y-matsuwitter/mackerel-fluentd
Copyright (c) 2015 Yuki Matsumoto

Current version is forked from the original version under the MIT license.
Copyright (c) 2015 Shinji Tanaka

