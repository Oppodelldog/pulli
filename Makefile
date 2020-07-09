setup: ## Install tools
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s v1.27.0
	mkdir .bin || true; mv bin/golangci-lint .bin/golangci-lint || true

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint: ## Run the linters
	golangci-lint run

test: ## Run all the tests
	go version
	go env
	go list ./... | xargs -n1 -I{} sh -c 'go test -race {}'

build: ## build binary to .build folder
	rm -f .build/*
	go build -o ".build/pulli" cmd/main.go

install: ## Install to <gopath>/src
	go install ./...

build-release: ## builds the checked out version into the .release/${tag} folder
	.release/build.sh

build-release-test: ## builds the checked out version into the .release/${tag} folder
	.release/build.sh test

# Self-Documented Makefile see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help