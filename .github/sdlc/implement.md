# Task: implement PROMPT.md

You are running unattended in CI. Make the repository match `PROMPT.md`.

## Rules

- Read `PROMPT.md` first. It is the source of truth. Implement only what it asks for.
- Work in the existing Go module if one exists. If there is no module, create one with a sensible path.
- Use only the Go standard library.
- The CLI entrypoint must be `./cmd/magician` (package `main`) and must declare `var version = "dev"`.
- Add or update focused, fast tests for the behavior you add. No network, sleeps, or external services.
- Update `docs/usage.md` when user-facing behavior changes. Do not edit `README.md`.
- Run `gofmt -w .`, `go mod tidy`, `go build ./...`, `go vet ./...`, and `go test ./...`. Fix every failure.
- Never modify `.github/**`, `.opencode/**`, `PROMPT.md`, `opencode.json`, or `README.md`.
- Never run `git commit`, `git push`, `gh`, or any command that publishes changes. The workflow handles git.

Aim for a small, idiomatic, reviewable diff. When you are done, reply with a short summary of what changed and the test commands you ran.
