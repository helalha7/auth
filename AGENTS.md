# AGENTS.md

Guidance for coding agents working in this repository.

## Project overview

This is a small Go authentication API backed by MySQL.

- `cmd/`: application entry point, configuration, database setup, and routes
- `internal/users/`: user model, password hashing, persistence, and registration
- `internal/tokens/`: authentication-token generation and persistence
- `internal/middlewares/`: bearer-token authentication middleware
- `internal/httpx/`: shared HTTP/JSON helpers
- `migrations/`: ordered MySQL migrations using Goose annotations

The server listens on `:8080`. It currently exposes:

- `POST /v1/users`
- `POST /v1/tokens/authentication`
- `GET /v1/test` (requires authentication)

## Development commands

Run these from the repository root:

```sh
go test ./...
go vet ./...
gofmt -w <changed-go-files>
go run ./cmd
```

The application requires `DB_DSN` in the root `.env` file. Never print, edit, or commit `.env` or other credentials. A running MySQL database with the migrations applied is required to start the server.

## Change guidelines

- Keep code within the existing package boundaries; wire dependencies and routes in `cmd/main.go`.
- Use the request context for all database calls.
- Keep SQL parameterized with `?` placeholders. Never interpolate request data into SQL.
- Define narrow store interfaces beside the handlers or middleware that consume them. This supports isolated tests with fakes.
- Return stable package-level sentinel errors for expected persistence outcomes, and translate them into HTTP responses in handlers.
- Use `internal/httpx.WriteJSON` and envelope-style JSON objects for JSON responses.
- Never expose password hashes, token hashes, or other secrets in API responses or logs.
- Store only token hashes in the database; return plaintext tokens only when they are initially created.
- Preserve the existing URL version prefix (`/v1`) when adding endpoints.
- Add a new sequential migration for schema changes; do not rewrite an existing migration that may already have been applied.
- Keep migrations reversible with both `-- +goose Up` and `-- +goose Down` sections.

## Testing expectations

- Add or update tests for behavior changes. Prefer table-driven tests where appropriate.
- Handler tests should use `httptest` and fake implementations of the local store interfaces.
- Store tests that require MySQL should be clearly separated from unit tests and must not rely on a developer's `.env` implicitly.
- Before finishing, run `gofmt` on changed Go files, then `go test ./...` and `go vet ./...`.
- If a check cannot run because MySQL or another external dependency is unavailable, report that explicitly.

## Scope and safety

- Make focused changes and avoid unrelated refactors.
- Do not discard or overwrite existing uncommitted work.
- Do not add generated binaries, local database data, secrets, or environment files to version control.
- Maintain backward compatibility for existing request and response fields unless the task explicitly requires an API change.
