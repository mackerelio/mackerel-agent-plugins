mackerel-plugin-proc-fd
===================

Opening fd custom metrics plugin for mackerel.io agent.  
This plugin monitors `/proc/[pid]/fd/*` to count file descriptors.  
Now this plugin posts only the maximum number of opening fd in matching processes.

## Synopsis

```shell
mackerel-plugin-proc-fd -process=<process name>
```

## Example of mackerel-agent.conf

```
[plugin.metrics.proc-fd]
command = "/path/to/mackerel-plugin-proc-fd -process='keepalived'"
```

