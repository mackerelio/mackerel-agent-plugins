mackerel-plugin-aws-elasticsearch
=======================

AWS Elasticsearch Service custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-aws-elasticsearch -domain=<aws-elasticsearch-domain> -client-id=<aws-client-id> [-region=<aws-region>] [-access-key-id=<aws-access-key-id>] [-secret-access-key=<aws-secret-access-key>] [-metric-key-prefix=<prefix>] [-metric-label-prefix=<prefix>] [-tempfile=<tmpfile>]
```

## AWS IAM Policy
the credential provided manually or fetched automatically by IAM Role should have the policy that includes an action, 'cloudwatch:GetMetricStatistics'

## Example of mackerel-agent.conf

```
[plugin.metrics.aws-elasticsearch]
command = "/path/to/mackerel-plugin-aws-elasticsearch -domain=your-es-domain -client-id=your-aws-client-id"
```
