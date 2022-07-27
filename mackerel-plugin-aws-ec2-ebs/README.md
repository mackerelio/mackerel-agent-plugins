mackerel-plugin-aws-ec2-ebs
=================================

AWS EC2 EBS custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-aws-ec2-ebs [-instance-id=<id>] [-region=<aws-region>] [-access-key-id=<id>] [-secret-access-key=<key>] [-tempfile=<tempfile>]
```
* collect data from all volumes which attached to the instance
* if you run on an ec2-instance, you probably don't have to specify `-instance-id` & `-region`
* if you run on an ec2-instance and the instance is associated with an appropriate IAM Role, you probably don't have to specify `-access-key-id` & `-secret-access-key`
* you can set keys by environment variables: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` (see https://github.com/aws/aws-sdk-go#configuring-credentials)

## AWS IAM Policy
the credential provided manually or fetched automatically with IAM Role, should have the policy that includes an action, `cloudwatch:GetMetricStatistics` and `ec2:DescribeVolumes`

## Example of mackerel-agent.conf

```
[plugin.metrics.aws-ec2_ebs]
command = ["/path/to/mackerel-plugin-aws-ec2-ebs"]
```
