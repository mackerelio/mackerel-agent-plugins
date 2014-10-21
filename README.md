mackerel-agent-plugins  [![Build Status](https://travis-ci.org/mackerelio/mackerel-agent-plugins.svg?branch=master)](https://travis-ci.org/mackerelio/mackerel-agent-plugins)
======================

Plugins for [mackerel-agent](https://github.com/mackerelio/mackerel-agent), which is resource aggregator for [Mackerel](https://mackerel.io).

Detailed specification of plugin can be shown at [mackerel-agent specification](http://help-ja.mackerel.io/entry/spec/agent).

Document of each plugin is located under each sub directory.

* [mackerel-plugin-apache2](./mackerel-plugin-apache2/README.md)
* [mackerel-plugin-memcached](./mackerel-plugin-memcached/README.md)
* [mackerel-plugin-mysql](./mackerel-plugin-mysql/README.md)
* [mackerel-plugin-nginx](./mackerel-plugin-nginx/README.md)
* [mackerel-plugin-plack](./mackerel-plugin-plack/README.md)
* [mackerel-plugin-postgres](./mackerel-plugin-postgres/README.md)
* [mackerel-plugin-redis](./mackerel-plugin-redis/README.md)

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

Contribution
============

* fork it
* develop plugin you want
* create pullrequest!



