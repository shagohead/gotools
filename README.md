gotools
=======

Tiny (without dependencies) `go run package@version` alternative that pre-installs binaries into (project) local directory. This eliminates waiting of tools building on every `go generate` call (gotools itself builds relatively fast).

Usage
-----

When gotools invoked it searches for `go.tools` file in current working directory or in one of it parents. It makes sense to store `go.tools` in the same directory as `go.mod` and `go.sum`. `go.tools` content example:

```
github.com/ogen-go/ogen/cmd/ogen@v1.2.1
github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0
```

Use these tools with `gotools` like this:

```go
//go:generate go run github.com/shagohead/gotools@v0.1.1 ogen -target=client -package=client -clean openapi.yaml
//go:generate go run github.com/shagohead/gotools@v0.1.1 sqlc generate
```
