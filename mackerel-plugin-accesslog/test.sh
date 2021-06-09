#!/bin/sh

prog=$(basename $0)
cd $(dirname $0)
PATH=$(pwd):$PATH
plugin=$(basename $(pwd))
if ! which $plugin >/dev/null
then
	echo "$prog: $plugin is not installed" >&2
	exit 2
fi

status=0
for i in lib/testdata/*.*
do
	if $plugin -no-posfile $i
	then
		echo OK: $i
	else
		status=$?
		echo FAIL: $i
	fi
done
exit $status
