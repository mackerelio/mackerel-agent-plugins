mackerel-plugin-td-table-count
=====================

Treasure Data Table Count custom metrics plugin for mackerel-agent.

## Usage

```shell
mackerel-plugin-td-table-count -api-key=<Master API Key> -database=<Database Name>
```

## Example of mackerel-agent.conf

```
[plugin.metrics.td-table-count]
command = "/path/to/mackerel-plugin-td-table-count -api-key=<Master API Key> -database=<Database Name>"
```

## Author

[Takuya Arita](https://github.com/ariarijp)
