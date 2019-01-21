#!/usr/bin/env bash

if ! which -s php-fpm
then
	echo "$(basename $0): php-fpm not found" >&2
	exit 2
fi

run_plugin()
{
	sleep 1
	go run ../main.go "$@"
	sleep 1
}

d=$(mktemp -d -t mackerel-plugin-php-fpm)
trap 'rm -rf $d; exit 1' 1 2 3 15

#
# tests
#
protocols=(tcp unix)
for proto in "${protocols[@]}"
do
	pidfile=$d/php-fpm.$proto.pid
	mkfifo $pidfile
	php-fpm -y php-fpm.$proto.conf -g $pidfile -D
	pid=$(cat $pidfile)
	run_plugin -fcgi -url 'http://localhost:9000/status?json'
	kill $pid
	while ps -p $pid >/dev/null
	do
		sleep 1
	done
done

rm -rf $d
