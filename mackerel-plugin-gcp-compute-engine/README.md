mackerel-plugin-gcp-compute-engine
====

## Overview

mackerel-plugin-gcp-compute-engine is mackerel plugin that collects metrics from Google Compute Engine by using [Google Monitoring API](https://cloud.google.com/monitoring/api/v3/).

## Description

This plugin uses Google Monitoring API to get more metrics than mackerel-agent defalut metrics.


Caution: This Plugin works only on Compute Engien instance that is enabled Stackdriver Monitoring API full access.

## Usage

```shell
mackerel-plugin-gcp-compute-engine -api-key=<api key> [-project=<project id or number>] [-instance-name=<target instance name>]
```

If `-project` or `-instance-name` are not specified, they are obtained from Google Compute Engine Metadata API for the instance executing this plugin.
It means you don't need to specify them to monitor the instance itself.

### Example of mackerel-agent.conf

```
[plugin.metrics.gcp-compute-engine]
command = "/path/to/mackerel-plugin-gcp-compute-engine -api-key=<YOUR-API-KEY>"
```

## Author

[littlekbt](https://github.com/littlekbt)
