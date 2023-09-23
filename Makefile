.PHONY: all
all: build

.PHONY: build
build:
	go build -v

.PHONY: lint
lint:
	docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:v1.45.2 golangci-lint run -v

.PHONY: test
test:
	go test -v

.PHONY: rel
rel:
	GOOS=darwin GOARCH=amd64 go build -o await-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o await-darwin-arm64
	GOOS=linux  GOARCH=amd64 go build -o await-linux-amd64
	GOOS=linux  GOARCH=arm64 go build -o await-linux-arm64
