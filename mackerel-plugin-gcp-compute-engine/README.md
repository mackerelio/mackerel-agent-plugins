mackerel-plugin-gcp-compute-engine
====

## Overview

mackerel-plugin-gcp-compute-engine is mackerel plugin that collects metrics from Google Compute Engine by using [Google Monitoring API](https://cloud.google.com/monitoring/api/v3/).

## Description

This plugin uses Google Monitoring API to get more metrics than mackerel-agent defalut metrics.   


Caution: This Plugin works only on Compute Engien instance that is enabled Stackdriver Monitoring API full access. 

## Usage

```shell
mackerel-plugin-gcp-compute-engine [-api-key=<APIKey>] [-project=<project_id_or_number>] [-instance-name=<target instance name>]
```


### Example of mackerel-agent.conf

```
[plugin.metrics.gcp-compute-engine]
command = "/path/to/mackerel-plugin-gcp-compute-engine/main" ... //arguments
```

## Author

[littlekbt](https://github.com/littlekbt)
