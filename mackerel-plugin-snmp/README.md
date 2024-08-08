mackerel-plugin-snmp
=====================

SNMP V2c custom metrics plugin for mackerel.io agent.

## Synopsis
can specify multiple metric-definitions in the form of `OID:NAME[:DIFF?][:STACK?]` args.

```shell
mackerel-plugin-snmp [-name=<graph-name>] [-unit=<graph-unit>] [-host=<host>] [-community=<snmp-v2c-community>] [-tempfile=<tempfile>] 'OID:NAME[:DIFF?][:STACK?]' ['OID:NAME[:DIFF?][:STACK?][:COUNTER?]' ...]
 
```

### What is the argument `COUNTER`

If the value is a 32-bit counter, specify uint32.

supported types:

- uint32
- uint64

## Example of mackerel-agent.conf

```
[plugin.metrics.pps]
command = "/path/to/mackerel-plugin-snmp -name='pps' -community='private' '.1.3.6.1.2.1.31.1.1.1.7.2:eth01in:1:0' '.1.3.6.1.2.1.31.1.1.1.11.2:eth01out:1:0'"
```
