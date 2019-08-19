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

## Author

[littlekbt](https://github.com/littlekbt)
