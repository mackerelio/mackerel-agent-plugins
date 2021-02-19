VERSION = 0.63.5
VERBOSE_FLAG = $(if $(VERBOSE),-verbose)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)

GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BINDIR  = build/$(GOOS)/$(GOARCH)

export GO111MODULE=on

.PHONY: all
all: lint cover testconvention rpm deb

.SECONDEXPANSION:
$(BINDIR)/mackerel-plugin-%: mackerel-plugin-%/main.go $$(wildcard mackerel-plugin-%/lib/*.go)
	@if [ ! -d $(BINDIR) ]; then mkdir -p $(BINDIR); fi
	go build -ldflags="-s -w" -o $@ ./`basename $@`

.PHONY: build
build:
	for i in mackerel-plugin-*; do \
	  $(MAKE) $(BINDIR)/$$i; \
	done

build/mackerel-plugin: $(patsubst %,depends_on,$(GOOS)$(GOARCH))
	mkdir -p build
	go build -ldflags="-s -w -X main.gitcommit=$(CURRENT_REVISION)" \
	  -o build/mackerel-plugin

.PHONY: depends_on
depends_on:
	@:

.PHONY: test
test: testgo lint testconvention

.PHONY: testgo
testgo: testdeps
	go test $(VERBOSE_FLAG) ./...

.PHONY: testconvention
testconvention:
	prove -r t/
	go generate ./... && git diff --exit-code -- . ':(exclude)go.*' || \
	  (echo 'please `go generate ./...` and commit them' && false)

.PHONY: testdeps
testdeps:
	cd && go get golang.org/x/lint/golint \
	  golang.org/x/tools/cmd/cover \
	  github.com/mattn/goveralls

.PHONY: check-release-deps
check-release-deps:
	@have_error=0; \
	for command in cpanm hub ghch gobump; do \
	  if ! command -v $$command > /dev/null; then \
	    have_error=1; \
	    echo "\`$$command\` command is required for releasing"; \
	  fi; \
	done; \
	test $$have_error = 0

.PHONY: lint
lint: testdeps
	golint -set_exit_status ./...

.PHONY: cover
cover: testdeps
	#go test -race -covermode=atomic -coverprofile=.profile.cov ./...
	go test -covermode=atomic -coverprofile=.profile.cov ./...

.PHONY: rpm
rpm: rpm-v1 rpm-v2

.PHONY: rpm-v1
rpm-v1:
	$(MAKE) build GOOS=linux GOARCH=386
	rpmbuild --define "_sourcedir `pwd`" --define "_bindir build/linux/386" \
	  --define "_version ${VERSION}" --define "buildarch noarch" \
	  --target noarch -bb packaging/rpm/mackerel-agent-plugins.spec
	$(MAKE) build GOOS=linux GOARCH=amd64
	rpmbuild --define "_sourcedir `pwd`" --define "_bindir build/linux/amd64" \
	  --define "_version ${VERSION}" --define "buildarch x86_64" \
	  --target x86_64 -bb packaging/rpm/mackerel-agent-plugins.spec

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

.PHONY: deb
deb: deb-v1 deb-v2

.PHONY: deb-v1
deb-v1:
	$(MAKE) build GOOS=linux GOARCH=386
	for i in `cat packaging/deb/debian/source/include-binaries`; do \
	  cp build/linux/386/`basename $$i` packaging/deb/debian/; \
	done
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

.PHONY: deb-v2
deb-v2: deb-v2-x86 deb-v2-arm

.PHONY: deb-v2-x86
deb-v2-x86:
	$(MAKE) build/mackerel-plugin GOOS=linux GOARCH=amd64
	cp build/mackerel-plugin packaging/deb-v2/debian/
	cd packaging/deb-v2 && debuild --no-tgz-check -rfakeroot -uc -us

.PHONY: deb-v2-arm
deb-v2-arm:
	$(MAKE) build/mackerel-plugin GOOS=linux GOARCH=arm64
	cp build/mackerel-plugin packaging/deb-v2/debian/
	cd packaging/deb-v2 && debuild --no-tgz-check -rfakeroot -uc -us -aarm64

.PHONY: release
release: check-release-deps
	(cd tool && cpanm -qn --installdeps .)
	perl tool/create-release-pullrequest

.PHONY: clean
clean:
	@if [ -d build ]; then rm -rfv build; fi

.PHONY: update
update:
	go get -u ./...
	go mod tidy
