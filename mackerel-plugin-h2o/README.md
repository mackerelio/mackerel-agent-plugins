mackerel-plugin-h2o
===================

H2O custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-h2o [-header=<header>] [-host=<host>] [-path=<path>] [-port=<port>] [-scheme=<'http'|'https'>] [-tempfile=<tempfile>] [-uri=<uri>]
```

## Requirements

- [Status Directives \- Configure \- H2O \- the optimized HTTP/2 server](https://h2o.examp1e.net/configure/status_directives.html#)

## Example of mackerel-agent.conf

```
[plugin.metrics.h2o]
command = "/path/to/mackerel-plugin-h2o"
```

## Graphs and Metrics

### h2o.uptime

- h2o.uptime.uptime

### h2o.connections

- h2o.connections.max-connections
- h2o.connections.connections

### h2o.listeners

- h2o.listeners.listeners

### h2o.worker-threads

- h2o.worker-threads.worker-threads

### h2o.num-sessions

- h2o.num-sessions.num-sessions

### h2o.requests

- h2o.requests.requests

### h2o.status-errors

- h2o.status-errors.status-errors_400
- h2o.status-errors.status-errors_403
- h2o.status-errors.status-errors_404
- h2o.status-errors.status-errors_405
- h2o.status-errors.status-errors_416
- h2o.status-errors.status-errors_417
- h2o.status-errors.status-errors_500
- h2o.status-errors.status-errors_502
- h2o.status-errors.status-errors_503

### h2o.http2-errors

- h2o.http2-errors.http2-errors_protocol
- h2o.http2-errors.http2-errors_internal
- h2o.http2-errors.http2-errors_flow-control
- h2o.http2-errors.http2-errors_settings-timeout
- h2o.http2-errors.http2-errors_frame-size
- h2o.http2-errors.http2-errors_refused-stream
- h2o.http2-errors.http2-errors_cancel
- h2o.http2-errors.http2-errors_compression
- h2o.http2-errors.http2-errors_connect
- h2o.http2-errors.http2-errors_enhance-your-calm
- h2o.http2-errors.http2-errors_inadequate-security

### h2o.read-closed

- h2o.read-closed.http2_read-closed

### h2o.write-closed

- h2o.write-closed.http2_write-closed

### h1o.connect-time

- h2o.connect-time.connect-time-0
- h2o.connect-time.connect-time-25
- h2o.connect-time.connect-time-50
- h2o.connect-time.connect-time-75
- h2o.connect-time.connect-time-99

### h2o.header-time

- h2o.header-time.header-time-0
- h2o.header-time.header-time-25
- h2o.header-time.header-time-50
- h2o.header-time.header-time-75
- h2o.header-time.header-time-99

### h2o.body-time

- h2o.body-time.body-time-0
- h2o.body-time.body-time-25
- h2o.body-time.body-time-50
- h2o.body-time.body-time-75
- h2o.body-time.body-time-99

### h2o.request-total-time

- h2o.request-total-time.request-total-time-0
- h2o.request-total-time.request-total-time-25
- h2o.request-total-time.request-total-time-50
- h2o.request-total-time.request-total-time-75
- h2o.request-total-time.request-total-time-99

### h2o.process-time

- h2o.process-time.process-time-0
- h2o.process-time.process-time-25
- h2o.process-time.process-time-50
- h2o.process-time.process-time-75
- h2o.process-time.process-time-99

### h2o.response-time

- h2o.response-time.response-time-0
- h2o.response-time.response-time-25
- h2o.response-time.response-time-50
- h2o.response-time.response-time-75
- h2o.response-time.response-time-99

### h2o.duration

- h2o.duration.duration-0
- h2o.duration.duration-25
- h2o.duration.duration-50
- h2o.duration.duration-75
- h2o.duration.duration-99
