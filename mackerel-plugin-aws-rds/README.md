mackerel-plugin-aws-rds
=======================

AWS RDS custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-aws-rds -identifier=<db-instance-identifer> [-region=<aws-region>] [-access-key-id=<id>] [-secret-access-key=<key>] [-tempfile=<tempfile>] [-engine=<mysql or aurora or mariadb or postgresql>] [-metric-key-prefix=<prefix>] [-metric-label-prefix=<label-prefix>]
```
* if you run on an ec2-instance, you probably don't have to specify `-region`
* if you run on an ec2-instance and the instance is associated with an appropriate IAM Role, you probably don't have to specify `-access-key-id` & `-secret-access-key`

## AWS IAM Policy
the credential provided manually or fetched automatically by IAM Role should have the policy that includes an action, 'cloudwatch:GetMetricStatistics'

## Example of mackerel-agent.conf

```
[plugin.metrics.aws-rds]
command = "/path/to/mackerel-plugin-aws-rds -identifier=mysql01 -engine=mysql"
```
