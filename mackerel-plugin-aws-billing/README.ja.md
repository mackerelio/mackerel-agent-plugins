mackerel-plugin-aws-billing
=======================

# Overview
このmackerelプラグインはAWSの料金をMackerel上にグラフ化することができます。

## Description
このプラグインは[AWSのCloudWatch Api](https://aws.amazon.com/ja/documentation/cloudwatch/)を使用してAWSアカウントに紐づくコスト情報を取得します。  
mackerelのホストに紐付かないグラフを表示する[サービスメトリック](https://mackerel.io/ja/features/service-metrics/)かホストに紐付けて表示することが可能です。
AWSの料金情報は他のサーバーのメトリックなどとは違い、短時間で変わるものではなく、数時間に一度更新されるものです。  

なので、このプラグインがAWSのAPIを使用して料金を取得するのは1時間に1度だけです。  

1時間に1度、ローカルのキャッシュファイルに書き込みを行い、そのキャッシュファイルを使用してデータを出力します。  

(※サービスメトリックにグラフを出力する場合は、mackerel-agentがデータを送信するのではなく、mackerel-agentにより動かされたプログラムがmackerelにデータを送信します。
  これはmackerel-agentがサービスメトリックにデータを送信することができないためです。
)

注意: 必ずAWSの請求情報をcloudwatchから取得できるよう、請求アラートを有効にしてください。(https://docs.aws.amazon.com/ja_jp/awsaccountbilling/latest/aboutv2/monitor-charges.html) 

## Usage

ホストに対してデータを送信する場合

```shell
mackerel-plugin-aws-billing [-access-key-id=<id>] [-secret-access-key=<key>] [-target=<aws-services>] [-currency=<currency>]
```

サービスメトリックにデータを送信する場合

```shell
mackerel-plugin-aws-billing [-access-key-id=<id>] [-secret-access-key=<key>] [-target=<aws-services>] [-currency=<currency>] [-mode=<SerivceMetric>] [-api-key=<api-key>] [-servicename=<servicename>]
```

- access-key-id(require *1)
  AWSから発行されるaccess key id です。AWSのAPIを利用するために必要です。

- secret-access-key(require *1)
  AWSから発行されるsecret access key です。AWSのAPIを利用するために必要です。

- target(optional)
  料金情報を取得する必要のあるAWSのサービスを指定します。  
  ,区切りで指定します。(例: AmazonEC2,AWSLambda)
  合計金額が知りたい時はAllを指定してください。

  料金情報を取得できるAWSのサービスはアカウントによって異なります。

  このパラメータを省略すると、取得できる全てのサービスの料金情報を出力します。

- currency(optional)
  デフォルトはUSドルです。

- dest(optional)
  現状はServiceMetricのみ指定できます。  
  この値を渡すとサービスメトリックにデータを出力します。  

- api-key(require *2)
  modeにServiceMetricが指定された時のみ必須です。  
  mackerelから発行されるAPIのKeyです。  
  このAPI KeyにはRead, Writeの権限が必要です。  

- servicename(require *2)
  modeにServiceMetricが指定された時のみ必須です。  

  グラフを出力するmackerel上のサービスの名前を指定してください。 

*1 環境変数に値を入れることも可能です。(https://github.com/aws/aws-sdk-go#configuring-credentials)  

*2 modeにServiceMetricが指定された時のみ必須です。  

### Example of mackerel-agent.conf
```
[plugin.metrics.aws-billing]
command = "/path/to/mackerel-plugin-aws-billing/main" ... //arguments
```
