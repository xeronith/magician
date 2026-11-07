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

`magician count <text>`

- Print the number of whitespace-separated words in the text.
- Join multiple arguments with a single space.
- Exit with a non-zero status and a clear message when `<text>` is missing.
- Whitespace-only text prints `0`.

`magician serve [addr]`

- Start an HTTP server and block until the process is interrupted.
- Default to `127.0.0.1:8080`; use `addr` when provided.
- Return a clear error and a non-zero exit if the server cannot start.
- Keep every existing command working.

## REST API

All responses that have a body are JSON (`Content-Type: application/json`).

Spells are stored in memory behind a mutex and reset when the process exits. IDs are
increasing integers starting at 1 and are never reused.

| Method | Path | Behavior |
| --- | --- | --- |
| GET | `/health` | `200` with `{"status":"ok"}` |
| GET | `/api/spells` | `200` with a JSON array of spells, `[]` when empty |
| POST | `/api/spells` | `201` with the created spell; invalid JSON or a blank `name` or `effect` returns `400` with `{"error":"..."}` |
| GET | `/api/spells/{id}` | `200` with the spell, or `404` with `{"error":"spell not found"}` |
| DELETE | `/api/spells/{id}` | `204`, or `404` with `{"error":"spell not found"}` |
| GET | `/openapi.json` | `200` with an OpenAPI 3.1 document covering every route and schema |
| GET | `/docs` | `200` with an HTML page that renders the OpenAPI document with Scalar |
| GET | `/` | `302` redirect to `/docs` |

A spell is `{"id": 1, "name": "fireball", "effect": "burns"}`.

A new spell request is `{"name": "fireball", "effect": "burns"}`. Trim both fields and treat
a blank field as invalid.

## Scalar documentation

- `/docs` imports Scalar from the pinned URL
  `https://cdn.jsdelivr.net/npm/@scalar/api-reference@1.77.0/esm.js` and calls
  `createApiReference('#app', { url: '/openapi.json' })`.
- The page needs no build step. The CDN script runs in the browser only; Go code and tests
  must not fetch it.
- Put the OpenAPI document at `cmd/magician/openapi.json` and embed it with `//go:embed`, so
  release binaries work from any directory.

## Requirements

- Go standard library only.
- Entrypoint: `./cmd/magician` with `var version = "dev"` in package `main`.
- Table-driven unit tests for the existing helpers and for every HTTP route and error path.
- HTTP tests must use `net/http/httptest` against the handler. Do not bind a fixed port,
  sleep, or fetch the Scalar CDN in tests.
- `gofmt`, `go build ./...`, `go vet ./...`, `go test ./...`, and `golangci-lint run ./...` must pass.
- A short usage document at `docs/usage.md` with examples for every command, including the
  server, `/docs`, and each REST endpoint.
