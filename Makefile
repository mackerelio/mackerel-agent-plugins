# This VERSION variable indicates the latest tag.
VERSION = $(subst v,,$(shell git describe --abbrev=0 --tags))
VERBOSE_FLAG = $(if $(VERBOSE),-verbose)

GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BINDIR  = build/$(GOOS)/$(GOARCH)

export GO111MODULE=on

.PHONY: all
all: testconvention rpm deb tar

.SECONDEXPANSION:
$(BINDIR)/mackerel-plugin-%: mackerel-plugin-%/main.go $$(wildcard mackerel-plugin-%/lib/*.go)
	@if [ ! -d $(BINDIR) ]; then mkdir -p $(BINDIR); fi
	cd `basename $@` && CGO_ENABLED=0 go build -ldflags="-s -w" -o ../$@

.PHONY: build
build:
	for i in mackerel-plugin-*; do \
	  $(MAKE) $(BINDIR)/$$i; \
	done

build/mackerel-plugin: $(patsubst %,depends_on,$(GOOS)$(GOARCH))
	mkdir -p build
	CGO_ENABLED=0 go build -ldflags="-s -w" \
	  -o build/mackerel-plugin

.PHONY: depends_on
depends_on:
	@:

.PHONY: test
test: testgo testconvention
	./test.bash

.PHONY: testgo
testgo:
	go test $(VERBOSE_FLAG) ./...

.PHONY: testconvention
testconvention:
	prove -r t/
	go generate ./... && git diff --exit-code -- . ':(exclude)go.*' || \
	  (echo 'please `go generate ./...` and commit them' && false)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: rpm
rpm: rpm-v2

.PHONY: rpm-v2
rpm-v2: rpm-v2-x86 rpm-v2-arm

.PHONY: rpm-v2-x86
rpm-v2-x86:
	$(MAKE) build/mackerel-plugin GOOS=linux GOARCH=amd64
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${VERSION}" \
	  --define "buildarch x86_64" --define "dist .el7.centos" \
	  --target x86_64 -bb packaging/rpm/mackerel-agent-plugins-v2.spec
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${VERSION}" \
	  --define "buildarch x86_64" --define "dist .amzn2" \
	  --target x86_64 -bb packaging/rpm/mackerel-agent-plugins-v2.spec

.PHONY: rpm-v2-arm
rpm-v2-arm:
	$(MAKE) build/mackerel-plugin GOOS=linux GOARCH=arm64
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${VERSION}" \
	  --define "buildarch aarch64" --define "dist .el7.centos" \
	  --target aarch64 -bb packaging/rpm/mackerel-agent-plugins-v2.spec
	rpmbuild --define "_sourcedir `pwd`"  --define "_version ${VERSION}" \
	  --define "buildarch aarch64" --define "dist .amzn2" \
	  --target aarch64 -bb packaging/rpm/mackerel-agent-plugins-v2.spec

NEXT_VERSION := $(shell grep -o -e "[0-9]\+.[0-9]\+.[0-9]\+-[0-9]" "packaging/deb-v2/debian/changelog" | head -1 | sed 's/-.*$$//')

.PHONY: deb
deb: deb-v2-x86 deb-v2-arm

.PHONY: deb-v2-x86
deb-v2-x86:
	git clean -f -d ./packaging
	$(MAKE) build/mackerel-plugin GOOS=linux GOARCH=amd64
	cp build/mackerel-plugin packaging/deb-v2/debian/
	cp -f packaging/dummy-empty.tar.gz packaging/mackerel-agent-plugins_${NEXT_VERSION}.orig.tar.gz
	cd packaging/deb-v2 && debuild --no-tgz-check -rfakeroot -uc -us

.PHONY: deb-v2-arm
deb-v2-arm:
	git clean -f -d ./packaging
	$(MAKE) build/mackerel-plugin GOOS=linux GOARCH=arm64
	cp build/mackerel-plugin packaging/deb-v2/debian/
	cp -f packaging/dummy-empty.tar.gz packaging/mackerel-agent-plugins_${NEXT_VERSION}.orig.tar.gz
	cd packaging/deb-v2 && debuild --no-tgz-check -rfakeroot -uc -us -aarm64

.PHONY: tar
tar: tar-x86 tar-arm

.PHONY: tar-x86
tar-x86:
	git clean -f -d ./packaging
	$(MAKE) build/mackerel-plugin GOOS=linux GOARCH=amd64
	mkdir -p packaging/tar/build/mackerel-agent-plugins-$(VERSION)-x86_64
	cp README.md CHANGELOG.md build/mackerel-plugin packaging/tar/build/mackerel-agent-plugins-$(VERSION)-x86_64/
	cd packaging/tar && VERSION=$(VERSION) ARCH=x86_64 ./build.sh

.PHONY: tar-arm
tar-arm:
	git clean -f -d ./packaging
	$(MAKE) build/mackerel-plugin GOOS=linux GOARCH=arm64
	mkdir -p packaging/tar/build/mackerel-agent-plugins-$(VERSION)-arm64
	cp README.md CHANGELOG.md build/mackerel-plugin packaging/tar/build/mackerel-agent-plugins-$(VERSION)-arm64/
	cd packaging/tar && VERSION=$(VERSION) ARCH=arm64 ./build.sh

.PHONY: clean
clean:
	@if [ -d build ]; then rm -rfv build packaging/tar/build; fi

.PHONY: update
update:
	go get -u ./...
	go mod tidy
