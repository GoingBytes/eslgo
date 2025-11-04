fmt:
    gofumpt -l -w .

lint:
    golangci-lint run --timeout 5m --tests=false

test:
    go test ./...
