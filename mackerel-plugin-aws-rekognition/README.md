mackerel-plugin-aws-rekognition
=======================

AWS Rekognition custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-aws-rekognition [-operation=<aws-rekognition-operation>] [-region=<aws-region>] [-access-key-id=<id>] [-secret-access-key=<key>] [-tempfile=<tempfile>]
```

## AWS IAM Policy
the credential provided manually or fetched automatically by IAM Role should have the policy that includes an action, 'cloudwatch:GetMetricStatistics'

## Example of mackerel-agent.conf

```
[plugin.metrics.aws-rekognition]
command = "/path/to/mackerel-plugin-aws-rekognition"
```
