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

exec $plugin
