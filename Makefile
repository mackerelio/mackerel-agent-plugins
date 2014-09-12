all: build deb rpm

build:
	mkdir build
	go get github.com/mitchellh/gox
	gox -build-toolchain -osarch="linux/386"
	gox -osarch="linux/386" -output build/mackerel-plugin-mysql github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-mysql
	gox -osarch="linux/386" -output build/mackerel-plugin-memcached github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-memcached
	gox -osarch="linux/386" -output build/mackerel-plugin-nginx github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-nginx
	gox -osarch="linux/386" -output build/mackerel-plugin-postgres github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-postgres
	cp build/mackerel-plugin-* packaging/deb/debian/

rpm:
	rpmbuild --define "_sourcedir `pwd`" -ba packaging/rpm/mackerel-agent-plugins.spec

deb:
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us
