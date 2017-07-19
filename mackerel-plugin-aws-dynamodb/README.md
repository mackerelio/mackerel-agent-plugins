mackerel-plugin-aws-dynamodb
=================================

AWS DynamoDB custom metrics plugin for mackerel.io agent.
Currently this plugin doesn't support following metrics:

- Metrics which take `GlobalSecondaryIndexName` Dimension
- Metrics related to DynamoDB Streams

## Synopsis

```shell
mackerel-plugin-aws-dynamodb -table-name=<table-name> -region=<aws-region> [-access-key-id=<id>] [-secret-access-key=<key>] [-metric-key-prefix=<key-prefix>]
```
* collect data from specified AWS DynamoDB
* you can set keys by environment variables: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`

## Example of mackerel-agent.conf

```toml
[plugin.metrics.aws-dynamodb]
command = "/path/to/mackerel-plugin-aws-dynamodb -table-name=MyTable -region=ap-northeast-1"
```
