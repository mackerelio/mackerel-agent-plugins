mackerel-plugin-apache2
====

## Description

Get apache "server-status" metrics for Mackerel and Sensu.

## Usage

### Set up your apache server

First, you should enabled mod_status module, and configure apache config file. For example is below.

```
ExtendedStatus On
<VirtualHost 127.0.0.1:1080>
    <Location /server-status>
        SetHandler server-status
        Order deny,allow
        Deny from all
        Allow from localhost
    </Location>
</VirtualHost>
```

### Build this program

Next, build this program.

```
go get github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-apache2
cd $GO_HOME/src/github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-apache2
go test
go build
```

### Execute this plugin

And, you can execute this program :-)

```
./mackerel-plugin-apache2 -p 1080
```

### Add mackerel-agent.conf

Finally, if you want to get apache2 metrics via Mackerel, please edit mackerel-agent.conf. For example is below.

```
[plugin.metrics.apache2]
command = "/path/to/mackerel-plugin-apache2 -p 1080"
type = "metric"
```

## For more information

Please execute 'mackerel-plugin-apache2 -h' and you can get command line options.

## References

This program to use as reference from [Percona Monitoring Plugins for Cacti](http://www.percona.com/doc/percona-monitoring-plugins/).

## Author

[Yuichiro Saito](https://github.com/koemu)
