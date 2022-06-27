GOPATH:=$(shell go env GOPATH)

.PHONY: init
init:
	make submodule
	make proto
	make config
	make gitlab

.PHONY: submodule
submodule:
	git submodule init
	git submodule update

.PHONY: proto
proto:
	mkdir -p api/swagger
	mkdir -p api/pb

	protoc -I pb/proto -I pb/aux \
	--go_out api/pb \
	--openapiv2_out=logtostderr=true,repeated_path_param_separator=ssv:./api/swagger \
	--openapiv2_opt use_go_templates=true \
	--openapiv2_opt logtostderr=true \
	--openapiv2_opt use_go_templates=true \
	--grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt paths=source_relative \
	--go_opt paths=source_relative \
	--go-grpc_out api/pb \
	--go-grpc_opt paths=source_relative \
	--grpc-gateway_out api/pb \
	--grpc-gateway_opt paths=source_relative pb/proto/base/*.proto

.PHONY: build
build:
	go build -o base *.go

.PHONY: test
test:
	go test -v ./... -cover -race

.PHONY: vendor
vendor:
	go get ./...
	go mod vendor
	go mod verify

.PHONY: config
config:
	cp -rf ./config.example.yaml ./config.yaml
	cp -rf ./config.example.yaml ./config.test.yaml

.PHONY: gitlab
gitlab:
	-cp -rf ./-gitlab-ci.yml ./.gitlab-ci.yml
	-rm -rf ./-gitlab-ci.yml