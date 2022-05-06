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
	GOOS=darwin GOARCH=amd64 go build -o ./bin/await-darwin-amd64 github.com/djui/await
	GOOS=linux  GOARCH=amd64 go build -o ./bin/await-linux-amd64  github.com/djui/await
	GOOS=linux  GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/await-linux-amd64-alpine -installsuffix cgo github.com/djui/await
