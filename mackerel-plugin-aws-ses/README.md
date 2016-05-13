mackerel-plugin-aws-ses
=================================

AWS SES custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-aws-ses -endpoint=<SES Endpoint URL> [-access-key-id=<id>] [-secret-access-key=<key>] [-tempfile=<tempfile>]
```
* SES Endpoint URL should be like "https://email.#{AWS_REGION}.amazonaws.com" (starting with "https://"). see "API (HTTPS) endpoint" column of http://docs.aws.amazon.com/ses/latest/DeveloperGuide/regions.html
* if you run on an ec2-instance and the instance is associated with an appropriate IAM Role, you probably don't have to specify `-access-key-id` & `-secret-access-key`

## AWS Policy
the credential provided manually or fetched automatically by IAM Role should have the policy that includes actions, 'ses:GetSendQuota' and 'ses:GetSendStatistics'

## Example of mackerel-agent.conf
```
[plugin.metrics.aws-ses]
command = "/path/to/mackerel-plugin-aws-ses -endpoint=https://email.us-west-2.amazonaws.com"
```
