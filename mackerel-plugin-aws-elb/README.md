mackerel-plugin-aws-elb
=======================

AWS ELB custom metrics plugin for mackerel.io agent.
As it stands, this can fetch only across-all-LBs metrics.

## Synopsis

```shell
mackerel-plugin-aws-elb [-region=<aws-region>] [-access-key-id=<id>] [-secret-access-key==<key>] [-tempfile=<tempfile>]
```
* if you run on an ec2-instance, you probably don't have to specify `-region`
* if you run on an ec2-instance and the instance is associated with an appropriate IAM Role, you probably don't have to specify `-access-key-id` & `-secret-access-key`

## AWS IAM Policy
the credential provided manually or fetched automatically by IAM Role should have the policy that includes actions, 'cloudwatch:GetMetricStatistics' and 'cloudwatch:ListMetrics'

## Example of mackerel-agent.conf

```
[plugin.metrics.aws-elb]
command = "/path/to/mackerel-plugin-aws-elb"
```
