all: build deb rpm

build: clean
	mkdir build
	go get -v ./...
	for i in mackerel-plugin-*; do \
	  echo gox -osarch="linux/386" -output build/$$i github.com/mackerelio/mackerel-agent-plugins/$$i; \
	  gox -osarch="linux/386" -output build/$$i github.com/mackerelio/mackerel-agent-plugins/$$i; \
	done
	cp build/mackerel-plugin-* packaging/deb/debian/

rpm:
	rpmbuild --define "_sourcedir `pwd`" -ba packaging/rpm/mackerel-agent-plugins.spec

deb:
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

clean:
	rm -f build/mackerel-plugin-*
	rmdir build

gox:
	go get github.com/mitchellh/gox
	gox -build-toolchain -osarch="linux/386"
