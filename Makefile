all: build deb rpm

build:
	go get github.com/mitchellh/gox
	gox -build-toolchain -osarch="linux/386"
	gox -osarch="linux/386" github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-mysql
	cp $HOME/gopath/bin/mackerel-plugin-mysql_linux_386 packaging/deb/debian/mackerel-plugin-mysql

rpm:
	rpmbuild --define "_sourcedir `pwd`" -ba packaging/rpm/mackerel-agent-plugins.spec 

deb:
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us
