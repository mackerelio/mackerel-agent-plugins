mackerel-plugin-php-fpm
=====================

PHP-FPM status custom metrics plugin for [mackerel-agent](https://github.com/mackerelio/mackerel-agent).

## Synopsis

```shell
mackerel-plugin-php-fpm [-metric-key-prefix=php-fpm] [-timeout=5] [-url=http://localhost/status?json] [-socket unix:///var/run/php-fpm.sock
```

### Socket option

If `-socket` option is set, the plugin reads status from standalone php-fpm service.
`-socket` option is available some notations.

* filepath (e.g., **/var/run/php-fpm.sock**)
* network address (e.g., **localhost:9000**)
* URL (e.g., **tcp://localhost:9000** or **unix:///var/run/php-fpm.sock**)

If not set, the plugin reads status via HTTP server such as Nginx or Apache.

## Example of mackerel-agent.conf

```
[plugin.metrics.php-fpm]
command = "/path/to/mackerel-plugin-php-fpm"
```

## Author
[ariarijp](https://github.com/ariarijp)
