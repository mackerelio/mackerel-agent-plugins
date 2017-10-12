mackerel-plugin-json
====================

Json custom metrics plugin for mackerel.io agent.

## Synopsis

```
mackerel-plugin-json -url=<url to get JSON> [-prefix=<prefix for a metric name>] [-insecure] [-include=<expression>] [-exclude=<expression>]
```

- `-url` always needs to be specified.
- If you want to skip a SSL certificate validation (e.g. using a self-signed certificate), you need to specify `-insecure`.
- If you want to get only metrics that matches a regular expression, specify the expression with `-include`. If you want to get only ones that doesn't match an expression, use `-exclude`.
- Metrics that have non-number value (string, `null`, etc.) are omitted.
- Arrays are supported. Metrics' name will contain a serial number like `custom.elements.0.xxx`, `custom.elements.1.xxx` and so on.

## Example of mackerel-agent.conf

```
[plugin.metrics.jolokia]
command = mackerel-plugin-json -url='http://127.0.0.1:8778/jolokia/read/java.lang:type=Memory' -include='HeapMemoryUsage' -prefix='custom.jolokia'
```

## Examples of output from the plugin

Given the following JSON:

```
curl -s http://127.0.0.1:8778/jolokia/read/java.lang:type=Memory | jq .
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

The output should be:

```
# Option : -url='http://127.0.0.1:8778/jolokia/read/java.lang:type=Memory'

custom.value.HeapMemoryUsage.used       994632048.000000        1503586907
custom.value.HeapMemoryUsage.init       1073741824.000000       1503586907
custom.value.HeapMemoryUsage.max        1069023232.000000       1503586907
custom.timestamp        1455079714.000000       1503586907
custom.status   200.000000      1503586907
custom.value.NonHeapMemoryUsage.committed       44040192.000000 1503586907
custom.value.NonHeapMemoryUsage.max     1350565888.000000       1503586907
custom.value.NonHeapMemoryUsage.used    43070016.000000 1503586907
custom.value.NonHeapMemoryUsage.init    2555904.000000  1503586907
custom.value.HeapMemoryUsage.committed  1069023232.000000       1503586907
```

```
# Option : -url='http://127.0.0.1:8778/jolokia/read/java.lang:type=Memory' -include='\.HeapMemoryUsage'

custom.value.HeapMemoryUsage.init       1073741824.000000       1503587081
custom.value.HeapMemoryUsage.committed  1069023232.000000       1503587081
custom.value.HeapMemoryUsage.max        1069023232.000000       1503587081
custom.value.HeapMemoryUsage.used       994632048.000000        1503587081
```

```
# Option : -url='http://127.0.0.1:8778/jolokia/read/java.lang:type=Memory' -exclude='(timestamp|status)'

custom.value.HeapMemoryUsage.init       1073741824.000000       1503587166
custom.value.HeapMemoryUsage.committed  1069023232.000000       1503587166
custom.value.HeapMemoryUsage.max        1069023232.000000       1503587166
custom.value.HeapMemoryUsage.used       994632048.000000        1503587166
custom.value.NonHeapMemoryUsage.max     1350565888.000000       1503587166
custom.value.NonHeapMemoryUsage.used    43070016.000000 1503587166
custom.value.NonHeapMemoryUsage.init    2555904.000000  1503587166
custom.value.NonHeapMemoryUsage.committed       44040192.000000 1503587166
```

You can also get metrics from any APIs that returns a JSON. For example:

```
# Option : -url='https://[your-github-token]@api.github.com/repos/doublemarket/private-repo' -include='(_count|watchers|issues)' -prefix='custom.github.some-private-repo'

custom.github.private-repo.open_issues_count    1171.000000     1503587879
custom.github.private-repo.open_issues  1171.000000     1503587879
custom.github.private-repo.network_count        15042.000000    1503587879
custom.github.private-repo.forks_count  15042.000000    1503587879
custom.github.private-repo.stargazers_count     36733.000000    1503587879
custom.github.private-repo.subscribers_count    2444.000000     1503587879
custom.github.private-repo.watchers_count       36733.000000    1503587879
custom.github.private-repo.watchers     36733.000000    1503587879
```
