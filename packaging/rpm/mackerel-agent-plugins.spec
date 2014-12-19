%define revision 1

%define _binaries_in_noarch_packages_terminate_build   0

%define __buildroot %{_builddir}/%{name}
#%define __targetdir %{_libexecdir}/mackerel/plugins
%define __targetdir /usr/local/bin

Summary: Monitoring program plugins for Mackerel
Name: mackerel-agent-plugins
Version: 0.6.1
Release: %{revision}
License: Apache-2
Group: Applications/System
URL: https://mackerel.io/
Requires: mackerel-agent >= 0.12.3

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

for i in apache2 aws-ec2-cpucredit aws-elb aws-rds elasticsearch haproxy jvm linux memcached mongodb munin mysql nginx php-apc plack postgres redis snmp squid varnish;do \
    %{__install} -m0755 %{_sourcedir}/build/mackerel-plugin-$i %{buildroot}%{__targetdir}/; \
done

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{__targetdir}

%changelog
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
