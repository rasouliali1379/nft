GOPATH:=$(shell go env GOPATH)

.PHONY: init
init:
	make config

.PHONY: config
config:
	cp -rf ./config.example.yaml ./config.yaml
	cp -rf ./config.example.yaml ./test/config.test.yaml