mackerel-plugin-aws-s3-requests
=================================

AWS S3 requests metrics plugin for mackerel.io agent.

## Requirement

You need a metrics configuration FilterID for the target S3 bucket to get CloudWatch metrics. If you haven't created any configurations yet, you can create by AWS CLI or AWS API. (ref: https://docs.aws.amazon.com/AmazonS3/latest/dev/metrics-configurations.html)

## Synopsis

```shell
mackerel-plugin-aws-s3-requests -bucket-name=<bucket-name> -filter-id=<filter-id> -region=<aws-region> -access-key-id=<id> -secret-access-key=<key> [-tempfile=<tempfile>] [-metric-key-prefix=<prefix>] [-metric-label-prefix=<prefix>]
```
* `filter-id` is the id for the metrics configuration, described in "Requirement" section.
* you can set some parameters by environment variables: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`.
  * If both of those environment variables and command line parameters are passed, command line parameters are used.
* You may omit `region` parameter if you're running this plugin on an EC2 instance running in same region with the target bucket.

## Example of mackerel-agent.conf

```
[plugin.metrics.aws-s3-requests]
command = "mackerel-plugin-aws-s3-requests -bucket-name=MyBucket -filter-id=SomeFilterId -region=ap-northeast-1"
```

## List of Metrics

You can customize graph labels and metric names by command line parameters `-metric-label-prefix` and `-metric-key-prefix`.

| Graph Label | Metric Label | Metric (Full) Name | CloudWatch Name | CloudWatch Stastics |
|-----------------------------------|--------------|----------------------------------------------------------------|---------------------|:-------------------:|
| S3 Requests | All | `custom.s3-requests.requests.AllRequests` | AllRequests | Sum |
| S3 Requests | Get | `custom.s3-requests.requests.GetRequests` | GetRequests | Sum |
| S3 Requests | Put | `custom.s3-requests.requests.PutRequests` | PutRequests | Sum |
| S3 Requests | Delete | `custom.s3-requests.requests.DeleteRequests` | DeleteRequests | Sum |
| S3 Requests | Head | `custom.s3-requests.requests.HeadRequests` | HeadRequests | Sum |
| S3 Requests | Post | `custom.s3-requests.requests.PostRequests` | PostRequests | Sum |
| S3 Requests | Listt | `custom.s3-requests.requests.ListRequests` | ListRequests | Sum |
| S3 Errors | 4xx | `custom.s3-requests.errors.4xxErrors` | 4xxErrors | Sum |
| S3 Errors | 5xx | `custom.s3-requests.errors.5xxErrors` | 5xxErrors | Sum |
| S3 Bytes | Downloaded | `custom.s3-requests.bytes.BytesDownloaded` | BytesDownloaded | Sum |
| S3 Bytes | Uploaded | `custom.s3-requests.bytes.BytesUploaded` | BytesUploaded | Sum |
| S3 FirstByteLatency [ms] | Average | `custom.s3-requests.first_byte_latency.FirstByteLatencyAvg` | FirstBytesLatency | Average |
| S3 FirstByteLatency [ms] | Maximum | `custom.s3-requests.first_byte_latency.FirstByteLatencyMax` | FirstBytesLatency | Maximum |
| S3 FirstByteLatency [ms] | Minimum | `custom.s3-requests.first_byte_latency.FirstByteLatencyMin` | FirstBytesLatency | Minimum |
| S3 TotalRequestLatency [ms] | Average | `custom.s3-requests.total_request_latency.TotalRequestLatencyAvg` | TotalRequestLatency | Average |
| S3 TotalRequestLatency [ms] | Maximum | `custom.s3-requests.total_request_latency.TotalRequestLatencyMax` | TotalRequestLatency | Maximum |
| S3 TotalRequestLatency [ms] | Minimum | `custom.s3-requests.total_request_latency.TotalRequestLatencyMin` | TotalRequestLatency | Minimum |
