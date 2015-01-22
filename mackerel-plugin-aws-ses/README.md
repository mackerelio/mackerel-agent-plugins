mackerel-plugin-aws-ses
=================================

AWS SES custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-aws-ses -endpoint=<SES Endpoint URL> [-access-key-id=<id>] [-secret-access-key==<key>] [-tempfile=<tempfile>]
```
* SES Endpoint URL should be like "https://email.#{AWS_REGION}.amazonaws.com" (starting with "https://"). see "API (HTTPS) endpoint" column of http://docs.aws.amazon.com/ses/latest/DeveloperGuide/regions.html
* currently, even if you run this on an ec2-instance and the instance is associated with an appropriate IAM Role, this plugin can't fetch any values because of authentication failure. so you need to manually create and specify access-key-id and secret-access-key.
* you can avoid specifying keys as command-line args. see `func GetAuth` in https://github.com/crowdmob/goamz/blob/master/aws/aws.go

## AWS Policy
the credential provided manually should have the policy that includes actions, 'ses:GetSendQuota' and 'ses:GetSendStatistics'

## Example of mackerel-agent.conf
```
[plugin.metrics.aws-ses]
command = "/path/to/mackerel-plugin-aws-ses -endpoint=https://email.us-west-2.amazonaws.com"
```
