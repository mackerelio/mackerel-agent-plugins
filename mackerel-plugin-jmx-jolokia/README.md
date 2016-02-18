mackerel-plugin-jmx-jolokia
===========================

Jolokia (https://jolokia.org/) custom metrics plugin for mackerel.io agent

## Synopsis

```shell
mackerel-plugin-jmx-jolokia [-host=<host>] [-port=<port>] [-tempfile=<tempfile>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.jolokia]
command = "/path/to/mackerel-plugin-jmx-jolokia"
```

## Example of jolokia response

```
curl -s http://127.0.0.1:8778/jolokia/read/java.lang:type=Memory/HeapMemoryUsage | jq .
{
  "request": {
    "mbean": "java.lang:type=Memory",
    "type": "read"
  },
  "value": {
    "ObjectPendingFinalizationCount": 0,
    "Verbose": false,
    "HeapMemoryUsage": {
      "init": 1073741824,
      "committed": 1069023232,
      "max": 1069023232,
      "used": 994632048
    },
    "NonHeapMemoryUsage": {
      "init": 2555904,
      "committed": 44040192,
      "max": 1350565888,
      "used": 43070016
    },
    "ObjectName": {
      "objectName": "java.lang:type=Memory"
    }
  },
  "timestamp": 1455079714,
  "status": 200
}
```