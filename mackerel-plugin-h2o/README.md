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

- h2o.connections.max_connections
- h2o.connections.connections

### h2o.listeners

- h2o.listeners.listeners

### h2o.worker_threads

- h2o.worker_threads.worker_threads

### h2o.num_sessions

- h2o.num_sessions.num_sessions

### h2o.requests

- h2o.requests.requests

### h2o.status_errors

- h2o.status_errors.status_errors_400
- h2o.status_errors.status_errors_403
- h2o.status_errors.status_errors_404
- h2o.status_errors.status_errors_405
- h2o.status_errors.status_errors_416
- h2o.status_errors.status_errors_417
- h2o.status_errors.status_errors_500
- h2o.status_errors.status_errors_502
- h2o.status_errors.status_errors_503

### h2o.http2_errors

- h2o.http2_errors.http2_errors_protocol
- h2o.http2_errors.http2_errors_internal
- h2o.http2_errors.http2_errors_flow_control
- h2o.http2_errors.http2_errors_settings_timeout
- h2o.http2_errors.http2_errors_frame_size
- h2o.http2_errors.http2_errors_refused_stream
- h2o.http2_errors.http2_errors_cancel
- h2o.http2_errors.http2_errors_compression
- h2o.http2_errors.http2_errors_connect
- h2o.http2_errors.http2_errors_enhance_your_calm
- h2o.http2_errors.http2_errors_inadequate_security

### h2o.read_closed

- h2o.read_closed.http2_read_closed

### h2o.write_closed

- h2o.write_closed.http2_write_closed

### h1o.connect_time

- h2o.connect_time.connect_time_0
- h2o.connect_time.connect_time_25
- h2o.connect_time.connect_time_50
- h2o.connect_time.connect_time_75
- h2o.connect_time.connect_time_99

### h2o.header_time

- h2o.header_time.header_time_0
- h2o.header_time.header_time_25
- h2o.header_time.header_time_50
- h2o.header_time.header_time_75
- h2o.header_time.header_time_99

### h2o.body_time

- h2o.body_time.body_time_0
- h2o.body_time.body_time_25
- h2o.body_time.body_time_50
- h2o.body_time.body_time_75
- h2o.body_time.body_time_99

### h2o.request_total_time

- h2o.request_total_time.request_total_time_0
- h2o.request_total_time.request_total_time_25
- h2o.request_total_time.request_total_time_50
- h2o.request_total_time.request_total_time_75
- h2o.request_total_time.request_total_time_99

### h2o.process_time

- h2o.process_time.process_time_0
- h2o.process_time.process_time_25
- h2o.process_time.process_time_50
- h2o.process_time.process_time_75
- h2o.process_time.process_time_99

### h2o.response_time

- h2o.response_time.response_time_0
- h2o.response_time.response_time_25
- h2o.response_time.response_time_50
- h2o.response_time.response_time_75
- h2o.response_time.response_time_99

### h2o.duration

- h2o.duration.duration_0
- h2o.duration.duration_25
- h2o.duration.duration_50
- h2o.duration.duration_75
- h2o.duration.duration_99
