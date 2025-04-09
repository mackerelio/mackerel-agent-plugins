#!/bin/sh

set -eu
cd build/mackerel-agent-plugins-${VERSION}-${ARCH}
for i in accesslog apache2 aws-cloudfront aws-dynamodb aws-ec2-cpucredit aws-ec2-ebs aws-elasticache aws-elasticsearch aws-elb aws-kinesis-streams aws-lambda aws-rds aws-s3-requests aws-ses conntrack docker elasticsearch fluentd gostats h2o haproxy inode jmx-jolokia jvm linux mailq memcached mongodb multicore munin mysql nginx openldap php-apc php-fpm php-opcache plack postgres proc-fd rabbitmq redis sidekiq snmp solr squid td-table-count trafficserver twemproxy unicorn uptime uwsgi-vassal varnish; do \
    ln -s ./mackerel-plugin mackerel-plugin-$i; \
done
cd ..
tar czf mackerel-agent-plugins-${VERSION}-${ARCH}.tar.gz mackerel-agent-plugins-${VERSION}-${ARCH}
