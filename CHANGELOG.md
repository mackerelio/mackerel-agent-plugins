# Changelog

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
