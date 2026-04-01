# Contributing to speedy-cli

Thanks for your interest in contributing.

## Development Setup

1. Install Go 1.22+
2. Clone repo
3. Run:
   - `go mod tidy`
   - `go test ./...`
   - `go build -o speedy-cli .`

## Pull Requests

- Keep PRs focused and small.
- Include tests for new behavior when possible.
- Update docs for user-facing changes.
- Follow conventional commit style (e.g. `feat:`, `fix:`, `docs:`).

## Code Style

- Run `gofmt` on changed files.
- Keep packages modular under `cmd/` and `internal/`.
- Avoid hardcoding secrets and credentials.
