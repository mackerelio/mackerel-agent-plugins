#!/bin/sh

set -eu
cd build/mackerel-agent-plugins-${VERSION}-${ARCH}
for i in accesslog apache2 aws-ec2-ebs conntrack docker elasticsearch fluentd gostats h2o haproxy inode jmx-jolokia jvm linux mailq memcached mongodb multicore munin mysql nginx openldap php-apc php-fpm php-opcache plack postgres proc-fd rabbitmq redis sidekiq snmp solr squid td-table-count trafficserver twemproxy unicorn uptime uwsgi-vassal varnish; do \
    ln -s ./mackerel-plugin mackerel-plugin-$i; \
done
cd ..
tar czf mackerel-agent-plugins-${VERSION}-${ARCH}.tar.gz mackerel-agent-plugins-${VERSION}-${ARCH}
