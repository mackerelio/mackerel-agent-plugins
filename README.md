mackerel-agent-plugins  [![Build Status](https://travis-ci.org/mackerelio/mackerel-agent-plugins.svg?branch=master)](https://travis-ci.org/mackerelio/mackerel-agent-plugins)
======================

This is the official plugin pack for [mackerel-agent](https://github.com/mackerelio/mackerel-agent), a piece of software which is installed on your hosts to collect metrics and events and send them to [Mackerel](https://mackerel.io).

Detailed specs of plugins can be viewed here: ENG [mackerel-agent specifications](https://mackerel.io/docs/entry/spec/agent), JPN [mackerel-agent 仕様](https://mackerel.io/ja/docs/entry/spec/agent).

Documentation for each plugin is located in its respective sub directory.

* [mackerel-plugin-apache2](./mackerel-plugin-apache2/README.md)
* [mackerel-plugin-aws-cloudfront](./mackerel-plugin-aws-cloudfront/README.md)
* [mackerel-plugin-aws-ec2](./mackerel-plugin-aws-ec2/README.md)
* [mackerel-plugin-aws-ec2-cpucredit](./mackerel-plugin-aws-ec2-cpucredit/README.md)
* [mackerel-plugin-aws-ec2-ebs](./mackerel-plugin-aws-ec2-ebs/README.md)
* [mackerel-plugin-aws-elasticache](./mackerel-plugin-aws-elasticache/README.md)
* [mackerel-plugin-aws-elasticsearch](./mackerel-plugin-aws-elasticsearch/README.md)
* [mackerel-plugin-aws-elb](./mackerel-plugin-aws-elb/README.md)
* [mackerel-plugin-aws-rds](./mackerel-plugin-aws-rds/README.md)
* [mackerel-plugin-aws-ses](./mackerel-plugin-aws-ses/README.md)
* [mackerel-plugin-conntrack](./mackerel-plugin-conntrack/README.md)
* [mackerel-plugin-docker](./mackerel-plugin-docker/README.md)
* [mackerel-plugin-elasticsearch](./mackerel-plugin-elasticsearch/README.md)
* [mackerel-plugin-fluentd](./mackerel-plugin-fluentd/README.md)
* [mackerel-plugin-gearmand](./mackerel-plugin-gearmand/README.md)
* [mackerel-plugin-gostats](./mackerel-plugin-gostats/README.md)
* [mackerel-plugin-graphite](./mackerel-plugin-graphite/README.md)
* [mackerel-plugin-haproxy](./mackerel-plugin-haproxy/README.md)
* [mackerel-plugin-inode](./mackerel-plugin-inode/README.md)
* [mackerel-plugin-jmx-jolokia](./mackerel-plugin-jmx-jolokia/README.md)
* [mackerel-plugin-jvm](./mackerel-plugin-jvm/README.md)
* [mackerel-plugin-linux](./mackerel-plugin-linux/README.md)
* [mackerel-plugin-mailq](./mackerel-plugin-mailq/README.md)
* [mackerel-plugin-memcached](./mackerel-plugin-memcached/README.md)
* [mackerel-plugin-mongodb](./mackerel-plugin-mongodb/README.md)
* [mackerel-plugin-munin](./mackerel-plugin-munin/README.md)
* [mackerel-plugin-murmur](./mackerel-plugin-murmur/README.md)
* [mackerel-plugin-mysql](./mackerel-plugin-mysql/README.md)
* [mackerel-plugin-nginx](./mackerel-plugin-nginx/README.md)
* [mackerel-plugin-php-apc](./mackerel-plugin-php-apc/README.md)
* [mackerel-plugin-php-fpm](./mackerel-plugin-php-fpm/README.md)
* [mackerel-plugin-php-opcache](./mackerel-plugin-php-opcache/README.md)
* [mackerel-plugin-plack](./mackerel-plugin-plack/README.md)
* [mackerel-plugin-postgres](./mackerel-plugin-postgres/README.md)
* [mackerel-plugin-proc-fd](./mackerel-plugin-proc-fd/README.md)
* [mackerel-plugin-rabbitmq](./mackerel-plugin-rabbitmq/README.md)
* [mackerel-plugin-redis](./mackerel-plugin-redis/README.md)
* [mackerel-plugin-snmp](./mackerel-plugin-snmp/README.md)
* [mackerel-plugin-solr](./mackerel-plugin-solr/README.md)
* [mackerel-plugin-squid](./mackerel-plugin-squid/README.md)
* [mackerel-plugin-td-table-count](./mackerel-plugin-td-table-count/README.md)
* [mackerel-plugin-trafficserver](./mackerel-plugin-trafficserver/README.md)
* [mackerel-plugin-twemproxy](./mackerel-plugin-twemproxy/README.md)
* [mackerel-plugin-unicorn](./mackerel-plugin-unicorn/README.md)
* [mackerel-plugin-uptime](./mackerel-plugin-uptime/README.md)
* [mackerel-plugin-varnish](./mackerel-plugin-varnish/README.md)
* [mackerel-plugin-windows-server-sessions](./mackerel-plugin-windows-server-sessions/README.md)
* [mackerel-plugin-xentop](./mackerel-plugin-xentop/README.md)

Installation
============

## Install mackerel-agent

ENG https://mackerel.io/docs/entry/howto/install-agent
JPN https://mackerel.io/ja/docs/entry/howto/install-agent

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

mackerel-agent-plugins will be installed to ```/usr/bin/mackerel-plugin-*```.

Caution
=======

Some plugins may not work on CentOS/RedHat 5 because the golang compiler (gc) doesn't support the old kernel.
(https://golang.org/doc/install)

Some plugins are not contained in rpm and deb packages. If you want to use them, build them.

Contribution
============

* fork it
* develop the plugin you want
* create a pullrequest!
