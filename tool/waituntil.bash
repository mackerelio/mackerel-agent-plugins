#!/bin/bash
# usage: waituntil [-n count] command [arg ...]

usage()
{
	echo "usage: $(basename $0) [-n count] command [arg ...]" >&2
	exit 2
}

count=0 # unlimited
while getopts :n: OPT
do
	case "$OPT" in
	:)	usage ;;
	n)	count="$OPTARG" ;;
	\?)	usage ;;
	esac
done
shift $((OPTIND - 1))
if (($# == 0))
then
	usage
fi

i=0
while (($count == 0 || $i < $count))
do
	if command "$@"
	then
		exit 0
	fi
	sleep 1
	((i++))
done
exit 1
