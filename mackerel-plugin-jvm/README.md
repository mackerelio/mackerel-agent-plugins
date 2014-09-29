mackerel-plugin-jvm
===================

JVM(jstat) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-jvm -javaname=<javaname> [-jstatpath=</path/to/jstat] [-jpspath=/path/to/jps] [-host=<host>] [-port=<port>]
```

## Requirements

- JVM 1.6 or higher

## Example of mackerel-agent.conf

```
[plugin.metrics.jvm]
command = "/path/to/mackerel-plugin-jvm -javaname=NettyServer -jstatpath=/usr/bin/jstat -jpspath=/usr/bin/jps"
```

## References

- https://github.com/sensu/sensu-community-plugins/blob/master/plugins/java/jstat-metrics.py
- http://docs.oracle.com/javase/7/docs/technotes/tools/share/jstat.html
- https://github.com/kazeburo/jstat2gf
