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

image=fluent/fluentd:v1.13-1
port0=24230
port1=24231
port2=24232
rule="rule-workers.txt"

docker run -d \
	--name "test-$plugin" \
	-p "$port0:$port0" \
	-p "$port1:$port1" \
	-p "$port2:$port2" \
	-v "$(pwd)/testdata:/fluentd/etc:ro" "$image" -c /fluentd/etc/fluentd-workers.conf
trap 'docker stop test-$plugin; docker rm test-$plugin; exit' 1 2 3 15
sleep 10

$plugin -port "$port0" -workers 3 | graphite-metric-test -f "$rule"
status=$?
docker stop "test-$plugin"
docker rm "test-$plugin"
exit $status
