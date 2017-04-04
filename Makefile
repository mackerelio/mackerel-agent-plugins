VERBOSE_FLAG = $(if $(VERBOSE),-verbose)
CURRENT_VERSION = $(shell git log --merges --oneline | perl -ne 'if(m/^.+Merge pull request \#[0-9]+ from .+\/bump-version-([0-9\.]+)/){print $$1;exit}')
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w"

all: lint cover testtool testconvention rpm deb

build: deps
	mkdir -p build
	for i in mackerel-plugin-*; do \
		go build  -ldflags="-s -w" -o build/$$i \
		`pwd | sed -e "s|${GOPATH}/src/||"`/$$i; \
	done

build/mackerel-plugin: deps
	mkdir -p build
	go build -ldflags="-s -w -X main.version=$(CURRENT_VERSION) -X main.gitcommit=$(CURRENT_REVISION)" \
	  -o build/mackerel-plugin

test: testgo lint testtool testconvention

testtool:
	prove tool/releng tool/autotag

testgo: testdeps
	go test $(VERBOSE_FLAG) ./...

testconvention:
	prove -r t/
	go generate ./... && git diff --exit-code || (echo 'please `go generate ./...` and commit them' && false)

deps:
	go get -d -v ./...

testdeps:
	go get -d -v -t ./...
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get github.com/pierrre/gotestcover
	go get github.com/mattn/goveralls

lint: testdeps
	go vet ./...
	golint -set_exit_status ./...

cover: testdeps
	gotestcover -v -covermode=count -coverprofile=.profile.cov -parallelpackages=4 ./...

rpm: build
	make build GOOS=linux GOARCH=386
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${CURRENT_VERSION}" --define "buildarch noarch" -bb packaging/rpm/mackerel-agent-plugins.spec
	make build GOOS=linux GOARCH=amd64
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${CURRENT_VERSION}" --define "buildarch x86_64" -bb packaging/rpm/mackerel-agent-plugins.spec

deb: build
	make build GOOS=linux GOARCH=386
	cp build/mackerel-plugin-* packaging/deb/debian/
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

clean:
	if [ -d build ]; then \
	  rm -f build/mackerel-plugin-*; \
	  rmdir build; \
	fi

release:
	tool/releng

.PHONY: all build test testgo deps testdeps rpm deb clean release lint cover testtool testconvention
