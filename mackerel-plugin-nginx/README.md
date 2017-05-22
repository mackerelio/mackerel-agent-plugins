mackerel-plugin-nginx
=====================

Nginx custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-nginx [-header=<header>] [-host=<host>] [-path=<path>] [-port=<port>] [-scheme=<'http'|'https'>] [-tempfile=<tempfile>] [-uri=<uri>]
```

## Requirements

- [ngx_http_stub_status_module](http://nginx.org/en/docs/http/ngx_http_stub_status_module.html)

## Example of mackerel-agent.conf

```
[plugin.metrics.nginx]
command = "/path/to/mackerel-plugin-nginx"
```
