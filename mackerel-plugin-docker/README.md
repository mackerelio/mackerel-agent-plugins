mackerel-plugin-docker
=========================

Docker (https://www.docker.com/) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-docker [-method=<method>] [-host=<host>] [-command=<docker>] [-tempfile=<tempfile>] [-name-format=<format>] [-label=<key>]
```

- `-method` Specify the method to collect stats, 'API' or 'File'. If not specified, a method is chosen based on docker API version. If the API version is under 1.17, 'File' is used. Otherwise, 'API' is used.
- `-host` Socket path. This option is same as `--host` option of docker command. The default value is `unix:///var/run/docker.sock`.
- `-command` Path to docker command. Without path, binary is searched in the directories named by the PATH environment variable. This is only used when method is 'File'. The default value is `docker`.
- `-tempfile` Temporary file stored metric values for calculating differentials.
- `-name-format` Set the name format from name, name_id, id, image, image_id, image_name or label (default "name_id")
- `-label` Use the value of the key as name in case that name-format is label.

## Example of mackerel-agent.conf

```
[plugin.metrics.docker]
command = "/path/to/mackerel-plugin-docker -method API"
```
