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

image=nginx:1.23
port=20080

docker run --name "test-$plugin" -p $port:80 -v $(pwd)/testdata/nginx.conf:/etc/nginx/nginx.conf:ro -d $image
trap 'docker stop test-$plugin; docker rm test-$plugin; exit 1' 1 2 3 15
sleep 10

# to store previous value to calculate a diff of metrics
$plugin -port $port >/dev/null 2>&1
sleep 1

$plugin -port $port | graphite-metric-test -f rule.txt
status=$?

docker stop "test-$plugin"
docker rm "test-$plugin"
exit $status
