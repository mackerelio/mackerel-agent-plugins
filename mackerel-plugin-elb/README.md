mackerel-plugin-elb
=====================

AWS ELB custom metrics plugin for mackerel.io agent.
As it stands, this can fetch only across-all-LBs metrics.

## Synopsis

```shell
mackerel-plugin-elb [-region=<aws-region>] [-access-key-id=<id>] [-secret-access-key==<key>] [-tempfile=<tempfile>]
```
* if you run on an ec2-instance, you probably don't have to specify `-region`
* if you run on an ec2-instance and the instance is associated with an appropriate IAM Role, you probably don't have to specify `-access-key-id` & `-secret-access-key`

## Example of mackerel-agent.conf

```
[plugin.metrics.elb]
command = "/path/to/mackerel-plugin-elb"
```
