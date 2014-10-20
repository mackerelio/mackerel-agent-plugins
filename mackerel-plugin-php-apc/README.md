mackerel-plugin-php-apc
====

## Description

Get PHP APC (Alternative PHP Cache) metrics for Mackerel and Sensu.

## Usage

### Set up your apache server

First, you should enabled mod_status module, and configure apache config file. For example is below.

```
<VirtualHost 127.0.0.1:1080>
    <Location /server-status>
        SetHandler server-status
    </Location>
</VirtualHost>
```

### Build this program

Next, build this program.

```
go get github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-php-apc
cd $GO_HOME/src/github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-php-apc
go test
go build
```

### Execute this plugin

And, you can execute this program :-)

```
./mackerel-plugin-php-apc -p 1080
```

### Add mackerel-agent.conf

Finally, if you want to get php-apc metrics via Mackerel, please edit mackerel-agent.conf. For example is below.

```
[plugin.metrics.php-apc]
command = "/path/to/mackerel-plugin-php-apc -p 1080"
type = "metric"
```

## For more information

Please execute 'mackerel-plugin-php-apc -h' and you can get command line options.

## Author

[Yuichiro Saito](https://github.com/koemu)
