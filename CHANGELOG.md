# Changelog

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
