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

`magician shout <text>`

- Print the text uppercased, followed by `!!!`.
- Join multiple words with a single space.
- Exit with a non-zero status and a clear message when `<text>` is missing.

## Requirements

- Go standard library only.
- Entrypoint: `./cmd/magician` with `var version = "dev"` in package `main`.
- Table-driven unit tests for the greeting, reversing, and shouting logic.
- `gofmt`, `go build ./...`, `go vet ./...`, `go test ./...`, and `golangci-lint run ./...` must pass.
- A short usage document at `docs/usage.md` with examples for every command.
