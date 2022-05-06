.PHONY: all
all: build

.PHONY: build
build:
	go build -v github.com/djui/await

.PHONY: lint
lint:
	docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:v1.45.2 golangci-lint run -v

.PHONY: test
test:
	go test -v

.PHONY: rel
rel:
	GOOS=darwin GOARCH=amd64 go build -o await-darwin-amd64 github.com/betalo-sweden/await
	GOOS=linux  GOARCH=amd64 go build -o await-linux-amd64  github.com/betalo-sweden/await
