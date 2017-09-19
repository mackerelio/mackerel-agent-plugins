VERSION = 0.32.0
VERBOSE_FLAG = $(if $(VERBOSE),-verbose)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)

GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BINDIR  = build/$(GOOS)/$(GOARCH)

all: lint cover testconvention rpm deb

$(BINDIR)/mackerel-plugin-%: mackerel-plugin-%/lib/*.go
	@if [ ! -d $(BINDIR) ]; then mkdir -p $(BINDIR); fi
	go build -ldflags="-s -w" -o $@ ./`basename $@`

build: deps
	for i in mackerel-plugin-*; do \
	  make $(BINDIR)/$$i; \
	done

build/mackerel-plugin: deps
	mkdir -p build
	go build -ldflags="-s -w -X main.gitcommit=$(CURRENT_REVISION)" \
	  -o $(BINDIR)/mackerel-plugin

test: testgo lint testconvention

testgo: testdeps
	go test $(VERBOSE_FLAG) ./...

testconvention:
	prove -r t/
	go generate ./... && git diff --exit-code || \
	  (echo 'please `go generate ./...` and commit them' && false)

deps:
	go get -d -v ./...

testdeps:
	go get -d -v -t ./...
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get github.com/pierrre/gotestcover
	go get github.com/mattn/goveralls

check-release-deps:
	@have_error=0; \
	for command in cpanm hub ghch gobump; do \
	  if ! command -v $$command > /dev/null; then \
	    have_error=1; \
	    echo "\`$$command\` command is required for releasing"; \
	  fi; \
	done; \
	test $$have_error = 0

lint: testdeps
	go vet ./...
	golint -set_exit_status ./...

cover: testdeps
	gotestcover -v -covermode=count -coverprofile=.profile.cov -parallelpackages=4 ./...

crossbuild-package:
	make build GOOS=linux GOARCH=386
	make build GOOS=linux GOARCH=amd64
	make build/mackerel-plugin GOOS=linux GOARCH=amd64

rpm: rpm-v1 rpm-v2

rpm-v1: crossbuild-package
	rpmbuild --define "_sourcedir `pwd`" --define "_bindir build/linux/386" \
	  --define "_version ${VERSION}" --define "buildarch noarch" \
	  -bb packaging/rpm/mackerel-agent-plugins.spec
	rpmbuild --define "_sourcedir `pwd`" --define "_bindir build/linux/amd64" \
	  --define "_version ${VERSION}" --define "buildarch x86_64" \
	  -bb packaging/rpm/mackerel-agent-plugins.spec

rpm-v2: crossbuild-package
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${VERSION}" \
	  --define "buildarch x86_64" --define "dist .el7.centos" --define "_bindir build/linux/amd64" \
	  -bb packaging/rpm/mackerel-agent-plugins-v2.spec

deb: deb-v1 deb-v2

deb-v1: crossbuild-package
	for i in `cat packaging/deb/debian/source/include-binaries`; do \
	  cp build/linux/386/`basename $$i` packaging/deb/debian/; \
	done
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

deb-v2: crossbuild-package
	cp build/linux/amd64/mackerel-plugin packaging/deb-v2/debian/
	cd packaging/deb-v2 && debuild --no-tgz-check -rfakeroot -uc -us

release: check-release-deps
	(cd tool && cpanm -qn --installdeps .)
	perl tool/create-release-pullrequest

clean:
	@if [ -d build ]; then rm -rfv build; fi

.PHONY: all build test testgo deps testdeps rpm rpm-v1 rpm-v2 deb deb-v1 deb-v2 clean release lint cover testconvention crossbuild-package
