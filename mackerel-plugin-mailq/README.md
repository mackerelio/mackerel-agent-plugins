## mackerel-plugin-mailq

This plugin collects the number of messages in the processing queue of MTAs.
[Exim](http://www.exim.org/), [Postfix](http://www.postfix.org/) and [qmail](https://cr.yp.to/qmail.html) are supported.

### Usage

```
mackerel-plugin-mailq -M <mta> [-metric-key-prefix <prefix>] [-metric-label-prefix <prefix>]
```

#### Options
```
-M   Name of MTA: exim, postfix, or qmail (mandatory)
-c   Path to queue-printing command, such as postqueue for postfix and qmail-qstat for qmail.
     Give this option if the command has non-standard name. Usually it can be guessed from the -M option
```

### Example agent configuration
```toml
[plugin.metrics.mailq]
command = "/usr/local/bin/mackerel-plugin-mailq -M postfix"
```
