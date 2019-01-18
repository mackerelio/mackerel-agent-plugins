mackerel-plugin-php-fpm
=====================

PHP-FPM status custom metrics plugin for [mackerel-agent](https://github.com/mackerelio/mackerel-agent).

## Synopsis

```shell
mackerel-plugin-php-fpm [-metric-key-prefix=php-fpm] [-timeout=5] [-url=http://localhost/status?json] [-socket unix:///var/run/php-fpm.sock
```

## Example of mackerel-agent.conf

```
[plugin.metrics.php-fpm]
command = "/path/to/mackerel-plugin-php-fpm"
```

## Author
[ariarijp](https://github.com/ariarijp)
