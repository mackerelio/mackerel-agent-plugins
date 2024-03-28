#!/bin/sh

prog=$(basename "$0")
if ! [ -S /var/run/docker.sock ]
then
	echo "$prog: there are no running docker" >&2
	exit 2
fi

cd "$(dirname "$0")" || exit
PATH=$(pwd):$PATH
plugin=$(basename "$(pwd)")
if ! which "$plugin" >/dev/null
then
	echo "$prog: $plugin is not installed" >&2
	exit 2
fi

image=local/test-$plugin
# mackerel-plugin-snmp is disallowed --port option.
port=161

docker build -t $image testdata/

docker run --name "test-$plugin" -v $(pwd)/testdata/snmpd.conf:/etc/snmp/snmpd.conf:ro -p $port:161/udp -d $image
trap 'docker stop test-$plugin; docker rm test-$plugin; exit 1' 1 2 3 15
sleep 10

$plugin '.1.3.6.1.2.1.25.1.5.0:hrSystemNumUsers:0:0' '.1.3.6.1.2.1.25.1.6.0:hrSystemProcesses:0:0' '.1.3.6.1.4.1.8072.1.3.2.3.1.2.4.101.99.104.111:echo:0:0' | graphite-metric-test -f rule.txt
status=$?

docker stop "test-$plugin"
docker rm "test-$plugin"
exit $status
