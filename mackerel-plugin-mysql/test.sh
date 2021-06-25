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

password=passpass
port=13306
image=mysql:8
rule=rule_mysql8_extend.txt

docker run -d \
	--name "test-$plugin" \
	-p $port:3306 \
	-e MYSQL_ROOT_PASSWORD=$password \
	"$image" --default-authentication-plugin=mysql_native_password
trap 'docker stop test-$plugin; docker rm test-$plugin; exit' 1 2 3 15 EXIT
sleep 10

#export MACKEREL_PLUGIN_WORKDIR=tmp

# wait until bootstrap mysqld..
for i in $(seq 5)
do
	echo "Connecting $i..."
	if $plugin -port $port -password $password -enable_extended >/dev/null 2>&1
	then
		break
	fi
	sleep 3
done
sleep 1

# to store previous value to calculate a diff of metrics
$plugin -port $port -password $password -enable_extended >/dev/null 2>&1
sleep 1

$plugin -port $port -password $password -enable_extended | graphite-metric-test -f "$rule"
