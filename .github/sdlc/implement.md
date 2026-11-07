# Task: implement PROMPT.md

You are running unattended in CI. Make the repository match `PROMPT.md`.

## Rules

- Read `PROMPT.md` first. It is the source of truth. Implement only what it asks for.
- Work in the existing Go module if one exists. If there is no module, create one with a sensible path.
- Use only the Go standard library.
- The CLI entrypoint must be `./cmd/magician` (package `main`) and must declare `var version = "dev"`. Keep every existing command working.
- Add or update focused, fast tests for the behavior you add. No external network, sleeps, or external services. HTTP handler tests may use `net/http/httptest`; never bind a fixed port.
- Embed any runtime data such as `cmd/magician/openapi.json` with `//go:embed`; release binaries must not read files from disk at runtime.
- Update `docs/usage.md` when user-facing behavior changes. Do not edit `README.md`.
- Handle every returned error or explicitly discard it; `errcheck` must report zero findings.
- Run `gofmt -w .`, `go mod tidy`, `go build ./...`, `go vet ./...`, `go test ./...`, and `golangci-lint run ./...`. Fix every failure and lint finding.
- Never modify `.github/**`, `.opencode/**`, `PROMPT.md`, `opencode.json`, or `README.md`.
- Never run `git commit`, `git push`, `gh`, or any command that publishes changes. The workflow handles git.

Aim for a small, idiomatic, reviewable diff. When you are done, reply with a short summary of what changed and the test commands you ran.
