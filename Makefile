all: build deb rpm

build:
	go get github.com/mitchellh/gox
	gox -build-toolchain -osarch="linux/386"
	gox -osarch="linux/386" -output mackerel-plugin-mysql github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-mysql
	cp mackerel-plugin-mysql packaging/deb/debian/

rpm:
	rpmbuild --define "_sourcedir `pwd`" -ba packaging/rpm/mackerel-agent-plugins.spec 

deb:
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us
