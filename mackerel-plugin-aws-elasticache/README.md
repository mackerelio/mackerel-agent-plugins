mackerel-plugin-aws-elasticache
=======================

AWS ElastiCache custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-aws-elasticache -elasticache-type=<type> -cache-cluster-id=<cluster-id> [-cache-node-id=<node-id>] [-region=<aws-region>] [-access-key-id=<id>] [-secret-access-key=<key>] [-tempfile=<tempfile>]
```
* if you run on an ec2-instance, you probably don't have to specify `-region`
* if you run on an ec2-instance and the instance is associated with an appropriate IAM Role, you probably don't have to specify `-access-key-id` & `-secret-access-key`

## AWS IAM Policy
the credential provided manually or fetched automatically by IAM Role should have the policy that includes an action, 'cloudwatch:GetMetricStatistics'

## Example of mackerel-agent.conf

```
[plugin.metrics.aws-elasticache]
command = "/path/to/mackerel-plugin-aws-elasticache -elasticache-type=memcached -cache-cluster-id=elasticache01"
```
