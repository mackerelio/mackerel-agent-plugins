mackerel-plugin-unicorn
=======================

Unicorn custom metrics plugin for mackerel.io agent.

Synopsis
--------

```sh
mackerel-plugin-unicorn [-pidfile=<path>] [-tempfile=<tempfile>] [-metric-key-prefix=<prefix>]
```

Example of mackerel-agent.conf
------------------------------

```conf
[plugin.metrics.unicorn]
command = "/path/to/mackerel-plugin-unicorn -pidfile=/var/www/app/shared/tmp/pids/unicorn.pid"
```
