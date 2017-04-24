mackerel-plugin-uptimerobot
=====================

[Uptime Robot](https://uptimerobot.com) plugin for mackerel-agent.

## Usage

```shell
./mackerel-plugin-uptimerobot -api-key=<API Key> -monitor-id=<Monitor ID> [-name=<Friendly Name>]
```

NOTICE: `-name` is optional. It is not related to Monitor's Friendly Name.

## Example of mackerel-agent.conf

```
[plugin.metrics.uptimerobot]
command = "/path/to/mackerel-plugin-uptimerobot -api-key=<API Key> -monitor-id=<Monitor ID>
```

## Author

[Takuya Arita](https://github.com/ariarijp)
