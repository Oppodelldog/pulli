setup: ## Install all the build and lint dependencies
	go get -u gopkg.in/alecthomas/gometalinter.v2
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	gometalinter --install --update

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint: ## Run all the linters
	gometalinter --vendor --disable-all \
		--enable=deadcode \
		--enable=gocyclo \
		--enable=ineffassign \
		--enable=gosimple \
		--enable=staticcheck \
		--enable=gofmt \
		--enable=golint \
		--enable=goimports \
		--enable=dupl \
		--enable=misspell \
		--enable=errcheck \
		--enable=vet \
		--enable=vetshadow \
		--enable=varcheck \
		--enable=structcheck \
		--enable=interfacer \
		--enable=goconst \
		--deadline=10m \
		./... | grep -v "mocks"

test: ## Run all the tests
	rm -f coverage.tmp && rm -f coverage.txt
	echo 'mode: atomic' > coverage.txt && go list ./... | xargs -n1 -I{} sh -c 'go test -race -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.txt' && rm coverage.tmp

ci: ## Run all the tests and code checks
	dep ensure
	go get ./...
	make test
	# goveralls -service drone.io -repotoken Hnq7byXbwVH2lKmNnrV0svSn8P6UOV9vZ
	go build src/cmd/main.go

build: ## build binary to .build folder
	rm -f .build/*
	go build -o ".build/pulli" src/cmd/main.go

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