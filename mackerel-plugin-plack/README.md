mackerel-plugin-plack
=====================

Plack custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-plack [-host=<host>] [-port=<port>] [-path=<path?json>] [-scheme=<http|https>]
```

## Requirements

This plugin requires [Plack::Middleware::ServerStatus::Lite](https://metacpan.org/release/Plack-Middleware-ServerStatus-Lite) > 0.21.
Enable `ServerStatus::Lite` as below.

```perl
use Plack::Builder;
builder {
    enable "Plack::Middleware::ServerStatus::Lite",
    path => '/server-status',
    allow => [ '127.0.0.1' ],
    counter_file => '/tmp/counter_file',
    scoreboard => '/var/run/server';
    $app;
};
```

## Example of mackerel-agent.conf

```
[plugin.metrics.plack]
command = "/path/to/mackerel-plugin-plack -port=8000 -path='/status?auto'"
```
