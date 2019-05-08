# Changelog

## 0.55.2 (2019-05-08)

* [mysql] add -debug option for troubleshooting #551 (lufia)


## 0.55.1 (2019-03-27)

* [accesslog] don't return any metrics on the first scan #549 (Songmu)
* [haproxy] fix example of haproxy.cfg in README.md #548 (Songmu)


## 0.55.0 (2019-02-13)

* [mongodb] apply for mongodb authenticationDatabase #542 (shibacow)
* [linux] consider the case where the width of the Netid column is only 5 in the output of ss #543 (Songmu)
*  [php-fpm] add option to read status from unix domain socket #541 (lufia)


## 0.54.0 (2019-01-10)

* [redis] Change evicted_keys.Diff to true #539 (lufia)
* [redis] Add evicted_keys metric #537 (lufia)
* Add redash api-key option #535 (kyoshidajp)
* [squid] Add metrics to squid #534 (nabeo)
* [postgres] Enable connection without password #533 (kyoshidajp)


## 0.53.0 (2018-11-12)

* [mysql] Use go-mackerel-plugin instead of go-mackerel-plugin-helper #531 (shibayu36)


## 0.52.0 (2018-10-17)

* Set (default) User-Agent header to HTTP requests #528 (astj)
* Build with Go 1.11 #527 (astj)
* Improve jvm error handling #526 (astj)


## 0.51.1 (2018-08-30)

* [postgres]Ignore error to support Aurora #519 (matsuu)


## 0.51.0 (2018-07-25)

* [mysql] Fix decoding transaction ids from mysql innodb status #521 (itchyny)
* add MSSQL plugin #520 (mattn)


## 0.50.0 (2018-06-20)

* [aws-kinesis-streams] Collect (Write|Read)ProvisionedThroughputExceeded metrics correctly #515 (shibayu36)
* [aws-s3-requests] CloudWatch GetMetricStatics parameters #514 (astj)


## 0.49.0 (2018-05-16)

* [aws-rds]support Aurora PostgreSQL engine #512 (matsuu)
* [aws-rds]fix unit for some metrics #511 (matsuu)
* [aws-rds]add BurstBalance metric #510 (matsuu)
* [linux] fix for collectiong ioDrive(FusionIO) diskstats #509 (hayajo)


## 0.48.0 (2018-04-18)

* [linux] collect disk stats of NVMe devices and ignore virtual/removable devices #506 (hayajo)


## 0.47.0 (2018-04-10)

* [aws-ec2-cpucredit] Add T2 unlimited CPU credit metrics #504 (astj)


## 0.46.0 (2018-03-15)

* [Redis] send uptime #502 (dozen)
* [redis] expired_keys change to diff true #501 (dozen)


## 0.45.0 (2018-03-01)

* [postgres] Add amount of xlog_location change #476 (kizkoh)


## 0.44.0 (2018-02-08)

* [aws-elasticsearch] support metric-{key,label}-prefix #495 (astj)
* [mongodb] Fix warning message on MongoDB 3.4, 3.6 #493 (hayajo)
* Add mackerel-plugin-aws-s3-requests #492 (astj)
* Migrate from `go-mgo/mgo` to `globalsign/mgo` #491 (hayajo)


## 0.43.0 (2018-01-23)

* Setting password via environment variable #488 (hayajo)
* update rpm-v2 task for building Amazon Linux 2 package #485 (hayajo)
* Support BSD #487 (miwarin)
* make `make build` works for some plugins which moved out from this repository #486 (astj)


## 0.42.0 (2018-01-10)

* Move mackerel-plugin-json to other repository #482 (shibayu36)
* Move mackerel-plugin-gearmand #480 (shibayu36)
* Move to mackerelio/mackerel-plugin-gcp-compute-engine #479 (shibayu36)
* [mongodb] fix connections_current metric mongodb-Replica-Set #467 (vfa-cancc)
* [haproxy]support unix domain socket #477 (hbadmin)
* [postgres]state may be null even in old versions #478 (matsuu)
* [uptime] use go-osstat/uptime instead of golib/uptime for getting more accurate uptime #475 (Songmu)
* [mysql] add a hint for -disable_innodb #474 (astj)


## 0.41.1 (2017-12-20)

* [mysql] set Diff: true for some stats which are actually counter values #472 (astj)


## 0.41.0 (2017-12-20)

* [mysql] Fix some InnoDB stats #469 (astj)
* [mysql] Fix message for socket option #468 (utisam)
* MySQL Plugin support Aurora reader node #462 (dozen)


## 0.40.0 (2017-12-12)

* Add h2o to package #464 (astj)
* Redis Plugin supports custom CONFIG command #463 (dozen)
* add mackerel-plugin-h2o #456 (hayajo)
* add defer to closing the response body, and change position it. #461 (qt-luigi)
* add that close the response body #460 (qt-luigi)
* [redis] Add Redis replication delay and lag metrics #455 (kizkoh)


## 0.39.0 (2017-11-28)

* Don't add plugins README which has been moved #453 (astj)
* Improve docker plugin #436 (astj)
* [jvm] Fix remote jvm monitoring #451 (astj)
* Changed README.md of mackerel-plugin-linux #452 (soudai)
* [json] Fix error handling #449 (astj)
* Fix license notice #448 (itchyny)
* [docker] Avoid concurrent map writes by multiple goroutines #446 (astj)
* [aws-ec2-ebs] Do not log "fetched no datapoints" error #445 (astj)
* [kinesis-streams] Use Sum aggregation for Kinesis streams statistics #435 (itchyny)


## 0.38.0 (2017-11-09)

* Improve mackerel-plugin-postgres #437 (astj)
* [docker] Add CPU Percentage metrics #424 (astj)
* [gostats] Use go-mackerel-plugin instead of go-mackerel-plugin-helper #429 (itchyny)
* [mysql]Fix makeBigint calculation #434 (matsuu)
* [cloudfront] add -metric-key-prefix option #433 (fujiwara)


## 0.37.1 (2017-10-26)

* [multicore] Refactor multicore plugin #430 (itchyny)


## 0.37.0 (2017-10-19)

* Implement mackerel-plugin-mcrouter #420 (waniji)
* [uptime] use go-mackerel-plugin instead of using go-mackerel-plugin-helper #428 (Songmu)


## 0.36.0 (2017-10-12)

* Add mackerel-plugin-json #395 (doublemarket)
* [awd-dynamodb] [incompatible] remove `.` from Metrics.Name #423 (astj)
* [unicorn] Support metric-key-prefix #425 (astj)
* [aws-elasticsearch] Improve CloudWatch Statistic type and add some metrics #387 (holidayworking)


## 0.35.0 (2017-10-04)

* [twemproxy] [incompatible] add `-enable-each-server-metrics` option #419 #421 (Songmu)


## 0.34.0 (2017-09-27)

* add mackerel-plugin-flume to package #415 (y-kuno)
* [mysql]add MyISAM related graphs #406 (matsuu)
* add mackerel-plugin-sidekiq to package #417 (syou6162)
* build with Go 1.9 #414 (astj)
* [OpenLDAP] fix get latestCSN #413 (masahide)
* [aws-dynamodb] Add ReadThrottleEvents metric and fill 0 when *ThrottleEvents metrics are not present #409 (astj)


## 0.33.0 (2017-09-20)

* add mackerel-plugin-nvidia-smi to package #411 (syou6162)
* [accesslog] Feature/accesslog/customize parser #410 (karupanerura)
* Fix redundant error by golint in redis.go #408 (shibayu36)
* add flume plugin #396 (y-kuno)
* [mysql]add handler graphs #402 (matsuu)


## 0.32.0 (2017-09-12)

* [memcached] add evicted.reclaimed and evicted.nonzero_evictions #388 (Songmu)
* [mysql]add missed metrics and fix graph definition #390 (matsuu)
* [Redis] fix expired keys #398 (edangelion)
* [accesslog] Fix for scanning long lines #400 (itchyny)


## 0.31.0 (2017-08-30)

* [redis] Change queries metric to diff of "total_commands_processed" #397 (edangelion)
* [aws-dynamodb] Refactor and parallelize CloudWatch request with errgroup #367 (astj)
* [plack] Don't raise errors when parsing JSON fields failed #394 (astj)
* [jmx-jolokia] add value to thread graph #393 (y-kuno)


## 0.30.0 (2017-08-23)

* add mackerel-plugin-openldap to package #391 (astj)
* Add Burst Balance metric for AWS EC2 EBS plugin #384 (ariarijp)
* Add openldap plugin  #374 (masahide)


## 0.29.1 (2017-08-02)

* [solr] Fix a graph definition for Apache Solr's cumulative metric #381 (supercaracal)
* [accesslog] Refine LTSV format detection logic https://github.com/Songmu/axslogparser/pull/8 (Songmu)
* [accesslog] Fix testcase (Percentile logic is Fixed up) #380 (Songmu)


## 0.29.0 (2017-07-26)

* [aws-dynamodb] Add TimeToLiveDeletedItemCount metrics #376 (astj)
* [aws-dynamodb] Adjust options and graph definitions #375 (astj)
* [mysql] Fix graph label prefixes #372 (koooge)


## 0.28.1 (2017-06-28)

* postgres: add metric-key-prefix #363 (edangelion)
* [accesslog] add mackerel-plugin-accesslog #359 (Songmu)
* add mackerel-plugin-aws-dynamodb to package #366 (astj)
* Use mackerelio/golib/logging as logger, not mackerelio/mackerel-agent/logging #365 (astj)
* postgres: collect dbsize only if connectable #361 (mechairoi)
* Support PostgreSQL 9.6 #360 (mechairoi)
* Add sidekiq plugin #354 (littlekbt)


## 0.28.0 (2017-06-14)

* Add aws-dynamodb plugin #349 (astj)
* Implemented mackerel-plugin-redash #355 (yoheimuta)
* Add mackerel-plugin-solr to package #356 (astj)
* Add test cases and fix issues for apache solr #341 (supercaracal)


## 0.27.2 (2017-06-07)

* disable diff on php-opcache.cache_size because they are gauge value #352 (matsuu)
* build with Go 1.8 #350 (Songmu)
* v2 packages (rpm and deb) #348 (Songmu)
* [aws-rds] Fix "Latency" metric label #347 (astj)
* Add AWS Kinesis Firehose Plugin #333 (holidayworking)
* Fixed mackerel-plugin-nginx/README.md #345 (kakakakakku)


## 0.27.1 (2017-05-09)

* [php-fpm] Implement PluginWithPrefix interfarce #338 (astj)
* Use SetTempfileByBasename to support MACKEREL_PLUGIN_WORKDIR #339 (astj)


## 0.27.0 (2017-04-27)

* Add uWSGI vassal plugin #335 (kizkoh)
* add mackerel-plugin-uwsgi-vassal to package #336 (astj)


## 0.26.0 (2017-04-19)

* Add AWS Rekognition Plugin #322 (holidayworking)
* Add aws-kinesis-streams plugin #326 (astj)
* Add AWS Lambda plugin #327 (astj)
* [redis] fix metrics lable #329 (y-kuno)
* Add aws-lambda and aws-kinesis-streams to package #330 (astj)
* Support twemproxy v0.3, Add total_server_error #332 (masahide)


## 0.25.6 (2017-04-06)

* Cross compile by go's native cross build, not by gox #321 (astj)
* fix a label of gostats plugin #323 (itchyny)


## 0.25.5 (2017-03-22)

* add `mackerel-plugin` command #315 (Songmu)
* Add AWS WAF Plugin #316 (holidayworking)
* use new bot token #318 (daiksy)
* use new bot token #319 (daiksy)


## 0.25.4 (2017-02-22)

* Improve gce plugin #313 (astj)


## 0.25.3 (2017-02-16)

* Feature/gcp compute engine #304 (littlekbt)
* [aws-rds] Make it possible to get metrics from Aurora. #307 (TakashiKaga)
* [multicore]fix tempfile path #311 (daiksy)


## 0.25.2 (2017-02-08)

* [aws-rds] fix metric name #306 (TakashiKaga)
* [aws-ses] ses.stats is unit type #308 (holidayworking)
* [aws-cloudfront] Fix regression #295 #309 (astj)


## 0.25.1 (2017-01-25)

* Make more plugins to support MACKEREL_PLUGIN_WORKDIR #301 (astj)
* [jvm] Fix the label and scale #302 (itchyny)
* [aws-rds] Support Aurora metrics and refactoring #303 (sioncojp)


## 0.25.0 (2017-01-04)

* Change directory structure convention of each plugin #289 (Songmu)
* [apache2] fix typo in graphdef #291 (astj)
* [apache2] Change metric name not to end with dot #293 (astj)
* add mackerel-plugin-windows-server-sessions #294 (daiksy)
* migrate from goamz to aws-sdk-go #295 (astj)
* [docker] Add timeout for API request #296 (astj)


## 0.24.0 (2016-11-29)

* Implement mackerel-plugin-aws-ec2 #248 (yyoshiki41)
* [postgres] support Pg9.1 #274 (Songmu)
* Add new nvidia-smi plugin #280 (ksauzz)
* [jvm] Add notice about user to README #281 (astj)
* Implement mackerel-plugin-twemproxy #283 (yoheimuta)
* fix cloudwatch dimensions for elb #284 (ki38sato)
* Change error strings to pass current golint #285 (astj)
* Add mackerel-plugin-twemproxy to package #286 (stefafafan)


## 0.23.1 (2016-10-27)

* [redis] Fix a bug to fetch no metrics of keys and expired #272 (yoheimuta)
* fix: "open file descriptors" property in elasticsearch  #273 (kamijin-fanta)
* [memcached] Supported memcached curr_items metric #275 (kakakakakku)
* [memcached] support new_items metrics #276 (Songmu)
* [redis] s/memoty/memory/ #277 (astj)


## 0.23.0 (2016-10-18)

* mackerel-plugin-linux: Allow to select multiple (but not all) sets of metrics #243 (astj)
* Fixed flag comment of mackerel-plugin-fluentd #257 (kakakakakku)
* Fix postgres.iotime.{blk_read_time,blk_write_time} #259 (mechairoi)
* [Plack] Adopt Plack::Middleware::ServerStatus::Lite 0.35's response #261 (astj)
* build with Go 1.7 #262 (astj)
* Add much graphs/metrics to mackerel-plugin-mysql #264 (netmarkjp)
* [apache2] Support -metric-key-prefix option and get rid of default Tempfile specification #265 (astj)
* [aws-rds] add `-engine` option #266 (Songmu)
* [elasticsearch] Add open_file_descriptors metric in elasticsearch plugin #267 (kamijin-fanta)
* Make *some* plugins to support MACKEREL_PLUGIN_WORKDIR #268 (astj)
* [redis] deal with MACKEREL_PLUGIN_WORKDIR #269 (astj)


## 0.22.1 (2016-09-06)

* Fixed README.md #253 (kakakakakku)
* [memcached] Support -metric-key-prefix option  #255 (astj)


## 0.22.0 (2016-07-14)

* add multicore plugin #246 (daiksy)
* add mackerel-plugin-multicore into package (daiksy)


## 0.21.2 (2016-07-07)

* Fix help message #244 (ariarijp)
* [apache2] update README.md. fix mod_status configuration #245 (Songmu)
* Add some plugins to README #247 (ariarijp)
* follow urfave/cli #249 (Songmu)
* [mysql] support -metric-key-prefix option #250 (Songmu)


## 0.21.1 (2016-06-28)

* build with go 1.6.2 #241 (Songmu)


## 0.21.0 (2016-06-23)

* Add PHP-FPM plugin #226 (ariarijp)
* Support password authentication of Redis #232 (hico-horiuchi)
* Add an option to specify type and id pattern to fluentd plugin #233 (waniji)
* Fix bug:aws-ses #234 (tjinjin)
* fix help link #235 (daiksy)
* xentop: get CPU %, not CPU time/min #236 (hagihala)
* add mackerel-plugin-php-fpm into package #238 (Songmu)


## 0.20.2 (2016-06-09)

* aws-ec2-ebs: Use wildcard in the graph definitions #230 (itchyny)


## 0.20.1 (2016-05-25)

* change signatures of doMain to follow recent codegangsta/cli #227 (Songmu)
* fix README.md of mackerel-plugin-jvm #228 (azusa)


## 0.20.0 (2016-05-10)

* [docker] use goroutine for fetching metrics via API #220 (stanaka)
* add graphite and proc-fd into package #223 (Songmu)


## 0.19.4 (2016-04-20)

* Add mackerel-plugin-graphite (#216) (taku-k)
* Add mackerel plugin proc fd (#207) (taku-k)
* Do not send fluentd metrics of other than the output plugin (#213) (waniji)

## 0.19.3 (2016-04-14)

* [redis] skip to calculate capacity when CONFIG command failed (#214) (Songmu)
* Revert "Revert "use /usr/bin/mackerel-plugin-*"" #208 (Songmu)
* fix: redis plugin panics when redis-server is not installed. #209 (stanaka)
* fix: rpm should not include dir #210 (stanaka)
* [nginx] fix typo #215 (y-kuno)
* Refactoring the release process #217 (stanaka)
