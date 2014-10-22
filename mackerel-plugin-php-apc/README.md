mackerel-plugin-php-apc
====

## Description

Get PHP APC (Alternative PHP Cache) metrics for Mackerel and Sensu.

## Usage (for Apache)

### Build this program

Next, build this program.

```
go get github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-php-apc
cd $GO_HOME/src/github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-php-apc
go test
go build
cp -a mackerel-plugin-php-apc /usr/local/bin/
mkdir DOCUMENT_ROOT/mackerel/
cp -a php-apc.php DOCUMENT_ROOT/mackerel/
```

### Set up your apache server

You should enable to execute PHP program (e.g. mod_php).

### Add apache config

Edit your apache config file to access metric from localhost only. For example is below.

```
<Directory "HTTP_HOME/mackerel/">
    Order deny,allow
    Deny from all
    Allow from 127.0.0.1 ::1
</Directory>
```
 
And, reload apache configuration.

```
sudo service httpd configtest
sudo service httpd reload
```

### Execute this plugin

And, you can execute this program :-)

```
./mackerel-plugin-php-apc
```

### Add mackerel-agent.conf

Finally, if you want to get php-apc metrics via Mackerel, please edit mackerel-agent.conf. For example is below.

```
[plugin.metrics.php-apc]
command = "/path/to/mackerel-plugin-php-apc"
type = "metric"
```

## For more information

Please execute 'mackerel-plugin-php-apc -h' and you can get command line options.

## Author

[Yuichiro Saito](https://github.com/koemu)
