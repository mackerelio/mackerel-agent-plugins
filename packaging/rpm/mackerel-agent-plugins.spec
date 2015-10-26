%define revision 1

%define _binaries_in_noarch_packages_terminate_build   0

%define __buildroot %{_builddir}/%{name}
#%define __targetdir %{_libexecdir}/mackerel/plugins
%define __targetdir /usr/local/bin

Summary: Monitoring program plugins for Mackerel
Name: mackerel-agent-plugins
Version: 0.14.1
Release: %{revision}
License: Apache-2
Group: Applications/System
URL: https://mackerel.io/

Source0: README.md
Packager:  Hatena
BuildArch: noarch
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
This package provides plugins for Mackerel.

%prep

%install
%{__rm} -rf %{buildroot}

%{__mkdir} -p %{buildroot}%{__targetdir}

for i in apache2 aws-ec2-cpucredit aws-elasticache aws-elasticsearch aws-elb aws-rds aws-ses elasticsearch haproxy jvm linux memcached mongodb munin mysql nginx php-apc php-opcache plack postgres redis snmp squid td-table-count trafficserver varnish xentop aws-cloudfront aws-ec2-ebs fluentd docker;do \
    %{__install} -m0755 %{_sourcedir}/build/mackerel-plugin-$i %{buildroot}%{__targetdir}/; \
done

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{__targetdir}

%changelog
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
