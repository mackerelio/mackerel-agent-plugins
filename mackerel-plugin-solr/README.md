mackerel-plugin-solr
=====================

Solr custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-solr [-host=<hostname>] [-port=<port>]
```

## Example of mackerel-agent.conf

### Default

```
[plugin.metrics.solr]
command = "/path/to/mackerel-plugin-solr"
```

### Custom

```
[plugin.metrics.solr]
command = "/path/to/mackerel-plugin-solr -host=192.168.33.10 -port=8984"
```

## Dependency Solr URL

- http://{host}:{port}/solr/admin/cores
- http://{host}:{port}/solr/{core name}/admin/mbeans

## Munin Solr plugin

- https://github.com/munin-monitoring/contrib/tree/master/plugins/solr

## Apache Solr official website

- http://lucene.apache.org/solr/

## Test

### Run

```
$ cd mackerel-plugin-solr/lib/
$ go test
```

### Add fixture json

```
$ docker pull solr:x.x
$ docker run --name solr_x -d -p 8983:8983 -t solr:x.x
$ docker exec -it --user=solr solr_x bin/solr create_core -c testcore
$ cd mackerel-plugin-solr/lib/
$ curl -s -S 'http://localhost:8983/solr/admin/cores?wt=json' | jq . > stats/x.x.x/cores.json
$ curl -s -S 'http://localhost:8983/solr/testcore/admin/mbeans?wt=json&stats=true&cat=QUERYHANDLER&key=/update/json&key=/select&key=/update/json/docs&key=/get&key=/update/csv&key=/replication&key=/update&key=/dataimport' | jq . > stats/x.x.x/query.json
$ curl -s -S 'http://localhost:8983/solr/testcore/admin/mbeans?wt=json&stats=true&cat=UPDATEHANDLER&key=/update/json&key=/select&key=/update/json/docs&key=/get&key=/update/csv&key=/replication&key=/update&key=/dataimport' | jq . > stats/x.x.x/update.json
$ curl -s -S 'http://localhost:8983/solr/testcore/admin/mbeans?wt=json&stats=true&cat=REPLICATION&key=/update/json&key=/select&key=/update/json/docs&key=/get&key=/update/csv&key=/replication&key=/update&key=/dataimport' | jq . > stats/x.x.x/replication.json
$ curl -s -S 'http://localhost:8983/solr/testcore/admin/mbeans?wt=json&stats=true&cat=CACHE&key=filterCache&key=perSegFilter&key=queryResultCache&key=documentCache&key=fieldValueCache' | jq . > stats/x.x.x/cache.json
```

### Documents

* [MBean Request Handler](https://cwiki.apache.org/confluence/display/solr/MBean+Request+Handler)
* [Performance Statistics Reference](https://cwiki.apache.org/confluence/display/solr/Performance+Statistics+Reference)
* [Docker Official Repository Solr](https://hub.docker.com/_/solr/)
