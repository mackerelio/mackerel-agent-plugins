mackerel-plugin-aws-waf
=======================

AWS WAF custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-aws-waf -web-acl-id=<aws-waf-web-acl-id> [-access-key-id=<id>] [-secret-access-key=<key>] [-tempfile=<tempfile>]
```

## AWS IAM Policy
the credential provided manually or fetched automatically by IAM Role should have the policy that includes an action, 'cloudwatch:GetMetricStatistics'

## Example of mackerel-agent.conf

```
[plugin.metrics.aws-waf]
command = "/path/to/mackerel-plugin-aws-waf -web-acl-id=your-web-acl-id"
```

## Notes

This plugin only supports AWS WAF for CloudFront, and not the metrics of WAF for ALB.
