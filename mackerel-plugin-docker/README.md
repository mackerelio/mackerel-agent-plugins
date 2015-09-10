mackerel-plugin-docker
=========================

Docker (https://www.docker.com/) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-docker [-host=<host>] [-command=<docker>] [-tempfile=<tempfile>]
```

- `-host` Socket path. This option is same as `--host` option of docker command.
- `-command` Path to docker command. Without path, binary is searched in the directories named by the PATH environment variable. The default value is `docker`.
- `-tempfile` Temporary file stored metric values for calculating differentials.

## Example of mackerel-agent.conf

```
[plugin.metrics.memcached]
command = "/path/to/mackerel-plugin-docker"
```
