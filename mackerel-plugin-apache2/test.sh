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

docker build -t test-$plugin testdata
docker run -d --name test-$plugin -p 10080:80 test-$plugin
trap 'docker stop test-$plugin; docker rm test-$plugin; docker rmi test-$plugin; exit' EXIT
sleep 10

if $plugin -p 10080 -status_page '/server-status?auto'
then
	echo OK
else
	echo FAIL
fi
