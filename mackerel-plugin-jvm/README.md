mackerel-plugin-jvm
===================

JVM(jstat) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-jvm -javaname=<javaname> [-pidfile=</path/to/pidfile>] [-jstatpath=</path/to/jstat] [-jpspath=/path/to/jps] [-jinfopath=/path/to/jinfo] [-host=<host>] [-port=<port>]
```

## Requirements

- JVM 1.6 or higher

## Example of mackerel-agent.conf

```
[plugin.metrics.jvm]
command = "/path/to/mackerel-plugin-jvm -javaname=NettyServer -jstatpath=/usr/bin/jstat -jpspath=/usr/bin/jps -jinfopath=/usr/bin/jinfo"
user = "SOME_USER_NAME"
```

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

## References

- https://github.com/sensu/sensu-community-plugins/blob/master/plugins/java/jstat-metrics.py
- http://docs.oracle.com/javase/7/docs/technotes/tools/share/jstat.html
- https://github.com/kazeburo/jstat2gf
