.PHONY: all help run test docker-compose-build-api docker-compose-up-api docker-compose-stop-api

RELEASE_VERSION := $(shell git rev-parse --short origin/master)
#RELEASE_VERSION := $(shell git describe) # describe last tag

all: help

## help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "${APP_NAME}":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## setup: get dependencies
setup:
	GO111MODULE=on go mod download

## build: build the application to linux
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build --ldflags="-X 'main.Version=${RELEASE_VERSION}'" -o golang_api_skeleton main.go

## test: run unit tests
test:	
	go test -race -cover -failfast -count=1 ./...

## lint: lints the whole application code
lint:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin v1.31.0
	@golangci-lint run -E golint -e "(.*Sync|.*buf\.Write)"

## docker-compose-build-api: build application docker image
docker-compose-build-api: 
	@docker-compose build

## docker-compose-up-api: up application docker image
docker-compose-up-api: 
	@docker-compose up

## docker-compose-stop-api: stop application docker container
docker-compose-stop-api: 
	@docker-compose stop

## run: run application locally using docker
run: docker-compose-build-api docker-compose-up-api 
