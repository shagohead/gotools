gotools
=======

Tiny (without dependencies) `go run package@version` alternative that pre-installs binaries into (project) local directory. This eliminates waiting of tools building on every `go generate` call (gotools itself builds relatively fast).

Why not use «tools.go» approach which tracks apps versions in go.mod? Because in that way tools dependencies becomes project dependencies and affects project's dependency graph and vice versa: your project dependencies affects tools dependencies, and even more: tools affects each other dependencies. This is wrong, because and your project and each tool should have their own dependencies isolated from others.

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
