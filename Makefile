
all: generate

generate:
	go run ./cmd/ipcat/main.go

aws:
	go run ./cmd/ipcat/main.go -aws

# todo change to golang
README.md: makestats.py datacenters.csv
	./makestats.py < datacenters.csv > README.md

test:
	find . -name '*.go' | xargs gofmt -w -s
	find . -name '*.go' | xargs goimports -w
	go vet ./...
	golint ./...
	go test .

clean:
	rm -f *~


ci: generate test

docker-ci:
	docker run --rm \
		-e COVERALLS_REPO_TOKEN=$(COVERALLS_REPO_TOKEN) \
		-v $(PWD):/go/src/github.com/client9/ipcat \
		-w /go/src/github.com/client9/ipcat \
		nickg/golang-dev-docker \
		make ci

.PHONY: ci docker-ci
