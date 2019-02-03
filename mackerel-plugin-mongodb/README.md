mackerel-plugin-mongodb
=====================

MongoDB custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-mongodb [-host=<host>] [-port=<port>] [-username=<username>] [-password=<password>] [-tempfile=<tempfile>] [-source=<authenticationDatabase>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.mongodb]
command = "/path/to/mackerel-plugin-mongodb"
```

## Add Role

newer mongodb requre `clusterMonitor` role when executed `db.serverStatus()` command.

so add role `clusterMonitor` to reporter.

```
db.grantRolesToUser(
  "user_id",
  [
  { role: "clusterMonitor", db:"admin"}
  ]
 );
 ```

see https://dba.stackexchange.com/questions/121832/db-serverstatus-got-not-authorized-on-admin-to-execute-command

