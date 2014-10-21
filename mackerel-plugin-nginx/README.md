mackerel-plugin-nginx
=====================

Nginx custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-nginx [-host=<host>] [-port=<port>] [-path=<path>] [-tempfile=<tempfile>]
```

## Requirements

- [ngx_http_status_module](http://nginx.org/en/docs/http/ngx_http_status_module.html)

## Example of mackerel-agent.conf

```
[plugin.metrics.nginx]
command = "/path/to/mackerel-plugin-nginx"
```
