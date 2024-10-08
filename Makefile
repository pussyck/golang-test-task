SHELL = /bin/sh
LDFLAGS = -s -w

DOCKER_BIN = $(shell command -v docker 2> /dev/null)
DC_BIN = $(shell command -v docker-compose 2> /dev/null)
DC_RUN_ARGS = --rm --user "$(shell id -u):$(shell id -g)"
APP_NAME = $(notdir $(CURDIR))

.PHONY : help build fmt lint gotest test cover shell redis-cli image clean
.DEFAULT_GOAL : help
.SILENT : test shell redis-cli

# Это выведет справку по каждой задаче.
help: ## Show this help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[32m%-11s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build app binary file
	$(DC_BIN) run $(DC_RUN_ARGS) --no-deps app go build -ldflags="$(LDFLAGS)" .

fmt: ## Run source code formatter tools
	$(DC_BIN) run $(DC_RUN_ARGS) --no-deps app sh -c 'GO111MODULE=off go get golang.org/x/tools/cmd/goimports && $$GOPATH/bin/goimports -d -w .'
	$(DC_BIN) run $(DC_RUN_ARGS) --no-deps app gofmt -s -w -d .

lint: ## Run app linters
	$(DOCKER_BIN) run --rm -t -v "$(shell pwd):/app" -w /app golangci/golangci-lint:v1.24-alpine golangci-lint run -v

gotest: ## Run app tests
	$(DC_BIN) run $(DC_RUN_ARGS) app go test -v -race ./...

test: lint gotest ## Run app tests and linters
	@printf "\n   \e[30;42m %s \033[0m\n\n" 'All tests passed!';

cover: ## Run app tests with coverage report
	$(DC_BIN) run $(DC_RUN_ARGS) app sh -c 'go test -race -covermode=atomic -coverprofile /tmp/cp.out ./... && go tool cover -html=/tmp/cp.out -o ./coverage.html'
	-sensible-browser ./coverage.html && sleep 2 && rm -f ./coverage.html

shell: ## Start shell into container with golang
	$(DC_BIN) run $(DC_RUN_ARGS) app sh

redis-cli: ## Start redis-cli
	$(DC_BIN) run --rm redis redis-cli -h redis -p 6379 \
		|| printf "\n   \e[1;41m %s \033[0m\n\n" "Probably you need to run \`docker-compose up -d\` before.";

image: ## Build docker image with app
	$(DOCKER_BIN) build -f ./Dockerfile -t $(APP_NAME):local .
	@printf "\n   \e[30;42m %s \033[0m\n\n" 'Now you can use image like `docker run --rm $(APP_NAME):local ...`';

clean: ## Make clean
	$(DC_BIN) down -v -t 1
	$(DOCKER_BIN) rmi $(APP_NAME):local -f
