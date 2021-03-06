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

user=postgres
password=passpass
docker run -d \
	--name test-$plugin \
	-p 15432:5432 \
	-e POSTGRES_PASSWORD=$password postgres:11
trap 'docker stop test-$plugin; docker rm test-$plugin; exit' EXIT
sleep 10

exec $plugin -port 15432 -user=$user -password $password
