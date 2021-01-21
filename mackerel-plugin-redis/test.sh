#!/bin/sh

prog=$(basename $0)
if ! [[ -S /var/run/docker.sock ]]
then
	echo "$prog: there are no running docker" >&2
	exit 2
fi

cd $(dirname $0)
PATH=$(pwd):$PATH
plugin=$(basename $(pwd))
if ! which -s $plugin
then
	echo "$prog: $plugin is not installed" >&2
	exit 2
fi

export REDIS_PASSWORD=passpass
docker run --name test-$plugin -p 16379:6379 -d redis:6 --requirepass $REDIS_PASSWORD
trap 'docker stop test-$plugin; docker rm test-$plugin; exit' EXIT
sleep 10

if $plugin -port 16379
then
	echo OK
else
	echo FAIL
fi
