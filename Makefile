CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep rmdeps
	if test -d src/github.com/whosonfirst/go-whosonfirst-sqlite-geojson; then rm -rf src/github.com/whosonfirst/go-whosonfirst-sqlite-geojson; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-sqlite-geojson
	cp -r index src/github.com/whosonfirst/go-whosonfirst-sqlite-geojson/
	cp -r tables src/github.com/whosonfirst/go-whosonfirst-sqlite-geojson/
	cp -r *.go src/github.com/whosonfirst/go-whosonfirst-sqlite-geojson/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

# if you're wondering about the 'rm -rf' stuff below it's because Go is
# weird... https://vanduuren.xyz/2017/golang-vendoring-interface-confusion/
# (20170912/thisisaaronland)

# see the way we're deleting the vendor-ed version of go-whosonfirst-sqlite
# from go-whosonfirst-index - if we don't do that everything fails with a 
# lot of duplicate symbol errors (20180206/thisisaaronland)

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-spatialite"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"
	@GOPATH=$(GOPATH) go get -u "github.com/twpayne/go-geom"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-index"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-log"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-sqlite"
	rm -rf src/github.com/mattn
	rm -rf src/github.com/shaxbee
	rm -rf src/github.com/whosonfirst/go-whosonfirst-sqlite/vendor/github.com/whosonfirst/go-whosonfirst-log
	rm -rf src/github.com/whosonfirst/go-whosonfirst-sqlite/vendor/github.com/whosonfirst/go-whosonfirst-index
	rm -rf src/github.com/whosonfirst/go-whosonfirst-index/vendor/github.com/whosonfirst/go-whosonfirst-sqlite/

vendor-deps: rmdeps deps
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt index/*.go
	go fmt tables/*.go

# remember: if you're doing anything that _queries_ a database you'll need to pass 
# '--tags json1' to go build... (20180331/thisisaaronland)

bin: 	self
	rm -rf bin/*
	@GOPATH=$(GOPATH) go build --tags json1 -o bin/wof-sqlite-index-geojson cmd/wof-sqlite-index-geojson.go
