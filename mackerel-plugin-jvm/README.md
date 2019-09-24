mackerel-plugin-jvm
===================

JVM(jstat) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-jvm -javaname=<javaname> [-pidfile=</path/to/pidfile>] [-jstatpath=</path/to/jstat] [-jpspath=/path/to/jps] [-jinfopath=/path/to/jinfo] [-remote=<host:port>]
```

## Requirements

- JVM 1.6 or higher

## Example of mackerel-agent.conf

```
[plugin.metrics.jvm]
command = "/path/to/mackerel-plugin-jvm -javaname=NettyServer -jstatpath=/usr/bin/jstat -jpspath=/usr/bin/jps -jinfopath=/usr/bin/jinfo"
user = "SOME_USER_NAME"
```

## Monitoring remote JVM

This plugin can retrieve metrics from remote jstatd with rmi protocol by setting `-remote` option.
In this case, following limitations are applied:
- jps and jstat commands must be executable localy from this plugin
- 'CMS Initiating Occupancy Fraction' metric cannot be retrieved remotely

## About javaname

You can check javaname by jps command.

```shell
# jps
14203 NettyServer
14822 Jps
```

Please choose an arbitrary name as `javaname` when you use `pidfile` option.
It is just used as a prefix of graph label.

## User to execute this plugin

This plugin (as well as the jps command explained above) must be executed by the user who executes the target Java application process, while mackerel-agent usually runs under root privilege.
Since the executing user may not be root, you are required to specify the user in `mackerel-agent-conf` as shown above.

## About the `PerfDisableSharedMem` JVM option issue

Since there is a performance issue called [the four month bug](https://www.evanjones.ca/jvm-mmap-pause.html), several middlewares specify the `-XX:+PerfDisableSharedMem` JVM option as default.
When the JVM option is enabled, this plugin is no longer able to work because which depends `jps` and `jstat` JDK tools.

## References

- https://github.com/sensu/sensu-community-plugins/blob/master/plugins/java/jstat-metrics.py
- http://docs.oracle.com/javase/7/docs/technotes/tools/share/jstat.html
- https://github.com/kazeburo/jstat2gf
