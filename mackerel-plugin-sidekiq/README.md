mackerel-plugin-sidekiq
====

## Overview
Sidekiq custom metrics plugin for mackerel.io agent.

## Description

This plugin makes two graphs: one shows processed job diff and failed job diff, and another one shows number of busy, enqueued, scheduled, retry and dead jobs.

## Usage

```
mackerel-plugin-sidekiq [-host=<host>] [-port=<port>] [-password=<password>] [-db=<db>] [-tempfile=<template file path>] [-redis-namespace=<redis namespace>]
```

### Example of mackerel-agent.conf

```
[plugin.metrics.sidekiq]
command = "/path/to/mackerel-plugin-sidekiq"
```

## Graphs and Metrics

### Sidekiq processed and failed count (sidekiq.ProcessedANDFailed)

* processed (sidekiq.ProcessedANDFailed.processed)
* failed (sidekiq.ProcessedANDFailed.failed)

### Sidekiq stats (sidekiq.Stats)

* busy (sidekiq.Stats.busy)
* enqueued (sidekiq.Stats.enqueued)
* scheduled (sidekiq.Stats.scheduled)
* retry (sidekiq.Stats.retry)
* dead (sidekiq.Stats.dead)

### sidekiq queue latency (sidekiq.QueueLatency)

* Latency sec (sidekiq.QueueLatency.#)

## Author

[littlekbt](https://github.com/littlekbt)
