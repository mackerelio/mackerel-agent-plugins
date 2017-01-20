mackerel-plugin-aws-billing
=======================

# Overview
mackerel-plugin-aws-billing is mackerel agent plugin that gets AWS cost and makes a graph.

## Description

It uses [AWS CloudWatch Api](https://aws.amazon.com/ja/documentation/cloudwatch/) to get AWS Billing Data related to AWS Account.  
It can make a graph on [Service Metric](https://mackerel.io/ja/features/service-metrics/).  
It gets AWS cost every hour by using AWS API because unlike other metrics, AWS updates cost once in a few hours.  
It writes cache file on server every hour, and uses cache to output data.  

(※In order to output data to Service Metric, this plugin sends data to mackerel instead of mackerel-agent because mackerel-agent can not send data to Service Metric.)  

Caution: You must must enable Billing Alerts.(see https://docs.aws.amazon.com/ja_jp/awsaccountbilling/latest/aboutv2/monitor-charges.html)  

## Usage

Send data to Host.  

```shell
mackerel-plugin-aws-billing [-access-key-id=<id>] [-secret-access-key=<key>] [-target=<aws-services>] [-currency=<currency>]
```

Send data to Service Metric.  

```shell
mackerel-plugin-aws-billing [-access-key-id=<id>] [-secret-access-key=<key>] [-target=<aws-services>] [-currency=<currency>] [-dest=<SerivceMetric or Host>] [-api-key=<api-key>] [-service-name=<servicename>]
```

- access-key-id(required)*1  
  access-key-id is published by AWS. It is required to use AWS API. 

- secret-access-key(required)*1  
  secret-access-key is published by AWS. It is required to use AWS API.. 

- target(optional)  
  If target is not specified, this plugin gets all available services.
  If target is specified, this plugin gets specified available services.
  If you want to specify multiple services, separate with comma.(ex: target=AmazonEC2,AWSLambda)
  If you want to get sum of costs, give All to target.(ex: target=All) 

  You can get cost of service that you use.

- currency(optional)  
  Defalut is in US Dollar.

- dest(required)  
  Use ServiceMetric or Host.

  If set ServiceMetric, This plugin outputs data to ServiceMetric.
  If set Host, This plugin outputs data to Host.

- api-key(required)*2  
  Required only if ServiceMetric is given to mode parameter.
  API Key is published by mackerel. 
  You must give read and write premissions to API Key.

- service-name(required)*2  
  Required only if ServiceMetric is given to mode parameter.
  Specify mackerel service name to make graph on Service Metric page.

*1 You can set keys by environment variables: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY (see https://github.com/aws/aws-sdk-go#configuring-credentials)  
*2 Required only if ServiceMetric is given to mode parameter.  

### Example of mackerel-agent.conf
```
[plugin.metrics.aws-billing]
command = "/path/to/mackerel-plugin-aws-billing/main" ... //arguments
```
