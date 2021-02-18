mackerel-plugin-aws-ses
=================================

AWS SES custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-aws-ses -region=<SES Region> [-access-key-id=<id>] [-secret-access-key=<key>] [-tempfile=<tempfile>]
```
* You may omit region parameter if you're running this plugin on an EC2 instance running in same region with the target SES
* if you run on an ec2-instance and the instance is associated with an appropriate IAM Role, you probably don't have to specify `-access-key-id` & `-secret-access-key`

## AWS Policy
the credential provided manually or fetched automatically by IAM Role should have the policy that includes actions, 'ses:GetSendQuota' and 'ses:GetSendStatistics'

## Example of mackerel-agent.conf
```
[plugin.metrics.aws-ses]
command = "/path/to/mackerel-plugin-aws-ses -region=us-west-2"
```
