mackerel-plugin-munin
=======================

munin (wrapper) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-munin -plugin=<munin-plugin-executable> [-plugin-conf-d=<munin-plugin-conf-dir>] [-name=<mackerel-metric-name>] [-tempfile=<tempfile>]
```

* this executes the munin-plugin, first with an `config` argument so that this will get graph definitions, then with no argument so that this will get values.
* when `-plugin-conf-d` specified, this reads `env.KEY VALUE` entries from files of the dir and set those as environment variables before plugin executions. (other kinds of plugin-conf-entry are not implemented)

## Example of mackerel-agent.conf

```
[plugin.metrics.nfsd]
command = "/path/to/mackerel-plugin-munin -plugin=/usr/share/munin/plugins/nfsd"
```

```
[plugin.metrics.bind9]
command = "/path/to/mackerel-plugin-munin -plugin=/etc/munin/plugins/bind9 -plugin-conf-d=/etc/munin/plugin-conf.d"
```

```
[plugin.metrics.postfix]
command = "MUNIN_LIBDIR=/usr/share/munin /path/to/mackerel-plugin-munin -plugin=/usr/share/munin/plugins/postfix_mailqueue -name=postfix.mailqueue"
```
(some munin-plugins sources `$MUNIN_LIBDIR/plugins/plugin.sh`)
