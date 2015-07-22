mackerel-agent-plugins  [![Build Status](https://travis-ci.org/mackerelio/mackerel-agent-plugins.svg?branch=master)](https://travis-ci.org/mackerelio/mackerel-agent-plugins)
======================

This is the official plugin pack for [mackerel-agent](https://github.com/mackerelio/mackerel-agent), a piece of software which is installed on your hosts to collect metrics and events and send them to [Mackerel](https://mackerel.io).

Detailed specs of plugins can be viewed here: ENG [mackerel-agent specifications](http://help.mackerel.io/entry/spec/agent), JPN [mackerel-agent 仕様](http://help-ja.mackerel.io/entry/spec/agent).

Documentation for each plugin is located in its respective sub directory.

* [mackerel-plugin-apache2](./mackerel-plugin-apache2/README.md)
* [mackerel-plugin-aws-cloudfront](./mackerel-plugin-aws-cloudfront/README.md)
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
* [mackerel-plugin-td-table-count](./mackerel-plugin-td-table-count/README.md)
* [mackerel-plugin-varnish](./mackerel-plugin-varnish/README.md)
* [mackerel-plugin-xentop](./mackerel-plugin-xentop/README.md)

Installation
============

## Install mackerel-agent

ENG http://help.mackerel.io/entry/howto/install-agent
JPN http://help-ja.mackerel.io/entry/howto/install-agent

If the mackerel-agent has already be installed this step can be ignored.

## Install mackerel-agent-plugins

Install the plugin pack from either the yum or the apt repository.

### CentOS 5/6

```shell
yum install mackerel-agent-plugins
```

### Debian 6/7

```shell
apt-get install mackerel-agent-plugins
```

mackerel-agent-plugins will be installed to ```/usr/local/bin/mackerel-plugin-*```.

Caution
=======

Some plugins may not work on CentOS/RedHat 5 because the golang compiler (gc) doesn't support the old kernel.
(https://golang.org/doc/install)

Contribution
============

* fork it
* develop the plugin you want
* create a pullrequest!



