VERBOSE_FLAG = $(if $(VERBOSE),-verbose)

VERSION = $$(git describe --tags --always --dirty) ($$(git name-rev --name-only HEAD))
CURRENT_VERSION = $(shell git log --merges --oneline | perl -ne 'if(m/^.+Merge pull request \#[0-9]+ from .+\/bump-version-([0-9\.]+)/){print $$1;exit}')

BUILD_FLAGS = -ldflags "\
	      -s -w \
	      -X main.Version \"$(VERSION)\" \
	      "

TARGET_OSARCH="linux/amd64"

check-variables:
	echo "CURRENT_VERSION: ${CURRENT_VERSION}"
	echo "TARGET_OSARCH: ${TARGET_OSARCH}"

all: lint cover testtool rpm deb

build: deps
	mkdir -p build
	for i in mackerel-plugin-*; do \
	  gox $(VERBOSE_FLAG) $(BUILD_FLAGS) \
	    -osarch=$(TARGET_OSARCH) -output build/$$i \
			`pwd | sed -e "s|${GOPATH}/src/||"`/$$i; \
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
	TARGET_OSARCH="linux/386" make build
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${CURRENT_VERSION}" --define "buildarch noarch" -bb packaging/rpm/mackerel-agent-plugins.spec
	TARGET_OSARCH="linux/amd64" make build
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${CURRENT_VERSION}" --define "buildarch x86_64" -bb packaging/rpm/mackerel-agent-plugins.spec

deb: build
	TARGET_OSARCH="linux/386" make build
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
