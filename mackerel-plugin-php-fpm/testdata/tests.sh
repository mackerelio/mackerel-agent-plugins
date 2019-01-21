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
cleanup()
{
	local status=$?
	rm -rf $d
	pkill php-fpm
	exit $status
}
trap cleanup HUP INT QUIT TERM EXIT

#
# tests
#
sock_tcp=tcp://localhost:9000
sock_unix=unix:///tmp/php-fpm.sock
protocols=(tcp unix)
for proto in "${protocols[@]}"
do
	pidfile=$d/php-fpm.$proto.pid
	mkfifo $pidfile
	php-fpm -y php-fpm.$proto.conf -g $pidfile -D
	pid=$(cat $pidfile)
	sock=sock_$proto
	run_plugin -socket "${!sock}" -url 'http://localhost:9000/status?json'
	kill $pid
	while ps -p $pid >/dev/null
	do
		sleep 1
	done
done
