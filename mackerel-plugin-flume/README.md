mackerel-plugin-flume
=====================

Flume custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-flume [-host=<host>] [-port=<port>] [-metric-key-prefix=<prefix>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.flume]
command = "/path/to/mackerel-plugin-flume"
```

## Documents

* [Monitoring](https://flume.apache.org/FlumeUserGuide.html#monitoring)
* [JSON Reporting](https://flume.apache.org/FlumeUserGuide.html#json-reporting)

