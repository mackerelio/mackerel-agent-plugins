VERBOSE_FLAG = $(if $(VERBOSE),-verbose)

VERSION = $$(git describe --tags --always --dirty) ($$(git name-rev --name-only HEAD))

BUILD_FLAGS = -ldflags "\
	      -X main.Version \"$(VERSION)\" \
	      "

TARGET_OSARCH="linux/386"

build: deps
	mkdir -p build
	for i in mackerel-plugin-*; do \
	  gox $(VERBOSE_FLAG) $(BUILD_FLAGS) \
	    -osarch=$(TARGET_OSARCH) -output build/$$i \
	    github.com/mackerelio/mackerel-agent-plugins/$$i; \
	done

test: testdeps
	go test $(VERBOSE_FLAG) ./...

deps:
	go get -d -v $(VERBOSE_FLAG) ./...

testdeps:
	go get -d -v -t $(VERBOSE_FLAG) ./...

rpm:
	rpmbuild --define "_sourcedir `pwd`" -ba packaging/rpm/mackerel-agent-plugins.spec

deb:
	cp build/mackerel-plugin-* packaging/deb/debian/
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

gox:
	go get github.com/mitchellh/gox
	gox -build-toolchain -osarch=$(TARGET_OSARCH)

clean:
	if [ -d build ]; then \
	  rm -f build/mackerel-plugin-*; \
	  rmdir build; \
	fi

.PHONY: build test deps testdeps rpm deb gox clean
