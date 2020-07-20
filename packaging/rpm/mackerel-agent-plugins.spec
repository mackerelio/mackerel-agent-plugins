%define revision 1

%define _binaries_in_noarch_packages_terminate_build   0

%define __buildroot %{_builddir}/%{name}
#%define __targetdir %{_libexecdir}/mackerel/plugins
%define __targetdir /usr/bin
%define __oldtargetdir /usr/local/bin

Summary: Monitoring program plugins for Mackerel
Name: mackerel-agent-plugins
Version: %{_version}
Release: %{revision}
License: Apache-2
Group: Applications/System
URL: https://mackerel.io/

Source0: README.md
Packager:  Hatena
BuildArch: %{buildarch}
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
This package provides plugins for Mackerel.

%prep

%install
%{__rm} -rf %{buildroot}

%{__mkdir} -p %{buildroot}%{__targetdir}

for i in accesslog apache2 aws-dynamodb aws-ec2-cpucredit aws-elasticache aws-elasticsearch aws-elb aws-kinesis-streams aws-lambda aws-rds aws-ses aws-s3-requests conntrack elasticsearch flume gostats graphite haproxy jmx-jolokia jvm linux mailq memcached mongodb multicore munin mysql nginx nvidia-smi openldap php-apc php-fpm php-opcache plack postgres proc-fd solr rabbitmq redis sidekiq snmp squid td-table-count trafficserver twemproxy uwsgi-vassal varnish xentop aws-cloudfront aws-ec2-ebs fluentd docker unicorn uptime inode h2o; do \
    %{__install} -m0755 %{_sourcedir}/%{_bindir}/mackerel-plugin-$i %{buildroot}%{__targetdir}/; \
done

%{__install} -d -m755 %{buildroot}%{__oldtargetdir}
for i in apache2 aws-ec2-cpucredit aws-elasticache aws-elasticsearch aws-elb aws-rds aws-ses conntrack elasticsearch gostats haproxy jmx-jolokia jvm linux mailq memcached mongodb munin mysql nginx php-apc php-opcache plack postgres rabbitmq redis snmp squid td-table-count trafficserver varnish xentop aws-cloudfront aws-ec2-ebs fluentd docker unicorn uptime inode; \
do \
    ln -s ../../bin/mackerel-plugin-$i %{buildroot}%{__oldtargetdir}/mackerel-plugin-$i; \
done

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{__targetdir}/*
%{__oldtargetdir}/*

%changelog
* Mon Jul 20 2020 <mackerel-developers@hatena.ne.jp> - 0.61.0
- [plugin-mysql] Fix to send Bytes_sent and Bytes_received correctly (by shibayu36)
- Bump github.com/lib/pq from 1.6.0 to 1.7.0 (by dependabot-preview[bot])
- Bump github.com/hashicorp/go-version from 1.2.0 to 1.2.1 (by dependabot-preview[bot])
- [plugin-accesslog] allow fields that cannot be parsed but unused by the plugin (by susisu)
- Bump github.com/go-redis/redis from 6.15.7+incompatible to 6.15.8+incompatible (by dependabot-preview[bot])
- Bump github.com/stretchr/testify from 1.6.0 to 1.6.1 (by dependabot-preview[bot])
- Update aws-sdk-go to 1.31.11 (by astj)
- Bump github.com/aws/aws-sdk-go from 1.30.27 to 1.31.7 (by dependabot-preview[bot])
- Bump github.com/lib/pq from 1.5.2 to 1.6.0 (by dependabot-preview[bot])
- Bump github.com/stretchr/testify from 1.5.1 to 1.6.0 (by dependabot-preview[bot])
- [plugin-postgres] add test.sh (by lufia)
- Bump github.com/aws/aws-sdk-go from 1.30.7 to 1.30.27 (by dependabot-preview[bot])

* Thu May 14 2020 <mackerel-developers@hatena.ne.jp> - 0.60.2
- Bump github.com/fsouza/go-dockerclient from 1.6.3 to 1.6.5 (by dependabot-preview[bot])
- Bump github.com/lib/pq from 1.3.0 to 1.5.2 (by dependabot-preview[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.1.8 to 3.1.10 (by dependabot-preview[bot])
- Bump github.com/Songmu/axslogparser from 1.2.0 to 1.3.0 (by dependabot-preview[bot])
- ignore diffs of go.mod and go.sum (by lufia)
- Bump go-mackerel-plugin{,-helper} (by astj)
- Bump github.com/aws/aws-sdk-go from 1.29.24 to 1.30.7 (by dependabot-preview[bot])
- Bump github.com/urfave/cli from 1.22.2 to 1.22.4 (by dependabot-preview[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.1.7 to 3.1.8 (by dependabot-preview[bot])
- Add documents for testing (by lufia)

* Fri Apr 03 2020 <mackerel-developers@hatena.ne.jp> - 0.60.1
- Bump github.com/jarcoal/httpmock from 1.0.4 to 1.0.5 (by dependabot-preview[bot])
- Bump github.com/aws/aws-sdk-go from 1.29.14 to 1.29.24 (by dependabot-preview[bot])
- Bump github.com/aws/aws-sdk-go from 1.28.13 to 1.29.14 (by dependabot-preview[bot])
- Bump github.com/montanaflynn/stats from 0.5.0 to 0.6.3 (by dependabot-preview[bot])
- Bump github.com/fsouza/go-dockerclient from 1.6.0 to 1.6.3 (by dependabot-preview[bot])
- Bump github.com/stretchr/testify from 1.4.0 to 1.5.1 (by dependabot-preview[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.1.5 to 3.1.7 (by dependabot-preview[bot])
- Bump github.com/go-redis/redis from 6.15.6+incompatible to 6.15.7+incompatible (by dependabot-preview[bot])
- Bump github.com/aws/aws-sdk-go from 1.28.5 to 1.28.13 (by dependabot-preview[bot])

* Wed Feb 05 2020 <mackerel-developers@hatena.ne.jp> - 0.60.0
- [varnish] remove unnecessary printf (by lufia)
- [varnish] Added metrics for Transient storage (by cohalz)
- [varnish] Added backend_reuse and backend_recycle metrics (by cohalz)
- rename: github.com/motemen/gobump -> github.com/x-motemen/gobump (by lufia)

* Wed Jan 22 2020 <mackerel-developers@hatena.ne.jp> - 0.59.1
- Bump github.com/aws/aws-sdk-go from 1.27.0 to 1.28.5 (by dependabot-preview[bot])
- Bump github.com/aws/aws-sdk-go from 1.26.6 to 1.27.0 (by dependabot-preview[bot])
- add .dependabot/config.yml (by lufia)
- refactor Makefile and update dependencies (by lufia)

* Thu Oct 24 2019 <mackerel-developers@hatena.ne.jp> - 0.59.0
- Build with Go 1.12.12
- Update dependencies (by astj)
- [solr] Fix several graph definitions (by supercaracal)
- [doc]add repository policy (by lufia)
- [jvm] Add the PerfDisableSharedMem JVM option issue to README (by supercaracal)

* Thu Aug 29 2019 <mackerel-developers@hatena.ne.jp> - 0.58.0
- [solr] Add version 7.x and 8.x support (by supercaracal)
- [sidekiq] add option for redis-namespace (by y-imaida)
- add fakeroot to build dependencies (by susisu)

* Mon Jul 22 2019 <mackerel-developers@hatena.ne.jp> - 0.57.0
- [jvm] fix jinfo command timed out error is not logged (by susisu)
- Build with Go 1.12 (by astj)
-  [plugin-jvm] added CGC and CGCT metrics and fixed parsing problem on them (by lufia)

* Tue Jun 11 2019 <mackerel-developers@hatena.ne.jp> - 0.56.0
- support go modules (by astj)
- [plugin-jvm] prefer ${JAVA_HOME}/bin/j** if JAVA_HOME is set (by astj)

* Wed May 08 2019 <mackerel-developers@hatena.ne.jp> - 0.55.2
- [mysql] add -debug option for troubleshooting (by lufia)

* Wed Mar 27 2019 <mackerel-developers@hatena.ne.jp> - 0.55.1
- [accesslog] don't return any metrics on the first scan (by Songmu)
- [haproxy] fix example of haproxy.cfg in README.md (by Songmu)

* Wed Feb 13 2019 <mackerel-developers@hatena.ne.jp> - 0.55.0
- [mongodb] apply for mongodb authenticationDatabase (by shibacow)
- [linux] consider the case where the width of the Netid column is only 5 in the output of ss (by Songmu)
-  [php-fpm] add option to read status from unix domain socket (by lufia)

* Thu Jan 10 2019 <mackerel-developers@hatena.ne.jp> - 0.54.0
- [redis] Change evicted_keys.Diff to true (by lufia)
- [redis] Add evicted_keys metric (by lufia)
- Add redash api-key option (by kyoshidajp)
- [squid] Add metrics to squid (by nabeo)
- [postgres] Enable connection without password (by kyoshidajp)

* Mon Nov 12 2018 <mackerel-developers@hatena.ne.jp> - 0.53.0
- [mysql] Use go-mackerel-plugin instead of go-mackerel-plugin-helper (by shibayu36)

* Wed Oct 17 2018 <mackerel-developers@hatena.ne.jp> - 0.52.0
- Set (default) User-Agent header to HTTP requests (by astj)
- Build with Go 1.11 (by astj)
- Improve jvm error handling (by astj)

* Thu Aug 30 2018 <mackerel-developers@hatena.ne.jp> - 0.51.1
- [postgres]Ignore error to support Aurora (by matsuu)

* Wed Jul 25 2018 <mackerel-developers@hatena.ne.jp> - 0.51.0
- [mysql] Fix decoding transaction ids from mysql innodb status (by itchyny)
- add MSSQL plugin (by mattn)

* Wed Jun 20 2018 <mackerel-developers@hatena.ne.jp> - 0.50.0
- [aws-kinesis-streams] Collect (Write|Read)ProvisionedThroughputExceeded metrics correctly (by shibayu36)
- [aws-s3-requests] CloudWatch GetMetricStatics parameters (by astj)

* Wed May 16 2018 <mackerel-developers@hatena.ne.jp> - 0.49.0
- [aws-rds]support Aurora PostgreSQL engine (by matsuu)
- [aws-rds]fix unit for some metrics (by matsuu)
- [aws-rds]add BurstBalance metric (by matsuu)
- [linux] fix for collectiong ioDrive(FusionIO) diskstats (by hayajo)

* Wed Apr 18 2018 <mackerel-developers@hatena.ne.jp> - 0.48.0
- [linux] collect disk stats of NVMe devices and ignore virtual/removable devices (by hayajo)

* Tue Apr 10 2018 <mackerel-developers@hatena.ne.jp> - 0.47.0
- [aws-ec2-cpucredit] Add T2 unlimited CPU credit metrics (by astj)

* Thu Mar 15 2018 <mackerel-developers@hatena.ne.jp> - 0.46.0
- [Redis] send uptime (by dozen)
- [redis] expired_keys change to diff true (by dozen)

* Thu Mar 01 2018 <mackerel-developers@hatena.ne.jp> - 0.45.0
- [postgres] Add amount of xlog_location change (by kizkoh)

* Thu Feb 08 2018 <mackerel-developers@hatena.ne.jp> - 0.44.0
- [aws-elasticsearch] support metric-{key,label}-prefix (by astj)
- [mongodb] Fix warning message on MongoDB 3.4, 3.6 (by hayajo)
- Add mackerel-plugin-aws-s3-requests (by astj)
- Migrate from `go-mgo/mgo` to `globalsign/mgo` (by hayajo)

* Tue Jan 23 2018 <mackerel-developers@hatena.ne.jp> - 0.43.0
- Setting password via environment variable (by hayajo)
- update rpm-v2 task for building Amazon Linux 2 package (by hayajo)
- Support BSD (by miwarin)
- make `make build` works for some plugins which moved out from this repository (by astj)

* Wed Jan 10 2018 <mackerel-developers@hatena.ne.jp> - 0.42.0
- Move mackerel-plugin-json to other repository (by shibayu36)
- Move mackerel-plugin-gearmand (by shibayu36)
- Move to mackerelio/mackerel-plugin-gcp-compute-engine (by shibayu36)
- [mongodb] fix connections_current metric mongodb-Replica-Set (by vfa-cancc)
- [haproxy]support unix domain socket (by hbadmin)
- [postgres]state may be null even in old versions (by matsuu)
- [uptime] use go-osstat/uptime instead of golib/uptime for getting more accurate uptime (by Songmu)
- [mysql] add a hint for -disable_innodb (by astj)

* Wed Dec 20 2017 <mackerel-developers@hatena.ne.jp> - 0.41.1
- [mysql] set Diff: true for some stats which are actually counter values (by astj)

* Wed Dec 20 2017 <mackerel-developers@hatena.ne.jp> - 0.41.0
- [mysql] Fix some InnoDB stats (by astj)
- [mysql] Fix message for socket option (by utisam)
- MySQL Plugin support Aurora reader node (by dozen)

* Tue Dec 12 2017 <mackerel-developers@hatena.ne.jp> - 0.40.0
- Add h2o to package (by astj)
- Redis Plugin supports custom CONFIG command (by dozen)
- add mackerel-plugin-h2o (by hayajo)
- add defer to closing the response body, and change position it. (by qt-luigi)
- add that close the response body (by qt-luigi)
- [redis] Add Redis replication delay and lag metrics (by kizkoh)

* Tue Nov 28 2017 <mackerel-developers@hatena.ne.jp> - 0.39.0
- Don't add plugins README which has been moved (by astj)
- Improve docker plugin (by astj)
- [jvm] Fix remote jvm monitoring (by astj)
- Changed README.md of mackerel-plugin-linux (by soudai)
- [json] Fix error handling (by astj)
- Fix license notice (by itchyny)
- [docker] Avoid concurrent map writes by multiple goroutines (by astj)
- [aws-ec2-ebs] Do not log "fetched no datapoints" error (by astj)
- [kinesis-streams] Use Sum aggregation for Kinesis streams statistics (by itchyny)

* Thu Nov 09 2017 <mackerel-developers@hatena.ne.jp> - 0.38.0
- Improve mackerel-plugin-postgres (by astj)
- [docker] Add CPU Percentage metrics (by astj)
- [gostats] Use go-mackerel-plugin instead of go-mackerel-plugin-helper (by itchyny)
- [mysql]Fix makeBigint calculation (by matsuu)
- [cloudfront] add -metric-key-prefix option (by fujiwara)

* Thu Oct 26 2017 <mackerel-developers@hatena.ne.jp> - 0.37.1
- [multicore] Refactor multicore plugin (by itchyny)

* Thu Oct 19 2017 <mackerel-developers@hatena.ne.jp> - 0.37.0
- Implement mackerel-plugin-mcrouter (by waniji)
- [uptime] use go-mackerel-plugin instead of using go-mackerel-plugin-helper (by Songmu)

* Thu Oct 12 2017 <mackerel-developers@hatena.ne.jp> - 0.36.0
- Add mackerel-plugin-json (by doublemarket)
- [awd-dynamodb] [incompatible] remove `.` from Metrics.Name (by astj)
- [unicorn] Support metric-key-prefix (by astj)
- [aws-elasticsearch] Improve CloudWatch Statistic type and add some metrics (by holidayworking)

* Wed Oct 04 2017 <mackerel-developers@hatena.ne.jp> - 0.35.0
- [twemproxy] [incompatible] add `-enable-each-server-metrics` option (by Songmu)

* Wed Sep 27 2017 <mackerel-developers@hatena.ne.jp> - 0.34.0
- add mackerel-plugin-flume to package (by y-kuno)
- [mysql]add MyISAM related graphs (by matsuu)
- add mackerel-plugin-sidekiq to package (by syou6162)
- build with Go 1.9 (by astj)
- [OpenLDAP] fix get latestCSN (by masahide)
- [aws-dynamodb] Add ReadThrottleEvents metric and fill 0 when *ThrottleEvents metrics are not present (by astj)

* Wed Sep 20 2017 <mackerel-developers@hatena.ne.jp> - 0.33.0
- add mackerel-plugin-nvidia-smi to package (by syou6162)
- [accesslog] Feature/accesslog/customize parser (by karupanerura)
- Fix redundant error by golint in redis.go (by shibayu36)
- add flume plugin (by y-kuno)
- [mysql]add handler graphs (by matsuu)

* Tue Sep 12 2017 <mackerel-developers@hatena.ne.jp> - 0.32.0
- [memcached] add evicted.reclaimed and evicted.nonzero_evictions (by Songmu)
- [mysql]add missed metrics and fix graph definition (by matsuu)
- [Redis] fix expired keys (by edangelion)
- [accesslog] Fix for scanning long lines (by itchyny)

* Wed Aug 30 2017 <mackerel-developers@hatena.ne.jp> - 0.31.0
- [redis] Change queries metric to diff of "total_commands_processed" (by edangelion)
- [aws-dynamodb] Refactor and parallelize CloudWatch request with errgroup (by astj)
- [plack] Don't raise errors when parsing JSON fields failed (by astj)
- [jmx-jolokia] add value to thread graph (by y-kuno)

* Wed Aug 23 2017 <mackerel-developers@hatena.ne.jp> - 0.30.0
- add mackerel-plugin-openldap to package (by astj)
- Add Burst Balance metric for AWS EC2 EBS plugin (by ariarijp)
- Add openldap plugin  (by masahide)

* Wed Aug 02 2017 <mackerel-developers@hatena.ne.jp> - 0.29.1
- [solr] Fix a graph definition for Apache Solr's cumulative metric (by supercaracal)
- [accesslog] Refine LTSV format detection logic (by Songmu)
- [accesslog] Fix testcase (Percentile logic is Fixed up) (by Songmu)

* Wed Jul 26 2017 <mackerel-developers@hatena.ne.jp> - 0.29.0
- [aws-dynamodb] Add TimeToLiveDeletedItemCount metrics (by astj)
- [aws-dynamodb] Adjust options and graph definitions (by astj)
- [mysql] Fix graph label prefixes (by koooge)

* Wed Jun 28 2017 <mackerel-developers@hatena.ne.jp> - 0.28.1
- postgres: add metric-key-prefix (by edangelion)
- [accesslog] add mackerel-plugin-accesslog (by Songmu)
- add mackerel-plugin-aws-dynamodb to package (by astj)
- Use mackerelio/golib/logging as logger, not mackerelio/mackerel-agent/logging (by astj)
- postgres: collect dbsize only if connectable (by mechairoi)
- Support PostgreSQL 9.6 (by mechairoi)
- Add sidekiq plugin (by littlekbt)

* Wed Jun 14 2017 <mackerel-developers@hatena.ne.jp> - 0.28.0
- Add aws-dynamodb plugin (by astj)
- Implemented mackerel-plugin-redash (by yoheimuta)
- Add mackerel-plugin-solr to package (by astj)
- Add test cases and fix issues for apache solr (by supercaracal)

* Wed Jun 07 2017 <mackerel-developers@hatena.ne.jp> - 0.27.2
- disable diff on php-opcache.cache_size because they are gauge value (by matsuu)
- build with Go 1.8 (by Songmu)
- v2 packages (rpm and deb) (by Songmu)
- [aws-rds] Fix "Latency" metric label (by astj)
- Add AWS Kinesis Firehose Plugin (by holidayworking)
- Fixed mackerel-plugin-nginx/README.md (by kakakakakku)

* Tue May 09 2017 <mackerel-developers@hatena.ne.jp> - 0.27.1-1
- [php-fpm] Implement PluginWithPrefix interfarce (by astj)
- Use SetTempfileByBasename to support MACKEREL_PLUGIN_WORKDIR (by astj)

* Thu Apr 27 2017 <mackerel-developers@hatena.ne.jp> - 0.27.0-1
- Add uWSGI vassal plugin (by kizkoh)
- add mackerel-plugin-uwsgi-vassal to package (by astj)

* Wed Apr 19 2017 <mackerel-developers@hatena.ne.jp> - 0.26.0-1
- Add AWS Rekognition Plugin (by holidayworking)
- Add aws-kinesis-streams plugin (by astj)
- Add AWS Lambda plugin (by astj)
- [redis] fix metrics lable (by y-kuno)
- Add aws-lambda and aws-kinesis-streams to package (by astj)
- Support twemproxy v0.3, Add total_server_error (by masahide)

* Thu Apr 06 2017 <mackerel-developers@hatena.ne.jp> - 0.25.6-1
- Cross compile by go's native cross build, not by gox (by astj)
- fix a label of gostats plugin (by itchyny)

* Wed Mar 22 2017 <mackerel-developers@hatena.ne.jp> - 0.25.5-1
- add `mackerel-plugin` command (by Songmu)
- Add AWS WAF Plugin (by holidayworking)
- use new bot token (by daiksy)
- use new bot token (by daiksy)

* Wed Feb 22 2017 <mackerel-developers@hatena.ne.jp> - 0.25.4-1
- Improve gce plugin (by astj)

* Thu Feb 16 2017 <mackerel-developers@hatena.ne.jp> - 0.25.3-1
- Feature/gcp compute engine (by littlekbt)
- [aws-rds] Make it possible to get metrics from Aurora. (by TakashiKaga)
- [multicore]fix tempfile path (by daiksy)

* Wed Feb 08 2017 <mackerel-developers@hatena.ne.jp> - 0.25.2-1
- [aws-rds] fix metric name (by TakashiKaga)
- [aws-ses] ses.stats is unit type (by holidayworking)
- [aws-cloudfront] Fix regression #295 (by astj)

* Wed Jan 25 2017 <mackerel-developers@hatena.ne.jp> - 0.25.1-1
- Make more plugins to support MACKEREL_PLUGIN_WORKDIR (by astj)
- [jvm] Fix the label and scale (by itchyny)
- [aws-rds] Support Aurora metrics and refactoring (by sioncojp)

* Wed Jan 04 2017 <mackerel-developers@hatena.ne.jp> - 0.25.0-1
- Change directory structure convention of each plugin (by Songmu)
- [apache2] fix typo in graphdef (by astj)
- [apache2] Change metric name not to end with dot (by astj)
- add mackerel-plugin-windows-server-sessions (by daiksy)
- migrate from goamz to aws-sdk-go (by astj)
- [docker] Add timeout for API request (by astj)

* Tue Nov 29 2016 <mackerel-developers@hatena.ne.jp> - 0.24.0-1
- Implement mackerel-plugin-aws-ec2 (by yyoshiki41)
- [postgres] support Pg9.1 (by Songmu)
- Add new nvidia-smi plugin (by ksauzz)
- [jvm] Add notice about user to README (by astj)
- Implement mackerel-plugin-twemproxy (by yoheimuta)
- fix cloudwatch dimensions for elb (by ki38sato)
- Change error strings to pass current golint (by astj)
- Add mackerel-plugin-twemproxy to package (by stefafafan)

* Thu Oct 27 2016 <mackerel-developers@hatena.ne.jp> - 0.23.1-1
- [redis] Fix a bug to fetch no metrics of keys and expired (by yoheimuta)
- fix: "open file descriptors" property in elasticsearch  (by kamijin-fanta)
- [memcached] Supported memcached curr_items metric (by kakakakakku)
- [memcached] support new_items metrics (by Songmu)
- [redis] s/memoty/memory/ (by astj)

* Tue Oct 18 2016 <mackerel-developers@hatena.ne.jp> - 0.23.0-1
- mackerel-plugin-linux: Allow to select multiple (but not all) sets of metrics (by astj)
- Fixed flag comment of mackerel-plugin-fluentd (by kakakakakku)
- Fix postgres.iotime.{blk_read_time,blk_write_time} (by mechairoi)
- [Plack] Adopt Plack::Middleware::ServerStatus::Lite 0.35's response (by astj)
- build with Go 1.7 (by astj)
- Add much graphs/metrics to mackerel-plugin-mysql (by netmarkjp)
- [apache2] Support -metric-key-prefix option and get rid of default Tempfile specification (by astj)
- [aws-rds] add `-engine` option (by Songmu)
- [elasticsearch] Add open_file_descriptors metric in elasticsearch plugin (by kamijin-fanta)
- Make *some* plugins to support MACKEREL_PLUGIN_WORKDIR (by astj)
- [redis] deal with MACKEREL_PLUGIN_WORKDIR (by astj)

* Tue Sep 06 2016 <mackerel-developers@hatena.ne.jp> - 0.22.1-1
- Fixed README.md (by kakakakakku)
- [memcached] Support -metric-key-prefix option  (by astj)

* Thu Jul 14 2016 <mackerel-developers@hatena.ne.jp> - 0.22.0-1
- add multicore plugin (by daiksy)
- add mackerel-plugin-multicore into package (by daiksy)

* Thu Jul 07 2016 <mackerel-developers@hatena.ne.jp> - 0.21.2-1
- Fix help message (by ariarijp)
- [apache2] update README.md. fix mod_status configuration (by Songmu)
- Add some plugins to README (by ariarijp)
- follow urfave/cli (by Songmu)
- [mysql] support -metric-key-prefix option (by Songmu)

* Tue Jun 28 2016 <mackerel-developers@hatena.ne.jp> - 0.21.1-1
- build with go 1.6.2 (by Songmu)

* Thu Jun 23 2016 <mackerel-developers@hatena.ne.jp> - 0.21.0-1
- Add PHP-FPM plugin (by ariarijp)
- Support password authentication of Redis (by hico-horiuchi)
- Add an option to specify type and id pattern to fluentd plugin (by waniji)
- Fix bug:aws-ses (by tjinjin)
- fix help link (by daiksy)
- xentop: get CPU %, not CPU time/min (by hagihala)
- add mackerel-plugin-php-fpm into package (by Songmu)

* Thu Jun 09 2016 <mackerel-developers@hatena.ne.jp> - 0.20.2-1
- aws-ec2-ebs: Use wildcard in the graph definitions (by itchyny)

* Wed May 25 2016 <mackerel-developers@hatena.ne.jp> - 0.20.1-1
- change signatures of doMain to follow recent codegangsta/cli (by Songmu)
- fix README.md of mackerel-plugin-jvm (by azusa)

* Tue May 10 2016 <mackerel-developers@hatena.ne.jp> - 0.20.0-1
- [docker] use goroutine for fetching metrics via API (by stanaka)
- add graphite and proc-fd into package (by Songmu)

* Wed Apr 20 2016 <mackerel-developers@hatena.ne.jp> - 0.19.4-1
- Add mackerel-plugin-graphite (#216) (by taku-k)
- Add mackerel plugin proc fd (#207) (by taku-k)
- Do not send fluentd metrics of other than the output plugin (#213) (by waniji)

* Thu Apr 14 2016 <mackerel-developers@hatena.ne.jp> - 0.19.3-1
- [redis] skip to calculate capacity when CONFIG command failed (by Songmu)
- Revert "Revert "use /usr/bin/mackerel-plugin-*"" (by Songmu)
- fix: redis plugin panics when redis-server is not installed. (by stanaka)
- fix: rpm should not include dir (by stanaka)
- [nginx] fix typo (by y-kuno)
- Refactoring the release process (by stanaka)

* Fri Mar 25 2016 <y.songmu@gmail.com> - 0.19.2
- Revert "use /usr/bin/mackerel-plugin-*" (by Songmu)

* Fri Mar 25 2016 <y.songmu@gmail.com> - 0.19.1
- use /usr/bin/mackerel-plugin-* (by naokibtn)
- use GOARCH=amd64 for now (by Songmu)

* Thu Mar 17 2016 <y.songmu@gmail.com> - 0.19.0
- [docker] Use Docker stats API (by stanaka)
- Add mailq plugin (by hanazuki)
- added mackerel-plugin-gearmand (by karupanerura)
- added capacity metrics for mysql (by karupanerura)
- added capacity metrics for redis (by karupanerura)
- support to metric-key-prefix/metric-label-prefix option for mackerel-plugin-plack (by karupanerura)
- Time out if jps,jinfo,jstat is hanged up (by tom--bo)
- add mailq into package (by Songmu)

* Thu Mar 10 2016 <y.songmu@gmail.com> - 0.18.1
- Fix helper.Tempfile in mysql.go (by hfm)

* Wed Mar 02 2016 <y.songmu@gmail.com> - 0.18.0
- [mysql] care innodb_buffer_pool_instances (by Songmu)
- Add uptime plugin (by Songmu)
- [inode] use `df -iP` on linux (care line break) (by Songmu)
- [mysql] support unix socket (by Songmu)

* Thu Feb 18 2016 <stefafafan@hatena.ne.jp> - 0.17.0
- Add mackerel-plugin-rabbitmq (by haramaki)
- Add metric key and label prefix option to Elasticsearch plugin (by yano3)
- Add jmx jolokia plugin (by y-kuno)
- Add nf(ip)_conntrack plugin (by hfm)
- [memcached] support unix socket (by Songmu)
- use plugin-helper for mackerel-plugin-apache2 (by stanaka)
- add conntrack, jmx-jolokia, rabbitmq into package (by Songmu)

* Thu Feb 04 2016 <y.songmu@gmail.com> - 0.16.0
- Add couple of metrics to Elasticsearch plugin (by ariarijp)
- [jvm] Add confirmation to error message (by tom--bo)
- Add scheme option to Elasticsearch plugin (by yano3)
- Add inode plugin (by itchyny)

* Thu Jan 07 2016 <y.songmu@gmail.com> - 0.15.2
- add unicorn plugin to package (by yano3)

* Thu Jan 07 2016 <y.songmu@gmail.com> - 0.15.1
- use mackerel-plugin-helper for mackerel-plugin-linux (by stanaka)

* Wed Jan 06 2016 <y.songmu@gmail.com> - 0.15.0
- Add mackerel-plugin-unicorn (by linyows)
- Update README (by y-kuno)
- Add mackerel-plugin-solr [not into package] (by supercaracal)
- Add mackerel-plugin-murmur (not into package) (by mikoim)
- add mackerel-plugin-gostats (by Songmu)
- Add graph definition for memcached cache size (by y-kuno)
- Squid: work with squid v3 (by naokibtn)
- add graphs to varnish plugin (by naokibtn)
- When Seconds_Behind_Master is NULL, agent-plugin doesn't send the Seconds_Behind_Master metric. (by norisu0313)
- support mongodb 3.2 (by stanaka)
- rename goserver2gostats and add README (by Songmu)

* Wed Nov 25 2015 <y.songmu@gmail.com> - 0.14.2
- Fix document (by tkuchiki)
- Get memory usage percentage and CMSInitiatingOccupancyFraction when CMS GC is running (by tom--bo)
- follow latest aws-sdk-go (by Songmu)

* Mon Oct 26 2015 <y.songmu@gmail.com> - 0.14.1
- fix index bug in plugin-xentop (by Songmu)

* Mon Oct 26 2015 <daiksy@hatena.ne.jp> - 0.14.0
- Apache Traffic Server Plugin (by naokibtn)
- added plugin for AWS Elasticsearch Service (by hiroakis)
- use wildcard definition & normalize xen names (by naokibtn)
- add graph definition for java8 metaspace (by Songmu)
- add plugins (aws-elasticsearch and trafficserver) into package (by Songmu)

* Thu Oct 15 2015 <itchyny@hatena.ne.jp> - 0.13.2
- reduce binary size (by Songmu)
- remove Config field from FluentPluginMetrics (by Songmu)
- support coreos and amazon linux for docker plugin (by stanaka)
- Add Key prefix option for AWS RDS plugin (by stanaka)

* Fri Sep 25 2015 <y.songmu@gmail.com> - 0.13.1
- [docker] resolve cgroup path in systemd environment (by Songmu)

* Wed Sep 16 2015 <itchyny@hatena.ne.jp> - 0.13.0
- add mackerel-plugin-fluentd (by stanaka)
- add mackerel-plugin-docker (by stanaka)

* Wed Sep 02 2015 <tomohiro68@gmail.com> - 0.12.0
- Plugin for AWS-EC2 EBS (by naokibtn)

* Thu Aug 20 2015 <y.songmu@gmail.com> - 0.11.2
- Fix/mongodb 2.4 or later (by stanaka)

* Thu Aug 13 2015 <tomohiro68@gmail.com> - 0.11.1
- [nginx] specify types of nginx metrics (by stanaka)

* Wed Jul 29 2015 <y.songmu@gmail.com> - 0.11.0
- [redis] support multiple redis instances on one server by using -metric-key-prefix option (by xorphitus)
- [Redis] fix tiny documentation typo (by hiroakis)
- [memcached/mysql/nginx/plack] care counter overflow and reset by using new helper. (by stanaka)
- [MySQL] Fix typo of metric name s/Thread_created/Threads_created/ (by yuuki1)
- [Redis] Support -socket option to specify unix socket (by Songmu)

* Wed Jul 08 2015 <tomohiro68@gmail.com> - 0.10.0
- Can specify database name option when postgresql does not has a database with the same name as the user name. (by azusa)
- Add mackerel-plugin-aws-cloudfront (by najeira)

* Wed Jun 17 2015 <tomohiro68@gmail.com> - 0.9.2
- elasticsearch: add memory size used by lucene segments, which were exposed in Elasticsearch 1.4 (by yshh)
- mysql: better error handling (by stanaka)

* Wed Jun 10 2015 <tomohiro68@gmail.com> - 0.9.1
- Fix an error in the number of requests of Varnish (by mono0x)

* Tue May 12 2015 <y.songmu@gmail.com> - 0.9.0
- Fix elasticache plugin (by ki38sato)
- Add td-table-count plugin (by ariarijp)
- additional help for nginx SSL sites (by obaratch)
- Add Throughput and Network Throughput metrics to AWS RDS Plugin (by ariarijp)
- Feature/haproxy basicauth (by stanaka)

* Wed Apr 08 2015 <y.songmu@gmail.com> - 0.8.1
- nginx: Add a 'header' parameter to command-line flags (by pandax381)
- remove mackerel-agent dependency (by Songmu)

* Thu Apr 02 2015 <y.songmu@gmail.com> - 0.8.0
- mackerel-plugin-xentop (by taketo957)
- Add Metrics for MySQL InnoDB (by koemu)
- apache2: Add a 'header' parameter to cli flags (by pandax381)
- Fix linux users, disk time metrics (by koemu)
- Add README for mackerel-plugin-xentop (by y-uuki)

* Wed Feb 25 2015 <y.songmu@gmail.com> - 0.7.0
- A plugin for PHP OPcache (by yucchiy)
- Add a parameter to specify the LoadBalancerName (by ki38sato)
- A plugin for elasticache (by ki38sato)
- A plugin for AWS SES (Quota, Send Statistics) (by naokibtn)
- Elasticsearch: change "Indices Docs" to stacked graph (by yshh)
- jvm plugin -pidfile option jstat failed (by ta9to)
- Fix some failed unit tests (by ariarijp)

* Tue Jan 20 2015 <y.songmu@gmail.com> - 0.6.3
- Elasticsearch: add evictions, lucene segments memory size (by yshh)
- Filter invalid float values (by krrrr38)

* Thu Dec 25 2014 <yuki.tsubo@gmail.com> - 0.6.2
- Fix the problem jvm plugin does'nt run (by y-uuki)

* Fri Dec 19 2014 <daiksy@hatena.ne.jp> - 0.6.1
- Fix spelling at MySQL Plugin (by koemu)
- Fix requirement module name (by y-uuki)
- Change scale KB into Byte about JVM memory usgae (by y-uuki)
- JVM pidfile option (by y-uuki)
- Add GC time percentage graph (by y-uuki)
- Elasticsearch: not exit but warn when some values cannot be fetched (by naokibtn)
- Fix MySQL int32 overflow (by y-uuki)

* Fri Dec 05 2014 <songmu@hatena.ne.jp> - 0.6.0
- Changed some parameter unit at Apache2 plugin.  (by koemu)
- add IAM Policy requirement (by naokibtn)
- Use `varnishstat` command (by naokibtn)
- plugin for munin (by naokibtn)

* Tue Oct 21 2014 <y_uuki@hatena.ne.jp> - 0.5.0
- Add plugin for HAProxy, Varnish, Squid, SNMP, EC2 CPU Credit, Elasticsearch, JVM, Linux procs, MongoDB, ELB, RDS and PHP-APC
- Fix Plack, Apache2

* Wed Sep 17 2014 <stanaka@hatena.ne.jp> - 0.4.2
- Fix memcached

* Tue Sep 16 2014 <stanaka@hatena.ne.jp> - 0.4.0
- Add plugin for apache2, nginx, plack, postgres, redis

* Tue Aug 27 2014 <stanaka@hatena.ne.jp> - 0.3.1
- Update version string

* Tue Aug 26 2014 <stanaka@hatena.ne.jp> - 0.1
- initial release
