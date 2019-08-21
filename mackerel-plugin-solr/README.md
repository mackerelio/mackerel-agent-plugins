mackerel-plugin-solr
=====================

Solr custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-solr [-host=<hostname>] [-port=<port>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.solr]
command = "/path/to/mackerel-plugin-solr"
```

You can explicitly specify a host IP address and a port number.

```
[plugin.metrics.solr]
command = "/path/to/mackerel-plugin-solr -host=192.168.33.10 -port=8984"
```

## Supported Versions

* `5.*.*`
* `6.*.*`
* `7.*.*`
* `8.*.*`

## Requirement APIs

* `http://{host}:{port}/solr/admin/info/system`
* `http://{host}:{port}/solr/admin/cores`
* `http://{host}:{port}/solr/{core name}/admin/mbeans`

## Unit tests

```
$ cd mackerel-plugin-solr/lib/
$ go test
```

Prepare fixture files to `mackerel-plugin-solr/lib/stats/x.x.x/*` if you'd like to support new version.

```
$ docker pull solr:x.x
$ docker run --name solr_x -d -p 8983:8983 -t solr:x.x
$ docker exec -it --user=solr solr_x bin/solr create_core -c testcore
$ cd mackerel-plugin-solr/lib/
```

```
$ curl -s -S 'http://localhost:8983/solr/admin/info/system?wt=json' | jq . > stats/x.x.x/system.json
```

```
$ curl -s -S 'http://localhost:8983/solr/admin/cores?wt=json' | jq . > stats/x.x.x/cores.json
```

```
### Solr5 or Solr6
$ curl -s -S 'http://localhost:8983/solr/testcore/admin/mbeans?wt=json&stats=true&cat=QUERYHANDLER&key=/update/json&key=/select&key=/update/json/docs&key=/get&key=/update/csv&key=/replication&key=/update&key=/dataimport' | jq . > stats/x.x.x/query.json
$ curl -s -S 'http://localhost:8983/solr/testcore/admin/mbeans?wt=json&stats=true&cat=UPDATEHANDLER&key=/update/json&key=/select&key=/update/json/docs&key=/get&key=/update/csv&key=/replication&key=/update&key=/dataimport' | jq . > stats/x.x.x/update.json

### Solr7 or Solr8
$ curl -s -S 'http://localhost:8983/solr/testcore/admin/mbeans?wt=json&stats=true&cat=QUERY&key=/update/json&key=/select&key=/update/json/docs&key=/get&key=/update/csv&key=/replication&key=/update&key=/dataimport' | jq . > stats/x.x.x/query.json
$ curl -s -S 'http://localhost:8983/solr/testcore/admin/mbeans?wt=json&stats=true&cat=UPDATE&key=/update/json&key=/select&key=/update/json/docs&key=/get&key=/update/csv&key=/replication&key=/update&key=/dataimport' | jq . > stats/x.x.x/update.json
```

```
$ curl -s -S 'http://localhost:8983/solr/testcore/admin/mbeans?wt=json&stats=true&cat=REPLICATION&key=/update/json&key=/select&key=/update/json/docs&key=/get&key=/update/csv&key=/replication&key=/update&key=/dataimport' | jq . > stats/x.x.x/replication.json
$ curl -s -S 'http://localhost:8983/solr/testcore/admin/mbeans?wt=json&stats=true&cat=CACHE&key=filterCache&key=perSegFilter&key=queryResultCache&key=documentCache&key=fieldValueCache' | jq . > stats/x.x.x/cache.json
```

## Documents

* [Apache Solr official website](http://lucene.apache.org/solr/)
* [MBean Request Handler](https://lucene.apache.org/solr/guide/8_1/mbean-request-handler.html)
* [Performance Statistics Reference](https://lucene.apache.org/solr/guide/8_1/performance-statistics-reference.html)
* [Docker Official Repository Solr](https://hub.docker.com/_/solr/)

## Other softwares

* [Prometheus Exporter](https://github.com/apache/lucene-solr/tree/master/solr/contrib/prometheus-exporter)
* [DataDog Integration](https://github.com/DataDog/integrations-core/tree/master/solr)
* [Munin Plugin](https://github.com/munin-monitoring/contrib/tree/master/plugins/solr)
