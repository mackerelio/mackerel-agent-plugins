mackerel-plugin-docker
=========================

Docker (https://www.docker.com/) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-docker [-host=<host>] [-tempfile=<tempfile>] [-name-format=<format>] [-label=<key>]
```

- `-host` Socket path. This option is same as `--host` option of docker command. The default value is `unix:///var/run/docker.sock`.
- `-tempfile` Temporary file stored metric values for calculating differentials.
- `-name-format` Set the name format from name, name_id, id, image, image_id, image_name or label (default "name_id")
- `-label` Use the value of the key as name in case that name-format is label.

## Current Status

* x: works

|Metric name                           |cgroup v1| cgroup v2|
|--------------------------------------|---------|----------|
|docker.cpuacct_percentage.#.user      | x       | x        |
|docker.cpuacct_percentage.#.system    | x       | x        |
|docker.cpuacct.#.user                 | x       | x        |
|docker.cpuacct.#.system               | x       | x        |
|docker.memory.#.cache                 | x       | x        |
|docker.memory.#.rss                   | x       | x        |
|docker.blkio.io_queued.#.read         | x       |          |
|docker.blkio.io_queued.#.write        | x       |          |
|docker.blkio.io_queued.#.sync         | x       |          |
|docker.blkio.io_queued.#.async        | x       |          |
|docker.blkio.io_serviced.#.read       | x       |          |
|docker.blkio.io_serviced.#.write      | x       |          |
|docker.blkio.io_serviced.#.sync       | x       |          |
|docker.blkio.io_serviced.#.async      | x       |          |
|docker.blkio.io_service_bytes.#.read  | x       |          |
|docker.blkio.io_service_bytes.#.write | x       |          |
|docker.blkio.io_service_bytes.#.sync  | x       |          |
|docker.blkio.io_service_bytes.#.async | x       |          |

## Example of mackerel-agent.conf

```
[plugin.metrics.docker]
command = "/path/to/mackerel-plugin-docker"
```
