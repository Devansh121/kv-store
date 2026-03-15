# kv-store

A simple Redis-compatible asynchronous in-memory key-value store.

## Local checks

Before pushing to `origin`, run the same core checks used by CI:

```bash
gofmt -w .
golangci-lint run
go vet ./...
go build ./...
go test ./...
```

If `golangci-lint` is not installed locally:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6
```
