mackerel-plugin-openldap
====

OpenLDAP custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-openldap <-bind=<bind dn>> <-pw=<password>> [-host=<hostname>] [-port=<port number>] [-tls]...

options:
  -bind string
    	bind dn ("cn=config" read user dn)
  -host string
    	Hostname (default "localhost")
  -insecureSkipVerify
    	TLS accepts any certificate.
  -metric-key-prefix string
    	Metric key prefix (default "openldap")
  -port string
    	Port (default "389")
  -pw string
    	bind password
  -replBase string
    	replication base dn
  -replLocalBind string
    	replicationlocalmaster bind dn
  -replLocalPW string
    	replication local bind password
  -replMasterBind string
    	replication master bind dn
  -replMasterHost string
    	replication master hostname
  -replMasterPW string
    	replication master bind password
  -replMasterPort string
    	replication master port (default "389")
  -replMasterTLS
    	replication master TLS(ldaps)
  -tempfile string
    	Temp file name
  -tls
    	TLS(ldaps)
```
## Example of mackerel-agent.conf

```
[plugin.metrics.openldap]
command = '''
	/path/to/mackerel-plugin-openldap 
		-host localhost \
		-port 636 \
		-tls -insecureSkipVerify \
		-bind "cn=Manager,dc=example,dc=net" \
		-pw "password" \
		-replBase "dc=example,dc=net" \
		-replMasterBind "uid=master_user,ou=system,dc=example,dc=net" \
		-replMasterPW "password" \
		-replMasterHost master-server-hostname \
		-replMasterPort 636 \
		-replMasterTLS \
		-replLocalBind "uid=local_user,ou=system,dc=example,dc=net" \
		-replLocalPW "password"
'''
```


