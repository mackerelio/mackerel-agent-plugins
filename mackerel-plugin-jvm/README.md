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

Normally mackerel-agent is executed by root, but this plugin (and jps command explained above) needs to be executed by the user who executes target java application process, and the two may differ.
In this case, you need to specify which user to execute this plugin to `mackerel-agent-conf` configuration, like the example.

## References

- https://github.com/sensu/sensu-community-plugins/blob/master/plugins/java/jstat-metrics.py
- http://docs.oracle.com/javase/7/docs/technotes/tools/share/jstat.html
- https://github.com/kazeburo/jstat2gf
