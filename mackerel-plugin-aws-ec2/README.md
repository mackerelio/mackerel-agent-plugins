mackerel-plugin-aws-ec2
=======================

AWS EC2 custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-aws-ec2 [-instance-id=<id>] [-region=<aws-region>] [-access-key-id=<id>] [-secret-access-key=<key>] [-tempfile=<tempfile>]
```
* if you run on an ec2-instance, you don't have to specify `-instance-id` & `-region`
* if you configure credentials (by using the `~/.aws/credentials` file or by setting the environment variables), you don't have to specify `-access-key-id` & `-secret-access-key`

※ For more information about credentials, see the [AWS SDK for Go](https://github.com/aws/aws-sdk-go#configuring-credentials).

## AWS IAM Policy
the credential provided manually or fetched automatically by IAM Role should have the policy that includes an action, 'cloudwatch:GetMetricStatistics'

## Example of mackerel-agent.conf

```
[plugin.metrics.aws-ec2]
command = "/path/to/mackerel-plugin-aws-ec2"
```
