%define _binaries_in_noarch_packages_terminate_build   0

%define __buildroot %{_builddir}/%{name}
%define __targetdir /usr/bin

Summary: Monitoring program metric plugins for Mackerel
Name: mackerel-agent-plugins
Version: %{_version}
Release: 1%{?dist}
License: Apache-2
Group: Applications/System
URL: https://mackerel.io/

Source0: README.md
Packager:  Hatena
BuildArch: %{buildarch}
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
This package provides metric plugins for Mackerel.

%prep

%install
%{__rm} -rf %{buildroot}

%{__mkdir} -p %{buildroot}%{__targetdir}

%{__install} -m0755 %{_sourcedir}/build/mackerel-plugin %{buildroot}%{__targetdir}/

for i in apache2 aws-ec2-cpucredit aws-elasticache aws-elasticsearch aws-elb aws-rds aws-ses conntrack elasticsearch gostats graphite haproxy jmx-jolokia jvm linux mailq memcached mongodb multicore munin mysql nginx php-apc php-fpm php-opcache plack postgres proc-fd rabbitmq redis snmp squid td-table-count trafficserver twemproxy varnish xentop aws-cloudfront aws-ec2-ebs fluentd docker unicorn uptime inode; do \
    ln -s ./mackerel-plugin %{buildroot}%{__targetdir}/mackerel-plugin-$i; \
done

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{__targetdir}/*

%changelog
* Tue May 09 2017 <mackerel-developers@hatena.ne.jp> - 0.27.1-1
- [php-fpm] Implement PluginWithPrefix interfarce (by astj)
- Use SetTempfileByBasename to support MACKEREL_PLUGIN_WORKDIR (by astj)
