# Image URL to use all building/pushing image targets
IMG ?= solidbroccoli:latest

all: build

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: fmt vet ## Build program.
	go build -a -o dist/broccoli main.go

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

docker-build: ## Build docker image with the manager.
	docker build -t ${IMG} .
