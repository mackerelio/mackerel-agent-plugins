%define revision 1

%define __buildroot %{_builddir}/%{name}
%define __targetdir %{_libexecdir}/mackerel/plugins

Summary: Monitoring program plugins for Mackerel
Name: mackerel-agent-plugins
Version: 0.1
Release: %{revision}%{?dist}
License: Apache-2
Group: Applications/System
URL: https://mackerel.io/

Source0: README.md
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
This package provides plugins for Mackerel.

%prep

%build
export GOPATH=${GOPATH:=/tmp/gopath}
rm -rf %{__buildroot}
mkdir -p %{__buildroot}
cd %{__buildroot}
go get github.com/mitchellh/gox
$GOPATH/bin/gox -build-toolchain -osarch="linux/386"
$GOPATH/bin/gox -osarch="linux/386" github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-mysql
$GOPATH/bin/gox -osarch="linux/386" github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-memcached

%install
%{__rm} -rf %{buildroot}

%{__mkdir} -p %{buildroot}%{__targetdir}
%{__install} -m0755 %{__buildroot}/mackerel-plugin-mysql_linux_386 %{buildroot}%{__targetdir}/mackerel-plugin-mysql
%{__install} -m0755 %{__buildroot}/mackerel-plugin-memcached_linux_386 %{buildroot}%{__targetdir}/mackerel-plugin-memcached

%post

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{__targetdir}

%changelog
* Tue Aug 26 2014 <stanaka@hatena.ne.jp> - 0.1
- initial release
