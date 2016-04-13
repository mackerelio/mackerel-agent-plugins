VERBOSE_FLAG = $(if $(VERBOSE),-verbose)

VERSION = $$(git describe --tags --always --dirty) ($$(git name-rev --name-only HEAD))

BUILD_FLAGS = -ldflags "\
	      -s -w \
	      -X main.Version \"$(VERSION)\" \
	      "

TARGET_OSARCH="linux/amd64"

all: lint cover testtool rpm deb

build: deps
	mkdir -p build
	for i in mackerel-plugin-*; do \
	  gox $(VERBOSE_FLAG) $(BUILD_FLAGS) \
	    -osarch=$(TARGET_OSARCH) -output build/$$i \
	    github.com/mackerelio/mackerel-agent-plugins/$$i; \
	done

test: testgo lint testtool

testtool:
	prove tool/releng tool/autotag

testgo: testdeps
	go test $(VERBOSE_FLAG) ./...

deps:
	go get -d -v $(VERBOSE_FLAG) ./...

testdeps:
	go get -d -v -t $(VERBOSE_FLAG) ./...
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get github.com/pierrre/gotestcover
	go get github.com/mattn/goveralls

LINT_RET = .golint.txt
lint: testdeps
	go vet ./...
	rm -f $(LINT_RET)
	golint ./... | tee -a $(LINT_RET)
	test ! -s $(LINT_RET)

cover: testdeps
	gotestcover -v -covermode=count -coverprofile=.profile.cov -parallelpackages=4 ./...

rpm: build
	rpmbuild --define "_sourcedir `pwd`" -ba packaging/rpm/mackerel-agent-plugins.spec

deb: build
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

release:
	tool/releng

.PHONY: all build test testgo deps testdeps rpm deb gox clean release lint cover testtool
