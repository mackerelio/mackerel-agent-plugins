VERBOSE_FLAG = $(if $(VERBOSE),-verbose)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)

all: lint cover testconvention rpm deb

build: deps
	mkdir -p build
	for i in mackerel-plugin-*; do \
		go build  -ldflags="-s -w" -o build/$$i \
		`pwd | sed -e "s|${GOPATH}/src/||"`/$$i; \
	done

build/mackerel-plugin: deps
	mkdir -p build
	go build -ldflags="-s -w -X main.gitcommit=$(CURRENT_REVISION)" \
	  -o build/mackerel-plugin

test: testgo lint testconvention

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

rpm: build
	make build GOOS=linux GOARCH=386
	rpmbuild --define "_sourcedir `pwd`"  --define "_version 0.27.1" --define "buildarch noarch" -bb packaging/rpm/mackerel-agent-plugins.spec
	make build GOOS=linux GOARCH=amd64
	rpmbuild --define "_sourcedir `pwd`"  --define "_version 0.27.1" --define "buildarch x86_64" -bb packaging/rpm/mackerel-agent-plugins.spec

deb: build
	make build GOOS=linux GOARCH=386
	cp build/mackerel-plugin-* packaging/deb/debian/
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

release: check-release-deps
	(cd tool && cpanm -qn --installdeps .)
	perl tool/create-release-pullrequest

clean:
	if [ -d build ]; then \
	  rm -f build/mackerel-plugin-*; \
	  rmdir build; \
	fi

.PHONY: all build test testgo deps testdeps rpm deb clean release lint cover testconvention
