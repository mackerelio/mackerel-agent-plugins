# mackerel-plugin-rack-stats
Rack servers metrics plugin for mackerel.io agent.  
Only work on Linux.  

## Synopsis

```shell
mackerel-plugin-rack-stats [-address=<url or unix domain socket>] [-path=<path>] [-metric-key=<metric-key>]  [-tempfile=<tempfile>]
```

```shell
$ ./mackerel-plugin-rack-stats --help
Usage of ./mackerel-plugin-rack-stats:
  -address string
        URL or Unix Domain Socket (default "http://localhost:8080")
  -metric-key string
        Metric Key
  -path string
        Path (default "/_raindrops")
  -tempfile string
        Temp file name
  -version
        Version
```

## Requirements

- [raindrops](https://rubygems.org/gems/raindrops)

## Example of mackerel-agent.conf

```
[plugin.metrics.rack_stats]
command = "/path/to/mackerel-plugin-rack-stats -address=unix:/path/to/unicorn.sock"
```

```
[plugin.metrics.rack_stats]
command = "/path/to/mackerel-plugin-rack-stats -address=http://localhost:8080"
```
