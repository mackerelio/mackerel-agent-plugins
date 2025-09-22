mackerel-agent-plugins  [![Build Status](https://github.com/mackerelio/mackerel-agent-plugins/workflows/test/badge.svg)](https://github.com/mackerelio/mackerel-agent-plugins/actions?workflow=test)
======================

This is the official plugin pack for [mackerel-agent](https://github.com/mackerelio/mackerel-agent), a piece of software which is installed on your hosts to collect metrics and events and send them to [Mackerel](https://mackerel.io).

Detailed specs of plugins can be viewed here: ENG [mackerel-agent specifications](https://mackerel.io/docs/entry/spec/agent), JPN [mackerel-agent 仕様](https://mackerel.io/ja/docs/entry/spec/agent).

Documentation for each plugin is located in its respective sub directory.

* [mackerel-plugin-accesslog](./mackerel-plugin-accesslog/README.md)
* [mackerel-plugin-apache2](./mackerel-plugin-apache2/README.md)
* [mackerel-plugin-aws-ec2-ebs](./mackerel-plugin-aws-ec2-ebs/README.md)
* [mackerel-plugin-conntrack](./mackerel-plugin-conntrack/README.md)
* [mackerel-plugin-docker](./mackerel-plugin-docker/README.md)
* [mackerel-plugin-elasticsearch](./mackerel-plugin-elasticsearch/README.md)
* [mackerel-plugin-fluentd](./mackerel-plugin-fluentd/README.md)
* [mackerel-plugin-gostats](./mackerel-plugin-gostats/README.md)
* [mackerel-plugin-h2o](./mackerel-plugin-h2o/README.md)
* [mackerel-plugin-haproxy](./mackerel-plugin-haproxy/README.md)
* [mackerel-plugin-inode](./mackerel-plugin-inode/README.md)
* [mackerel-plugin-jmx-jolokia](./mackerel-plugin-jmx-jolokia/README.md)
* [mackerel-plugin-jvm](./mackerel-plugin-jvm/README.md)
* [mackerel-plugin-linux](./mackerel-plugin-linux/README.md)
* [mackerel-plugin-mailq](./mackerel-plugin-mailq/README.md)
* [mackerel-plugin-memcached](./mackerel-plugin-memcached/README.md)
* [mackerel-plugin-mongodb](https://github.com/mackerelio/mackerel-plugin-mongodb/blob/main/README.md)
* [mackerel-plugin-mssql](./mackerel-plugin-mssql/README.md)
* [mackerel-plugin-multicore](./mackerel-plugin-multicore/README.md)
* [mackerel-plugin-munin](./mackerel-plugin-munin/README.md)
* [mackerel-plugin-mysql](https://github.com/mackerelio/mackerel-plugin-mysql/blob/main/README.md)
* [mackerel-plugin-nginx](./mackerel-plugin-nginx/README.md)
* [mackerel-plugin-openldap](./mackerel-plugin-openldap/README.md)
* [mackerel-plugin-php-apc](./mackerel-plugin-php-apc/README.md)
* [mackerel-plugin-php-fpm](./mackerel-plugin-php-fpm/README.md)
* [mackerel-plugin-php-opcache](./mackerel-plugin-php-opcache/README.md)
* [mackerel-plugin-plack](./mackerel-plugin-plack/README.md)
* [mackerel-plugin-postgres](./mackerel-plugin-postgres/README.md)
* [mackerel-plugin-proc-fd](./mackerel-plugin-proc-fd/README.md)
* [mackerel-plugin-rabbitmq](./mackerel-plugin-rabbitmq/README.md)
* [mackerel-plugin-redis](./mackerel-plugin-redis/README.md)
* [mackerel-plugin-sidekiq](./mackerel-plugin-sidekiq/README.md)
* [mackerel-plugin-snmp](./mackerel-plugin-snmp/README.md)
* [mackerel-plugin-solr](./mackerel-plugin-solr/README.md)
* [mackerel-plugin-squid](./mackerel-plugin-squid/README.md)
* [mackerel-plugin-td-table-count](./mackerel-plugin-td-table-count/README.md)
* [mackerel-plugin-trafficserver](./mackerel-plugin-trafficserver/README.md)
* [mackerel-plugin-twemproxy](./mackerel-plugin-twemproxy/README.md)
* [mackerel-plugin-unicorn](./mackerel-plugin-unicorn/README.md)
* [mackerel-plugin-uptime](./mackerel-plugin-uptime/README.md)
* [mackerel-plugin-uwsgi-vassal](./mackerel-plugin-uwsgi-vassal/README.md)
* [mackerel-plugin-varnish](./mackerel-plugin-varnish/README.md)
* [mackerel-plugin-windows-process-stats](./mackerel-plugin-windows-process-stats/README.md)
* [mackerel-plugin-windows-server-sessions](./mackerel-plugin-windows-server-sessions/README.md)

Installation
============

## Install mackerel-agent

ENG https://mackerel.io/docs/entry/howto/install-agent
JPN https://mackerel.io/ja/docs/entry/howto/install-agent

If the mackerel-agent has already be installed this step can be ignored.

## Install mackerel-agent-plugins

Install the plugin pack from either the yum or the apt repository.
To setup these package repositories, see the documentation regarding the installation of mackerel-agent ([rpm](https://mackerel.io/docs/entry/howto/install-agent/rpm) / [deb](https://mackerel.io/docs/entry/howto/install-agent/deb)).

mackerel-agent-plugins will be installed to ```/usr/bin/mackerel-plugin-*``` (and some plugins are symlinked as ```/usr/local/bin/mackerel-plugin-*```, for backward compatibility.).

### yum

```shell
yum install mackerel-agent-plugins
```

### apt

```shell
apt-get install mackerel-agent-plugins
```

### Go

**mackerel-agent-plugins** supports two newer versions of Go.

```shell
go install github.com/mackerelio/mackerel-agent-plugins/...@latest
```

Caution
=======

Some plugins may not work on CentOS/RedHat 5 because the golang compiler (gc) doesn't support the old kernel.
(https://golang.org/doc/install)

Some plugins are not contained in rpm and deb packages. If you want to use them, build them.

Contribution
============

* see [The official plugin registry](https://mackerel.io/blog/entry/feature/20171116#The-official-plugin-registry) and [Pull requests to the existing central repository](https://mackerel.io/blog/entry/feature/20171116#Pull-Requests-to-the-existing-central-repository)
* fork it
* develop the plugin you want
* create a pullrequest!

License
=======
```
Copyright 2014 Hatena Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
