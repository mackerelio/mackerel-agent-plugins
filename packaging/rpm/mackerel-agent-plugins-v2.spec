%define _binaries_in_noarch_packages_terminate_build   0

%define __buildroot %{_builddir}/%{name}
%define __targetdir /usr/bin

Summary: Monitoring program metric plugins for Mackerel
Name: mackerel-agent-plugins
Version: %{_version}
Release: 1%{?dist}
License: ASL 2.0
Group: Applications/System
URL: https://mackerel.io/

Source0: README.md
Packager:  Hatena
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
This package provides metric plugins for Mackerel.

%prep

%install
%{__rm} -rf %{buildroot}

%{__mkdir} -p %{buildroot}%{__targetdir}

%{__install} -m0755 %{_sourcedir}/build/mackerel-plugin %{buildroot}%{__targetdir}/

for i in accesslog apache2 aws-ec2-ebs conntrack docker elasticsearch fluentd gostats h2o haproxy inode jmx-jolokia jvm linux mailq memcached mongodb multicore munin mysql nginx openldap php-apc php-fpm php-opcache plack postgres proc-fd rabbitmq redis sidekiq snmp solr squid td-table-count trafficserver twemproxy unicorn uptime uwsgi-vassal varnish; do \
    ln -s ./mackerel-plugin %{buildroot}%{__targetdir}/mackerel-plugin-$i; \
done

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{__targetdir}/*

%changelog
* Mon Oct 20 2025 <mackerel-developers@hatena.ne.jp> - 0.89.2
- added dependabot cooldown (by yseto)
- Bump github.com/mackerelio/mackerel-plugin-json from 1.2.5 to 1.2.7 in /mackerel-plugin-json in the mackerelio group across 1 directory (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.10.0 to 1.11.1 (by dependabot[bot])
- Bump the aws-aws-sdk-go-v2 group across 1 directory with 6 updates (by dependabot[bot])
- grouping some libraries (by yseto)
- Bump github.com/mackerelio/go-mackerel-plugin-helper from 0.1.3 to 0.1.4 (by dependabot[bot])
- Bump github.com/gosnmp/gosnmp from 1.40.0 to 1.42.1 (by dependabot[bot])
- rewrite files not included in the package (by yseto)
- Bump github.com/redis/go-redis/v9 from 9.11.0 to 9.14.0 (by dependabot[bot])
- Bump github.com/urfave/cli from 1.22.16 to 1.22.17 (by dependabot[bot])

* Fri Sep 19 2025 <mackerel-developers@hatena.ne.jp> - 0.89.1
- Bump github.com/mackerelio/mackerel-plugin-mongodb from 1.1.1 to 1.1.2 in /mackerel-plugin-mongodb (by dependabot[bot])
- Bump github.com/mackerelio/go-osstat from 0.2.5 to 0.2.6 (by dependabot[bot])
- Bump github.com/jarcoal/httpmock from 1.3.1 to 1.4.1 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-plugin-mongodb from 1.1.1 to 1.1.2 (by dependabot[bot])
- Bump golang.org/x/text from 0.27.0 to 0.29.0 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.9.6 to 1.12.2 (by dependabot[bot])

* Tue Sep 9 2025 <mackerel-developers@hatena.ne.jp> - 0.89.0
- fix tarball version (by yseto)
- Bump golang.org/x/crypto from 0.31.0 to 0.35.0 in /mackerel-plugin-mongodb (by dependabot[bot])
- Bump actions/setup-go from 5 to 6 (by dependabot[bot])
- Bump actions/checkout from 4 to 5 (by dependabot[bot])
- Bump actions/download-artifact from 4 to 5 (by dependabot[bot])
- add -port option to mackerel-plugin-snmp plugin (by kmuto)
- Bump github.com/mackerelio/mackerel-plugin-mysql from 1.3.0 to 1.3.2 in /mackerel-plugin-mysql (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-plugin-mysql from 1.3.0 to 1.3.2 (by dependabot[bot])
- Bump golang.org/x/oauth2 from 0.7.0 to 0.27.0 in /mackerel-plugin-gcp-compute-engine (by dependabot[bot])
- update github.com/go-ldap/ldap/v3 (by yseto)
- replace to github.com/yusufpapurcu/wmi (by yseto)
- Bump golang.org/x/text from 0.21.0 to 0.27.0 (by dependabot[bot])
- Bump golang.org/x/sync from 0.10.0 to 0.16.0 (by dependabot[bot])
- Remove aws-* plugins (by yseto)
- fix error handling (by yseto)
- Bump github.com/redis/go-redis/v9 from 9.7.0 to 9.11.0 (by dependabot[bot])
- Bump mackerelio/workflows from 1.4.0 to 1.5.0 (by dependabot[bot])
- Bump golang.org/x/net from 0.36.0 to 0.38.0 in /mackerel-plugin-gcp-compute-engine (by dependabot[bot])
- Bump golang.org/x/crypto from 0.23.0 to 0.35.0 (by dependabot[bot])
- Bump github.com/containerd/containerd from 1.6.26 to 1.6.38 (by dependabot[bot])

* Tue Jul 1 2025 <mackerel-developers@hatena.ne.jp> - 0.88.3
- The graph definition for mackerel-plugin-squid was broken. (by do-su-0805)
- update aws-sdk-go-v2 on mackerel-plugin-aws-ec2-ebs (by yseto)
- remove some plugin version data (by yseto)
- Delete CI for Windows Server 2022 (by appare45)

* Fri May 16 2025 <mackerel-developers@hatena.ne.jp> - 0.88.2
- read VERSION from git (by yseto)
- git commit, version from runtime/debug (by yseto)
- use Go 1.24 (by yseto)
- Bump github.com/gosnmp/gosnmp from 1.38.0 to 1.40.0 (by dependabot[bot])
- Add tar.gz packaging (by fujiwara)
- Bump mackerelio/workflows from 1.3.0 to 1.4.0 (by dependabot[bot])

* Mon Mar 31 2025 <mackerel-developers@hatena.ne.jp> - 0.88.1
- replace to newer runner-images (by yseto)
- Bump golang.org/x/net from 0.33.0 to 0.36.0 in /mackerel-plugin-gcp-compute-engine (by dependabot[bot])

* Tue Mar 4 2025 <mackerel-developers@hatena.ne.jp> - 0.88.0
- support memory peak value which is introduced in PHP 8.4 on mackerel-plugin-php-fpm (by kmuto)
- Bump mackerelio/workflows from 1.2.0 to 1.3.0 (by dependabot[bot])

* Mon Jan 27 2025 <mackerel-developers@hatena.ne.jp> - 0.87.0
- Fix CI build (by ne-sachirou)
- Bump golang.org/x/net from 0.23.0 to 0.33.0 in /mackerel-plugin-gcp-compute-engine (by dependabot[bot])
- Bump golang.org/x/crypto from 0.17.0 to 0.31.0 in /mackerel-plugin-mongodb (by dependabot[bot])
- use mackerelio/workflows@v1.2.0 (by yseto)
- Bump golang.org/x/text from 0.17.0 to 0.21.0 (by dependabot[bot])
- Bump golang.org/x/sync from 0.8.0 to 0.10.0 (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.8.1 to 1.10.0 in /mackerel-plugin-nvidia-smi (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.9.0 to 1.10.0 (by dependabot[bot])
- [jvm] add prefix option (by masarasi)
- Bump github.com/redis/go-redis/v9 from 9.5.1 to 9.7.0 (by dependabot[bot])
- Bump github.com/urfave/cli from 1.22.15 to 1.22.16 (by dependabot[bot])
- Bump github.com/mackerelio/go-mackerel-plugin-helper from 0.1.2 to 0.1.3 in /mackerel-plugin-nvidia-smi (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-plugin-mongodb from 1.1.0 to 1.1.1 in /mackerel-plugin-mongodb (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-plugin-json from 1.2.4 to 1.2.5 in /mackerel-plugin-json (by dependabot[bot])
- Bump github.com/opencontainers/runc from 1.1.12 to 1.1.14 (by dependabot[bot])
- Bump github.com/gosnmp/gosnmp from 1.35.0 to 1.38.0 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.51.26 to 1.55.5 (by dependabot[bot])

* Tue Oct 8 2024 <mackerel-developers@hatena.ne.jp> - 0.86.0
- Bump github.com/mackerelio/mackerel-plugin-mysql from 1.2.1 to 1.3.0 in /mackerel-plugin-mysql (by dependabot[bot])
- Watch every module go.mod (by yohfee)
- Bump github.com/mackerelio/mackerel-plugin-mysql from 1.2.2 to 1.3.0 (by dependabot[bot])
- Fix lint (by yohfee)
- [plugin-aws-ec2-ebs] add actions to README need to be allowed in the iam policy (by miztch)

* Thu Aug 8 2024 <mackerel-developers@hatena.ne.jp> - 0.85.0
- Bump golang.org/x/text from 0.15.0 to 0.17.0 (by dependabot[bot])
- [mackerel-plugin-snmp] support 32bit counter overflow. (by yseto)
- Bump github.com/mackerelio/go-mackerel-plugin from 0.1.4 to 0.1.5 (by dependabot[bot])
- Bump github.com/mackerelio/go-mackerel-plugin-helper from 0.1.2 to 0.1.3 (by dependabot[bot])
- Bump github.com/mackerelio/go-osstat from 0.2.4 to 0.2.5 (by dependabot[bot])
- Bump github.com/jmoiron/sqlx from 1.3.5 to 1.4.0 (by dependabot[bot])

* Mon Jul 1 2024 <mackerel-developers@hatena.ne.jp> - 0.84.0
- [mackerel-plugin-php-fpm] add slow_requests_delta metrics because slow_requests is a counter (by Arthur1)

* Wed Jun 12 2024 <mackerel-developers@hatena.ne.jp> - 0.83.0
- [plugin-mailq] mailq/postfix: Fix mail queue count regexp capture when it's `1` (by yongjiajun)
- Bump github.com/urfave/cli from 1.22.14 to 1.22.15 (by dependabot[bot])
- Bump golang.org/x/sync from 0.6.0 to 0.7.0 (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.8.4 to 1.9.0 (by dependabot[bot])

* Tue Apr 23 2024 <mackerel-developers@hatena.ne.jp> - 0.82.1
- Bump github.com/aws/aws-sdk-go from 1.45.11 to 1.51.26 (by dependabot[bot])
- Bump golang.org/x/net from 0.7.0 to 0.23.0 in /mackerel-plugin-gcp-compute-engine (by dependabot[bot])
- Bump github.com/docker/docker from 23.0.1+incompatible to 24.0.9+incompatible (by dependabot[bot])
- Bump google.golang.org/protobuf from 1.28.1 to 1.33.0 in /mackerel-plugin-gcp-compute-engine (by dependabot[bot])
- Bump google.golang.org/protobuf from 1.27.1 to 1.33.0 in /mackerel-plugin-murmur (by dependabot[bot])
- Bump github.com/redis/go-redis/v9 from 9.1.0 to 9.5.1 (by dependabot[bot])
- Bump github.com/opencontainers/runc from 1.1.2 to 1.1.12 (by dependabot[bot])
- Bump github.com/containerd/containerd from 1.6.18 to 1.6.26 (by dependabot[bot])
- Bump golang.org/x/crypto from 0.0.0-20220622213112-05595931fe9d to 0.17.0 in /mackerel-plugin-mongodb (by dependabot[bot])
- Bump google.golang.org/grpc from 1.51.0 to 1.56.3 in /mackerel-plugin-gcp-compute-engine (by dependabot[bot])
- Bump github.com/go-redis/redismock/v9 from 9.0.3 to 9.2.0 (by dependabot[bot])
- Bump github.com/jarcoal/httpmock from 1.3.0 to 1.3.1 (by dependabot[bot])

* Fri Apr 5 2024 <mackerel-developers@hatena.ne.jp> - 0.82.0
- add STRING number support (by kmuto)
- Allow mackerel-plugin-redis to specify username (by mkadokawa-idcf)

* Tue Mar 5 2024 <mackerel-developers@hatena.ne.jp> - 0.81.0
- [mackerel-plugin-windows-server-sessions] Remove the dependency for WMIC command (by Arthur1)

* Tue Feb 27 2024 <mackerel-developers@hatena.ne.jp> - 0.80.0
- Update mysql, mongodb plugins. (by yseto)
- update go version -> 1.22 (by lufia)
- added TLS support on mackerel-plugin-redis (by yseto)
- Bump golang.org/x/crypto from 0.0.0-20220622213112-05595931fe9d to 0.17.0 (by dependabot[bot])
- Bump actions/upload-artifact from 3 to 4 (by dependabot[bot])
- Bump actions/download-artifact from 3 to 4 (by dependabot[bot])
- Bump actions/setup-go from 4 to 5 (by dependabot[bot])

* Fri Sep 22 2023 <mackerel-developers@hatena.ne.jp> - 0.79.0
- Bump github.com/aws/aws-sdk-go from 1.44.239 to 1.45.11 (by dependabot[bot])
- Bump actions/checkout from 3 to 4 (by dependabot[bot])
- use go-redis (by yseto)
- Bump golang.org/x/text from 0.9.0 to 0.13.0 (by dependabot[bot])
- Bump golang.org/x/sync from 0.1.0 to 0.3.0 (by dependabot[bot])
- Bump github.com/urfave/cli from 1.22.12 to 1.22.14 (by dependabot[bot])
- Bump github.com/montanaflynn/stats from 0.7.0 to 0.7.1 (by dependabot[bot])
- Bump github.com/lib/pq from 1.10.7 to 1.10.9 (by dependabot[bot])

* Mon Sep 4 2023 <mackerel-developers@hatena.ne.jp> - 0.78.4
- Fixed Docker CPU Percentage's unusual spike when restarting containers (by Arthur1)
- Remove old rpm packaging (by yseto)

* Thu Jul 13 2023 <mackerel-developers@hatena.ne.jp> - 0.78.3

* Wed Jun 14 2023 <mackerel-developers@hatena.ne.jp> - 0.78.2
- Bump github.com/mackerelio/mackerel-plugin-mysql from 1.2.0 to 1.2.1 in /mackerel-plugin-mysql (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-plugin-mysql from 1.2.0 to 1.2.1 (by dependabot[bot])
- [fluentd] add README.md with extended_metrics option (by sakamossan)

* Wed Apr 12 2023 <mackerel-developers@hatena.ne.jp> - 0.78.1
- Bump golang.org/x/text from 0.8.0 to 0.9.0 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.44.229 to 1.44.239 (by dependabot[bot])
- Fix timeout/error labels of mackerel-plugin-twemproxy (by kmuto)
- [ci] refactor .github/workflows/test.yml (by lufia)
- Bump github.com/aws/aws-sdk-go from 1.44.219 to 1.44.229 (by dependabot[bot])
- Bump github.com/mackerelio/go-osstat from 0.2.3 to 0.2.4 (by dependabot[bot])

* Mon Mar 13 2023 <mackerel-developers@hatena.ne.jp> - 0.78.0
- Bump github.com/fsouza/go-dockerclient from 1.9.5 to 1.9.6 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.44.199 to 1.44.219 (by dependabot[bot])
- Bump golang.org/x/text from 0.7.0 to 0.8.0 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-plugin-mongodb from 1.0.0 to 1.1.0 in /mackerel-plugin-mongodb (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-plugin-mongodb from 1.0.0 to 1.1.0 (by dependabot[bot])
- Fix support APCu and PHP7 or later environment. (by uzulla)

* Mon Feb 27 2023 <mackerel-developers@hatena.ne.jp> - 0.77.0
- Bump github.com/mackerelio/mackerel-plugin-mysql from 1.1.0 to 1.2.0 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-plugin-mysql from 1.1.0 to 1.2.0 in /mackerel-plugin-mysql (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.8.1 to 1.8.2 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.9.4 to 1.9.5 (by dependabot[bot])
- fix syntax for aws rds (by heleeen)
- Bump golang.org/x/net from 0.3.0 to 0.7.0 in /mackerel-plugin-gcp-compute-engine (by dependabot[bot])
- Bump github.com/containerd/containerd from 1.6.14 to 1.6.18 (by dependabot[bot])
- added multiple-os tests (by yseto)
- added import. (by yseto)
- added external repository mongodb (by yseto)

* Wed Feb 15 2023 <mackerel-developers@hatena.ne.jp> - 0.76.1
- Bump golang.org/x/text from 0.6.0 to 0.7.0 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.9.3 to 1.9.4 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.44.190 to 1.44.199 (by dependabot[bot])
- Remove `circle.yml` (by wafuwafu13)
- plugin-mongodb is external repository (by yseto)
- Bump github.com/jarcoal/httpmock from 1.2.0 to 1.3.0 (by dependabot[bot])

* Wed Feb 1 2023 <mackerel-developers@hatena.ne.jp> - 0.76.0
- Bump github.com/aws/aws-sdk-go from 1.44.184 to 1.44.190 (by dependabot[bot])
- Bump actions/setup-go from 2 to 3 (by dependabot[bot])
- Bump actions/cache from 2 to 3 (by dependabot[bot])
- Bump actions/upload-artifact from 2 to 3 (by dependabot[bot])
- Bump actions/download-artifact from 2 to 3 (by dependabot[bot])
- Bump actions/checkout from 2 to 3 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-plugin-mysql from 1.0.0 to 1.1.0 in /mackerel-plugin-mysql (by dependabot[bot])
- Bump github.com/urfave/cli from 1.22.10 to 1.22.12 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-plugin-mysql from 1.0.0 to 1.1.0 (by dependabot[bot])
- Enables Dependabot version updates for GitHub Actions (by Arthur1)
- Remove debian package v1 process. (by yseto)
- fix staticcheck (by yseto)
- Bump github.com/fsouza/go-dockerclient from 1.9.0 to 1.9.3 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.44.116 to 1.44.184 (by dependabot[bot])
- fix gosimple (by yseto)
- fix errcheck, ineffassign. (by yseto)
- ci: enable `gofmt` (by wafuwafu13)
- Bump golang.org/x/text from 0.5.0 to 0.6.0 (by dependabot[bot])
- Bump github.com/montanaflynn/stats from 0.6.6 to 0.7.0 (by dependabot[bot])

* Wed Jan 18 2023 <mackerel-developers@hatena.ne.jp> - 0.75.0
- fix build on current working directory. (by yseto)
- added compile option. (by yseto)
- plugin-mysql is external repository (by yseto)
- packaging: make compression format xz (by lufia)
- accesslog: use `reqtime_microsec` if exists (by wafuwafu13)
- split go.mod for plugins that have previously split repositories (by yseto)

* Tue Dec 20 2022 <mackerel-developers@hatena.ne.jp> - 0.74.0
- use ubuntu-20.04 (by yseto)
- fix packaging process on ci (by yseto)
- refine file rewrite process (by yseto)
- remove xentop on backward compatibility symlink (by yseto)
- fix packaging (by yseto)
- [plugin-elasticsearch] Fix the test for elasticsearch (by lufia)
- added external plugin support (by yseto)
- sort plugins on packaging files. (by yseto)
- Purge less-used plugins from mackerel-agent-plugins package (by lufia)
- Purge mackerel-plugin-nvidia-smi (by lufia)
- Update dependencies (by lufia)
- added test for elasticsearch (by yseto)

* Wed Nov 9 2022 <mackerel-developers@hatena.ne.jp> - 0.73.0
- Fix Elasticsearch plugin. (by fujiwara)
- Bump github.com/stretchr/testify from 1.8.0 to 1.8.1 (by dependabot[bot])
- [mongodb] add metric-key-prefix option (by tukaelu)

* Thu Oct 20 2022 <mackerel-developers@hatena.ne.jp> - 0.72.2
- Bump github.com/aws/aws-sdk-go from 1.44.56 to 1.44.116 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.8.3 to 1.9.0 (by dependabot[bot])
- added timeout on lint (by yseto)
- Bump github.com/urfave/cli from 1.22.9 to 1.22.10 (by dependabot[bot])
- Bump github.com/lib/pq from 1.10.5 to 1.10.7 (by dependabot[bot])
- [uptime] Add tests (by wafuwafu13)
- go.mod from 1.16 to 1.18 (by yseto)
- Bump github.com/mackerelio/go-mackerel-plugin from 0.1.3 to 0.1.4 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.8.1 to 1.8.3 (by dependabot[bot])
- Bump github.com/mackerelio/go-osstat from 0.2.2 to 0.2.3 (by dependabot[bot])
- Improve test (by yseto)
- Bump github.com/go-ldap/ldap/v3 from 3.4.3 to 3.4.4 (by dependabot[bot])
- [aws-ec2-ebs] fix calcurate procedure of Nitro instance (by yseto)
- [plugin-aws-ec2-ebs] fix misuse of period (by lufia)
- Bump github.com/mackerelio/go-mackerel-plugin-helper from 0.1.1 to 0.1.2 (by dependabot[bot])
- Bump github.com/gosnmp/gosnmp from 1.34.0 to 1.35.0 (by dependabot[bot])
- Bump github.com/jarcoal/httpmock from 1.1.0 to 1.2.0 (by dependabot[bot])

* Wed Jul 20 2022 <mackerel-developers@hatena.ne.jp> - 0.72.1
- Bump github.com/aws/aws-sdk-go from 1.44.37 to 1.44.56 (by dependabot[bot])
- Bump github.com/gomodule/redigo from 1.8.8 to 1.8.9 (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.7.1 to 1.8.0 (by dependabot[bot])
- Bump github.com/hashicorp/go-version from 1.4.0 to 1.6.0 (by dependabot[bot])
- Bump github.com/mackerelio/go-mackerel-plugin from 0.1.2 to 0.1.3 (by dependabot[bot])
- Bump github.com/urfave/cli from 1.22.5 to 1.22.9 (by dependabot[bot])
- Bump github.com/jmoiron/sqlx from 1.3.4 to 1.3.5 (by dependabot[bot])

* Wed Jun 22 2022 <mackerel-developers@hatena.ne.jp> - 0.72.0
- Bump github.com/aws/aws-sdk-go from 1.43.36 to 1.44.37 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.8.0 to 1.8.1 (by dependabot[bot])
- [plugin-docker] update README and a description of -command option (by lufia)
- [plugin-docker] fix CPU/Memory metrics on Docker hosts uses cgroup2 (by xruins)
- [plugin-docker] (Breaking) drop 'File' method support (by lufia)
- Bump github.com/fsouza/go-dockerclient from 1.7.10 to 1.8.0 (by dependabot[bot])

* Thu Apr 21 2022 <mackerel-developers@hatena.ne.jp> - 0.71.0
- [plugin-mysql] fix panic when parsing aio stats (by lufia)
- Fix: Input 'job-number' has been deprecated with message: use flag-name instead (by ne-sachirou)

* Thu Apr 14 2022 <mackerel-developers@hatena.ne.jp> - 0.70.6
- Bump github.com/go-ldap/ldap/v3 from 3.4.2 to 3.4.3 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.43.26 to 1.43.36 (by dependabot[bot])
- Bump github.com/lib/pq from 1.10.4 to 1.10.5 (by dependabot[bot])
- [linux] users メトリックに 0 が投稿されない問題を修正 (by masarasi)

* Wed Mar 30 2022 <mackerel-developers@hatena.ne.jp> - 0.70.5
- Bump github.com/mackerelio/go-osstat from 0.2.1 to 0.2.2 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.43.11 to 1.43.26 (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.7.0 to 1.7.1 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.7.9 to 1.7.10 (by dependabot[bot])

* Tue Mar 15 2022 <mackerel-developers@hatena.ne.jp> - 0.70.4
- Bump github.com/aws/aws-sdk-go from 1.42.52 to 1.43.11 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.7.8 to 1.7.9 (by dependabot[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.4.1 to 3.4.2 (by dependabot[bot])

* Wed Feb 16 2022 <mackerel-developers@hatena.ne.jp> - 0.70.3
- Bump github.com/aws/aws-sdk-go from 1.40.59 to 1.42.52 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.7.4 to 1.7.8 (by dependabot[bot])
- Bump github.com/gomodule/redigo from 1.8.6 to 1.8.8 (by dependabot[bot])
- Bump github.com/hashicorp/go-version from 1.3.0 to 1.4.0 (by dependabot[bot])

* Wed Jan 12 2022 <mackerel-developers@hatena.ne.jp> - 0.70.2
- Bump github.com/jarcoal/httpmock from 1.0.8 to 1.1.0 (by dependabot[bot])
- Bump github.com/gomodule/redigo from 1.8.5 to 1.8.6 (by dependabot[bot])
- Bump github.com/gosnmp/gosnmp from 1.32.0 to 1.34.0 (by dependabot[bot])
- Bump github.com/lib/pq from 1.10.3 to 1.10.4 (by dependabot[bot])

* Wed Dec 1 2021 <mackerel-developers@hatena.ne.jp> - 0.70.1
- upgrade to Go 1.17 and others (by lufia)
- Bump github.com/mackerelio/go-osstat from 0.2.0 to 0.2.1 (by dependabot[bot])

* Thu Nov 18 2021 <mackerel-developers@hatena.ne.jp> - 0.70.0
- [plugin-sidekiq] add queue latency metric (by ch1aki)

* Thu Oct 14 2021 <mackerel-developers@hatena.ne.jp> - 0.69.1
- Bump github.com/aws/aws-sdk-go from 1.39.4 to 1.40.59 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.7.3 to 1.7.4 (by dependabot[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.3.0 to 3.4.1 (by dependabot[bot])
- Bump github.com/lib/pq from 1.10.2 to 1.10.3 (by dependabot[bot])
- Bump github.com/mackerelio/golib from 1.2.0 to 1.2.1 (by dependabot[bot])

* Wed Sep 29 2021 <mackerel-developers@hatena.ne.jp> - 0.69.0
- [plugin-redis] migrate redis client library to redigo (by pyto86pri)
- replace library gosnmp in mackerel-plugin-snmp (by yseto)
- [mysql] adapt for mysql5.7's InnoDB SEMAPHORES metrics (by do-su-0805)

* Tue Aug 24 2021 <mackerel-developers@hatena.ne.jp> - 0.68.0
- [plugin-postgres] suppress fetchXlogLocation when wal_level != 'logical' (by handlename)
- added "pending log flushes" (by yseto)
- fix pending_normal_aio_* (by yseto)

* Thu Aug 5 2021 <mackerel-developers@hatena.ne.jp> - 0.67.1
- [accesslog] use Seek to skip the log (by lufia)
- [ci][fluentd] add test.sh (by lufia)

* Wed Jul 28 2021 <mackerel-developers@hatena.ne.jp> - 0.67.0
- [fluentd] add workers option. (by sugy)
- Bump github.com/aws/aws-sdk-go from 1.38.70 to 1.39.4 (by dependabot[bot])

* Thu Jul 15 2021 <mackerel-developers@hatena.ne.jp> - 0.66.0
- [fluentd] add metric-key-prefix option. (by sugy)

* Tue Jul 06 2021 <mackerel-developers@hatena.ne.jp> - 0.65.0
- bump go-mackerel-plugin and go-mackerel-plugin-helper (by astj)
- Bump github.com/aws/aws-sdk-go from 1.38.40 to 1.38.70 (by dependabot[bot])
- Bump github.com/gomodule/redigo from 1.8.4 to 1.8.5 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.7.2 to 1.7.3 (by dependabot[bot])
- Bump github.com/jmoiron/sqlx from 1.3.1 to 1.3.4 (by dependabot[bot])
- [ci][plugin-mysql] check metrics (by lufia)
- [plugin-aws-cloudfront]Replace label name of graph (by wafuwafu13)
- [plugin-redis] avoid to store +Inf (by lufia)

* Wed Jun 23 2021 <mackerel-developers@hatena.ne.jp> - 0.64.3
- [ci][plugin-redis] check metrics (by lufia)
- [ci] run tests on the workflow (by lufia)

* Thu Jun 10 2021 <mackerel-developers@hatena.ne.jp> - 0.64.2
- [plugin-mysql] Fix plugin-mysql to be able to collect metrics by ignoring non-numerical values (by lufia)

* Thu Jun 03 2021 <mackerel-developers@hatena.ne.jp> - 0.64.1
- Bump github.com/montanaflynn/stats from 0.6.5 to 0.6.6 (by dependabot[bot])
- Bump github.com/lib/pq from 1.10.0 to 1.10.2 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.38.1 to 1.38.40 (by dependabot[bot])
- Bump github.com/mackerelio/go-osstat from 0.1.0 to 0.2.0 (by dependabot[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.2.4 to 3.3.0 (by dependabot[bot])
- Use latest go-mackerel-plugin(-helper) (by astj)
- Bump github.com/hashicorp/go-version from 1.2.1 to 1.3.0 (by dependabot[bot])
- upgrade Go 1.14 -> 1.16 (by lufia)

* Wed Mar 24 2021 <mackerel-developers@hatena.ne.jp> - 0.64.0
- Bump github.com/jmoiron/sqlx from 1.2.0 to 1.3.1 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.37.33 to 1.38.1 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.36.28 to 1.37.33 (by dependabot[bot])
- Bump github.com/gomodule/redigo from 1.8.3 to 1.8.4 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.7.0 to 1.7.2 (by dependabot[bot])
- Bump github.com/lib/pq from 1.9.0 to 1.10.0 (by dependabot[bot])
- Bump github.com/montanaflynn/stats from 0.6.4 to 0.6.5 (by dependabot[bot])
- Bump github.com/mackerelio/golib from 1.1.0 to 1.2.0 (by dependabot[bot])
- Bump github.com/jarcoal/httpmock from 1.0.7 to 1.0.8 (by dependabot[bot])
- Support to extended metrics for fluentd 1.6. (by fujiwara)

* Fri Feb 19 2021 <mackerel-developers@hatena.ne.jp> - 0.63.5
- fix incorrect endpoint inspect condition (by yseto)

* Fri Feb 19 2021 <mackerel-developers@hatena.ne.jp> - 0.63.4
- migrate from goamz,go-ses to aws-sdk-go (by yseto)
- replace mackerel-github-release (by yseto)
- [plugin-redis] add test.sh (by lufia)

* Thu Jan 21 2021 <mackerel-developers@hatena.ne.jp> - 0.63.3
- Revert "delete unused Makefile tasks" (by lufia)
- Bump github.com/aws/aws-sdk-go from 1.35.33 to 1.36.28 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.6.6 to 1.7.0 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-plugin-json from 1.2.0 to 1.2.2 (by dependabot[bot])
- Bump github.com/montanaflynn/stats from 0.6.3 to 0.6.4 (by dependabot[bot])
- Bump github.com/jarcoal/httpmock from 1.0.6 to 1.0.7 (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.6.1 to 1.7.0 (by dependabot[bot])
- Bump github.com/mackerelio/golib from 1.0.0 to 1.1.0 (by dependabot[bot])
- Fix labels of innodb_tables_in_use and innodb_locked_tables (by syou6162)
- Bump github.com/fsouza/go-dockerclient from 1.6.5 to 1.6.6 (by dependabot[bot])
- migrate garyburd/redigo -> gomodule/redigo (by lufia)
- migrate CI to GitHub Actions (by lufia)
- [plugin-apache2] add test.sh (by lufia)

* Wed Dec 09 2020 <mackerel-developers@hatena.ne.jp> - 0.63.2
- Bump github.com/lib/pq from 1.8.0 to 1.9.0 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-plugin-json from 1.1.0 to 1.2.0 (by dependabot[bot])
- Enables to parse responses which include string and number from Plack. (by fujiwara)
- Bump github.com/aws/aws-sdk-go from 1.34.32 to 1.35.33 (by dependabot[bot])
- Bump github.com/urfave/cli from 1.22.4 to 1.22.5 (by dependabot-preview[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.2.3 to 3.2.4 (by dependabot-preview[bot])
- Update Dependabot config file (by dependabot-preview[bot])

* Thu Oct 01 2020 <mackerel-developers@hatena.ne.jp> - 0.63.1
- Bump github.com/aws/aws-sdk-go from 1.34.22 to 1.34.32 (by dependabot-preview[bot])
- fix build arch for make rpm-v2-arm (by astj)

* Tue Sep 15 2020 <mackerel-developers@hatena.ne.jp> - 0.63.0
- add arm64 architecture packages, and fix Architecture field of deb (by lufia)
- Bump github.com/aws/aws-sdk-go from 1.33.17 to 1.34.22 (by dependabot-preview[bot])
- Bump github.com/jarcoal/httpmock from 1.0.5 to 1.0.6 (by dependabot-preview[bot])
- Bump github.com/go-redis/redis from 6.15.8+incompatible to 6.15.9+incompatible (by dependabot-preview[bot])
- Build with Go 1.14 (by lufia)
- Bump github.com/aws/aws-sdk-go from 1.33.12 to 1.33.17 (by dependabot-preview[bot])
- Bump github.com/garyburd/redigo from 1.6.0 to 1.6.2 (by dependabot-preview[bot])
- Bump github.com/lib/pq from 1.7.1 to 1.8.0 (by dependabot-preview[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.1.10 to 3.2.3 (by dependabot-preview[bot])
- [plugin-mongodb] add test.sh (by lufia)

* Wed Jul 29 2020 <mackerel-developers@hatena.ne.jp> - 0.62.0
- [plugin-linux] Do not skip diskstat metrics with Linux kernel 4.18+ (by astj)
- Bump github.com/aws/aws-sdk-go from 1.31.11 to 1.33.12 (by dependabot-preview[bot])
- Bump github.com/lib/pq from 1.7.0 to 1.7.1 (by dependabot-preview[bot])

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
