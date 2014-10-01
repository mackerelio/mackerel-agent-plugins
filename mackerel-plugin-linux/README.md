mackerel-plugin-linux
============================

## Description

Get linux process metrics for Mackerel and Sensu.

For example ...

- Page swapped (swap)
- Connections each network service primitives (netstat)
- Disk read time (diskstats)
- Interrupts (proc_stat)
- Context switches
- Forks
- Login users (users)

## Required

- Linux Kernel 2.6.32 or above.

## Usage

### Build this program

First, build this program.

```
go get github.com/mackerelio/mackerel-agent-plugins
cd $GO_HOME/src/github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-linux
go test
go build
```

### Execute this plugin

Next, you can execute this program :-)

```
./mackerel-plugin-linux
```

### Add mackerel-agent.conf

Finally, if you want to get linux metrics via Mackerel, please edit mackerel-agent.conf. For example is below.

```
[plugin.metrics.linux]
command = "/path/to/mackerel-plugin-linux"
type = "metric"
```

## Optional: Selecting get metrics

## For more information

Please execute 'mackerel-plugin-linux -h' and you can get command line options.

## References

This program to use as reference from [Percona Monitoring Plugins for Cacti](http://www.percona.com/doc/percona-monitoring-plugins/).

## Author

[Yuichiro Saito](https://github.com/koemu)
