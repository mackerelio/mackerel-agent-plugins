mackerel-plugin-haproxy
=====================

HAProxy custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-haproxy [-host=<host>] [-port=<port>] [-path=<stats-path>] [-scheme=<http|https>] [-username=<username] [-password=<password>] [-tempfile=<tempfile>]
or
mackerel-plugin-haproxy [-uri=<uri>] [-username=<username] [-password=<password>] [-tempfile=<tempfile>]
```

For Basic Auth, set username.

## Example of mackerel-agent.conf

```
[plugin.metrics.haproxy]
command = "/path/to/mackerel-plugin-haproxy -port=8000"
```

## Example of haproxy.conf

This plugin requires to enable stats of haproxy.
Example configuration is as follow.

```
listen hastats *:8088
    mode http
    maxconn 64
    timeout connect 5000
    timeout client 10000
    timeout server 10000
    stats enable
    stats show-legends
    stats uri /haproxy?hastats

    # basic auth
    stats auth admin:adminadmin
```

See haproxy_test.go for example configuration.