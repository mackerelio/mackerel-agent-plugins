mackerel-plugin-mysqlrouter
============================
[![Build Status](https://cloud.drone.io/api/badges/rluisr/mysqlrouter_exporter/status.svg)](https://cloud.drone.io/rluisr/mysqlrouter_exporter)

Usage
-----
```
[plugin.metrics.mysqlrouter]
command = "MYSQLROUTER_URL=http://localhost:8080 MYSQLROUTER_USER=luis MYSQLROUTER_PASS=luis /path/to/mackerel-plugin-mysqlrouter"
```

FYI
---
You should set alert of health per route like below.

![](https://raw.githubusercontent.com/rluisr/image-store/master/mackerel-plugin-mysqlrouter/mackerel-plugin-mysqlrouter01.png)

This `0.1` is no meaning. If route is no heath return `0`.

Todo
----
- [ ] Show graph per route connections