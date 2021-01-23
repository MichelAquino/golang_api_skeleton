.PHONY: all help run test docker-compose-build-api docker-compose-up-api docker-compose-stop-api

all: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup:
	GO111MODULE=on go mod download

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o golang_api_skeleton main.go

test:	## Start unit tests
	go test -race -cover -failfast ./...

docker-compose-build-api: ## Build application's image
	@docker-compose build

docker-compose-up-api: ## Up application's image
	@docker-compose up

docker-compose-stop-api: ## Stop application's containers
	@docker-compose stop

run: docker-compose-build-api docker-compose-up-api ## Run application locally
