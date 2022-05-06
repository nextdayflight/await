.PHONY: all
all: build

.PHONY: build
build:
	go build -v github.com/djui/await

.PHONY: lint
lint:
	@if [ $$(gofmt -l . | wc -l) != 0 ]; then \
	    echo "gofmt: code not formatted"; \
	    gofmt -l . | grep -v vendor/; \
	    exit 1; \
	fi

	@gometalinter \
	             --vendor \
	             --tests \
	             --disable=gocyclo \
	             --disable=dupl \
	             --disable=deadcode \
	             --disable=gotype \
	             --disable=maligned \
	             --disable=interfacer \
	             --disable=varcheck \
	             --disable=gosec \
	             --disable=megacheck \
	             ./...

.PHONY: test
test:
	go test -v

.PHONY: rel
rel:
	GOOS=darwin GOARCH=amd64 go build -o await-darwin-amd64 github.com/betalo-sweden/await
	GOOS=linux  GOARCH=amd64 go build -o await-linux-amd64  github.com/betalo-sweden/await
