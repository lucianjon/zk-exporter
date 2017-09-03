export SHELL:=/usr/bin/env bash -O extglob -c
export GO15VENDOREXPERIMENT:=1
export OS=$(shell uname | tr '[:upper:]' '[:lower:]')

build: GOOS ?= ${OS}
build: GOARCH ?= amd64
build: test
	rm -f zk-exporter
	GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "-X main.buildTime=`date --iso-8601=s` -X main.buildVersion=`git rev-parse HEAD | cut -c-7`" -o zk-exporter .

release-linux:
	GOOS=linux $(MAKE) build
	tar Jcf zk-exporter-`git describe --abbrev=0 --tags`-linux-amd64.txz zk-exporter

release-darwin:
	GOOS=darwin $(MAKE) build
	tar Jcf zk-exporter-`git describe --abbrev=0 --tags`-darwin-amd64.txz zk-exporter

release: clean release-linux release-darwin

test:
	go test -v

clean:
	rm -f zk-exporter
	rm -f zk-exporter-*.txz

run: build
	./zk-exporter
