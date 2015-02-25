mackerel-agent-plugins  [![Build Status](https://travis-ci.org/mackerelio/mackerel-agent-plugins.svg?branch=master)](https://travis-ci.org/mackerelio/mackerel-agent-plugins)
======================

Plugins for [mackerel-agent](https://github.com/mackerelio/mackerel-agent), which is resource aggregator for [Mackerel](https://mackerel.io).

Detailed specification of plugin can be shown at [mackerel-agent specification](http://help-ja.mackerel.io/entry/spec/agent).

Document of each plugin is located under each sub directory.

* [mackerel-plugin-apache2](./mackerel-plugin-apache2/README.md)
* [mackerel-plugin-aws-ec2-cpucredit](./mackerel-plugin-aws-ec2-cpucredit/README.md)
* [mackerel-plugin-aws-elasticache](./mackerel-plugin-aws-elasticache/README.md)
* [mackerel-plugin-aws-elb](./mackerel-plugin-aws-elb/README.md)
* [mackerel-plugin-aws-rds](./mackerel-plugin-aws-rds/README.md)
* [mackerel-plugin-aws-ses](./mackerel-plugin-aws-ses/README.md)
* [mackerel-plugin-elasticsearch](./mackerel-plugin-elasticsearch/README.md)
* [mackerel-plugin-haproxy](./mackerel-plugin-haproxy/README.md)
* [mackerel-plugin-jvm](./mackerel-plugin-jvm/README.md)
* [mackerel-plugin-linux](./mackerel-plugin-linux/README.md)
* [mackerel-plugin-memcached](./mackerel-plugin-memcached/README.md)
* [mackerel-plugin-mongodb](./mackerel-plugin-mongodb/README.md)
* [mackerel-plugin-munin](./mackerel-plugin-munin/README.md)
* [mackerel-plugin-mysql](./mackerel-plugin-mysql/README.md)
* [mackerel-plugin-nginx](./mackerel-plugin-nginx/README.md)
* [mackerel-plugin-php-apc](./mackerel-plugin-php-apc/README.md)
* [mackerel-plugin-php-opcache](./mackerel-plugin-php-opcache/README.md)
* [mackerel-plugin-plack](./mackerel-plugin-plack/README.md)
* [mackerel-plugin-postgres](./mackerel-plugin-postgres/README.md)
* [mackerel-plugin-redis](./mackerel-plugin-redis/README.md)
* [mackerel-plugin-snmp](./mackerel-plugin-snmp/README.md)
* [mackerel-plugin-squid](./mackerel-plugin-squid/README.md)
* [mackerel-plugin-varnish](./mackerel-plugin-varnish/README.md)

Installation
============

## Install mackerel-agent

http://help-ja.mackerel.io/entry/howto/install-agent

Skip this process when you have already installed mackerel-agent.

## Install mackerel-agent-plugins

Install plugins via yum or apt repository.

### CentOS 5/6

```shell
yum install mackerel-agent-plugins
```

### Debian 6/7

```shell
apt-get install mackerel-agent-plugins
```

mackerel-agent-plugins are installed to ```/usr/local/bin/mackerel-plugin-*```.

Caution
=======

Some plugins may not work on CentOS/RedHat 5 because golang compiler (gc) doesn't support old kernel.
(https://golang.org/doc/install)

Contribution
============

* fork it
* develop plugin you want
* create pullrequest!



