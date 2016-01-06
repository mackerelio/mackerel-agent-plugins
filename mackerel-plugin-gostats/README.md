mackerel-plugin-gostats
=====================

gostats custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-gostats [-host=<host>] [-port=<port>] [-path=<path>] [-scheme=<http|https>] [-uri=<URI>] [-metric-key-prefix=gostats]
```

## Requirements

This plugin requires [github.com/fukata/golang-stats-api-handler](https://github.com/fukata/golang-stats-api-handler)
Enable `github.com/fukata/golang-stats-api-handler` as below.

```
import (
    "net/http"
    "log"
    "github.com/fukata/golang-stats-api-handler"
)
func main() {
    http.HandleFunc("/api/stats", stats_api.Handler)
    log.Fatal( http.ListenAndServe(":8080", nil) )
}
```

## Example of mackerel-agent.conf

```
[plugin.metrics.gostats]
command = "/path/to/mackerel-plugin-gostats -port=8000 -path=/api/stats"
```
