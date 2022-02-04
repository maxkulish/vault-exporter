DOCKER_IMAGE_NAME ?= grapeshot/vault_exporter
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR = $(ROOT_DIR)/bin
PROJ_NAME = vault-exporter

help: __help__

__help__:
	@echo make build - build go executables in the ./bin folder
	@echo make clean - delete executables, download project from github and build
	@echo make coverage - run test coverage and open html file with results
	@echo make benchmark - run benchmark tests with memory alocation
	@echo make test - run all tests in a project

vault_exporter: main.go
	go build -o $@ ./

.uptodate: vault_exporter Dockerfile
	docker build -t $(DOCKER_IMAGE_NAME):$(shell git rev-parse --short HEAD) .

lint:
	gometalinter --enable-all --vendor --deadline=5m

update-vendor:
	dep ensure
	dep prune

build: clean
	make build_mac
	make build_linux

build_mac:
	cd $(ROOT_DIR)
	GOOS=darwin GOARCH=amd64 go build --race -o $(BIN_DIR)/macos/$(PROJ_NAME) main.go

build_linux:
	cd $(ROOT_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/linux/$(PROJ_NAME) main.go

clean:
	rm -rf ./bin .uptodate vault_exporter

test:
	cd $(ROOT_DIR)
	go test -v ./...

benchmark:
	cd $(ROOT_DIR)
	go test -bench . -benchmem

coverage:
	cd $(ROOT_DIR)
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

zip_linux:
	cd $(ROOT_DIR)
	zip $(PROJ_NAME)_linux.zip $(BIN_DIR)/linux/$(PROJ_NAME)

zip_mac:
	cd $(ROOT_DIR)
	zip $(PROJ_NAME)_mac.zip $(BIN_DIR)/macos/$(PROJ_NAME)