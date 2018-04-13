# mackerel-plugin-aws-step-function
AWS Step Functions custom metrics plugin.

## Synopsis

```shell
mackerel-plugin-aws-step-functions [-access-key-id=<id>] [-secret-access-key=<key>] [-region=<region>] [-state-machine-arn=<arn>] [-metric-key-prefix=<prefix>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.aws-step-functions]
command = "/path/to/mackerel-plugin-aws-step-functions"
```
