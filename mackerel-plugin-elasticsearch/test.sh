#!/bin/sh

prog=$(basename $0)
if ! [ -S /var/run/docker.sock ]
then
	echo "$prog: there are no running docker" >&2
	exit 2
fi

cd $(dirname $0)
PATH=$(pwd):$PATH
plugin=$(basename $(pwd))
if ! which $plugin >/dev/null
then
	echo "$prog: $plugin is not installed" >&2
	exit 2
fi

password=passpass
port=9200
docker run -d \
	--name test-$plugin \
	-p $port:$port \
	-e "ELASTIC_PASSWORD=$password" \
	-e "discovery.type=single-node" \
	-e "ingest.geoip.downloader.enabled=false" \
	elasticsearch:8.5.0
trap 'docker stop test-$plugin; docker rm test-$plugin; exit 1' 1 2 3 15

# to store previous value to calculate a diff of metrics
../tool/waituntil.bash -n 300 $plugin -scheme https -port $port -user=elastic -password $password -insecure -suppress-missing-error >/dev/null 2>&1

sleep 1

$plugin -scheme https -port $port -user=elastic -password $password -insecure -suppress-missing-error | graphite-metric-test -f rule.txt
status=$?

docker stop "test-$plugin"
docker rm "test-$plugin"
exit $status
