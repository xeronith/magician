# Project: magician

Build a small, dependency-free Go CLI called `magician` in the module `github.com/xeronith/magician`.

## Commands

`magician version`

- Print the version. The release pipeline injects it with `-ldflags "-X main.version=<tag>"`.
- The default version when built locally is `dev`.

`magician greet [name]`

- Print a friendly greeting.
- Use `world` when no name is provided.
- Trim surrounding whitespace from the name.
- Treat an empty or whitespace-only name as `world`.

`magician spell <word>`

- Print the word reversed.
- Exit with a non-zero status and a clear message when `<word>` is missing.

## Requirements

- Go standard library only.
- Entrypoint: `./cmd/magician` with `var version = "dev"` in package `main`.
- Table-driven unit tests for the greeting and reversing logic.
- `gofmt`, `go build ./...`, `go vet ./...`, and `go test ./...` must pass.
- A short usage document at `docs/usage.md` with examples for every command.
