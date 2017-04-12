mackerel-plugin-aws-lambda
=================================

AWS Lambda custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-aws-lambda [-function-name=<function-name>] -region=<aws-region> -access-key-id=<id> -secret-access-key=<key> [-tempfile=<tempfile>] [-metric-key-prefix=<prefix>]
```
* If `function-name` is supplied, collect data from specified Lambda function.
  * If not, whole Lambda stastics in the region is collected.
* you can set some parameters by environment variables: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`.
  * If both of those environment variables and command line parameters are passed, command line parameters are used.
* You may omit `region` parameter if you're running this plugin on an EC2 instance running in same region with the target Lambda function

## Example of mackerel-agent.conf

```
[plugin.metrics.aws-lambda]
command = "/path/to/mackerel-plugin-aws-lambda -function-name=MyFunc -region=ap-northeast-1"
```
