mackerel-plugin-xentop
======================

Xen custom metrics plugin for mackerel.io agent.
Xen metrics are retrieved by `xentop` command.

## Synopsis

```shell
mackerel-plugin-xentop [-xenversion=<xenversion>] [-tempfile=<tempfile>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.xentop]
command = "/path/to/mackerel-plugin-xentop"
```

