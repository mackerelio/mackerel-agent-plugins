%define revision 1

%define _binaries_in_noarch_packages_terminate_build   0

%define __buildroot %{_builddir}/%{name}
#%define __targetdir %{_libexecdir}/mackerel/plugins
%define __targetdir /usr/local/bin

Summary: Monitoring program plugins for Mackerel
Name: mackerel-agent-plugins
Version: 0.3.1
Release: %{revision}%{?dist}
License: Apache-2
Group: Applications/System
URL: https://mackerel.io/
Requires: mackerel-agent >= 0.12.1

Source0: README.md
Packager:  Hatena
BuildArch: noarch
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
This package provides plugins for Mackerel.

%prep

%build
rm -rf %{__buildroot}
mkdir -p %{__buildroot}
cd %{__buildroot}
go get github.com/mitchellh/gox
gox -build-toolchain -osarch="linux/386"
gox -osarch="linux/386" github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-mysql
gox -osarch="linux/386" github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-memcached
gox -osarch="linux/386" github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-nginx

%install
%{__rm} -rf %{buildroot}

%{__mkdir} -p %{buildroot}%{__targetdir}
%{__install} -m0755 %{__buildroot}/mackerel-plugin-mysql_linux_386 %{buildroot}%{__targetdir}/mackerel-plugin-mysql
%{__install} -m0755 %{__buildroot}/mackerel-plugin-memcached_linux_386 %{buildroot}%{__targetdir}/mackerel-plugin-memcached
%{__install} -m0755 %{__buildroot}/mackerel-plugin-memcached_linux_386 %{buildroot}%{__targetdir}/mackerel-plugin-nginx

%post

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{__targetdir}

%changelog
* Tue Aug 27 2014 <stanaka@hatena.ne.jp> - 0.3.1
- Update version string

* Tue Aug 26 2014 <stanaka@hatena.ne.jp> - 0.1
- initial release
