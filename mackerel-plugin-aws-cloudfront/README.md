mackerel-plugin-aws-cloudfront
==============================

AWS CloudFront custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-aws-cloudfront -identifier=<cloudfront-distribution-id> [-access-key-id=<id>] [-secret-access-key=<key>] [-tempfile=<tempfile>] [-metric-key-prefix=<cloudfront>] [-metric-label-prefix=<CloudFront>]
```

* if you run on an ec2-instance and the instance is associated with an appropriate IAM Role, you probably don't have to specify `-access-key-id` & `-secret-access-key`

## AWS IAM Policy

the credential provided manually or fetched automatically by IAM Role should have the policy that includes an action, 'cloudwatch:GetMetricStatistics'

## Example of mackerel-agent.conf

```
[plugin.metrics.aws-cloudfront]
command = "/path/to/mackerel-plugin-aws-cloudfront -identifier=yourdistributionid"
```
