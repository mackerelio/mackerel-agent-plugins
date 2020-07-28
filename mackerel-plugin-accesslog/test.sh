#!/bin/sh

prog=$(basename $0)
cd $(dirname $0)
plugin=$(basename $(pwd))
if ! which -s $plugin
then
	echo "$prog: $plugin is not installed" >&2
	exit 2
fi

for i in lib/testdata/*.*
do
	if $plugin -no-posfile $i >/dev/null 2>&1; then
		echo OK: $i
	else
		echo FAIL: $i
	fi
done
