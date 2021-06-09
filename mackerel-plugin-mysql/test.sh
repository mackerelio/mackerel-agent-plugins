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

password=passpass
port=13306
docker run -d \
	--name test-$plugin \
	-p $port:3306 \
	-e MYSQL_ROOT_PASSWORD=$password \
	mysql:8 --default-authentication-plugin=mysql_native_password
trap 'docker stop test-$plugin; docker rm test-$plugin; exit' 1 2 3 15 EXIT
sleep 10

# wait until bootstrap mysqld..
while (( i++ < 3 )) && ! $plugin -port $port -password $password -enable_extended
do
	sleep 3
done
sleep 1
#export MACKEREL_PLUGIN_WORKDIR=tmp
$plugin -port $port -password $password -enable_extended
