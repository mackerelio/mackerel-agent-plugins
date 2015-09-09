mackerel-plugin-docker
=========================

Docker (https://www.docker.com/) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-docker [-host=<host>] [-port=<port>] [-tempfile=<tempfile>]
```

`-host` and `-port` options are not implemented yet.

## Example of mackerel-agent.conf

```
[plugin.metrics.memcached]
command = "/path/to/mackerel-plugin-docker"
```
