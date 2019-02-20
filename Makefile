
all: generate

generate:
	go run ./cmd/ipcat/main.go

aws:
	go run ./cmd/ipcat/main.go -aws

azure:
	go run ./cmd/ipcat/main.go -azure

appengine:
	go run ./cmd/ipcat/main.go -appengine

cloudflare:
	go run ./cmd/ipcat/main.go -cloudflare

install:
	go get golang.org/x/tools/cmd/goimports

test:
	find . -name '*.go' | xargs gofmt -w -s
	find . -name '*.go' | xargs goimports -w
	go vet ./...
	golint ./...
	go test .

misspell:
	misspell README.md
	find . -name '*.go' | misspell

clean:
	rm -f *~
	git gc

ci: generate test

docker-ci:
	docker run --rm \
		-e COVERALLS_REPO_TOKEN=$(COVERALLS_REPO_TOKEN) \
		-v $(PWD):/go/src/github.com/client9/ipcat \
		-w /go/src/github.com/client9/ipcat \
		nickg/golang-dev-docker \
		make ci

.PHONY: ci docker-ci
