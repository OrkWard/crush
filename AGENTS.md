# Crush Development Guide

## Build/Test/Lint Commands

Use `task --list` for all available tasks.

- **Build**: `task build` or `go build .`
- **Run**: `task run` or `go run .`
- **Test**: `task test` (runs `go test -race -failfast ./...`)
  - Single test: `go test ./internal/ui/diffview -run TestDiffView`
- **Update Golden Files**: `go test ./... -update`
  - Specific package: `go test ./internal/ui/diffview -update`
  - Re-record agent VCR: `task test:record`
- **Lint**: `task lint` (includes `lint:log` + `golangci-lint`)
  - Install linter: `task lint:install`
  - Fix: `task lint:fix`
- **Format**: `task fmt` (`gofumpt -w .`) + `task fmt:html`
- **Modernize**: `task modernize`
- **Dev**: `task dev` (profiling enabled)
- **Install**: `task install`
- **Schema**: `task schema` (generates `schema.json`)
- **Hyper Providers**: `task hyper` (`go generate ./internal/agent/hyper/...`)
- **Release**: `task release` (on main, clean git, CI passed)

## Linting

Configured in `.golangci.yml` (v2):
- Enabled: bodyclose, goprintffuncname, misspell, noctx, nolintlint, rowserrcheck, sqlclosecheck, staticcheck, tparallel, whitespace
- Formatters: gofumpt, goimports
- Log messages must start with capital (enforced by `lint:log`)

## Code Style Guidelines

- **Imports**: `goimports` grouped (stdlib, external, internal)
- **Formatting**: `gofumpt` (stricter gofmt)
- **Naming**: PascalCase exported, camelCase unexported
- **Types**: Explicit types, aliases for clarity (`type AgentName string`)
- **Errors**: `fmt.Errorf` wrapping
- **Context**: First param in ops
- **Interfaces**: In consumer pkg, small/focused
- **Structs**: Embed for composition
- **Constants**: Typed iota enums
- **JSON tags**: snake_case
- **Permissions**: Octal (0o755, 0o644)
- **Comments**: Capital start, period end (78 cols)
- **Logs**: Capital first letter

## Testing

- `testify/require`, `t.Parallel()`, `t.Setenv()`, `t.TempDir()` (no cleanup)
- Golden files: `./internal/ui/diffview/testdata/**/*.golden`
- Mock providers:
  ```go
  original := config.UseMockProviders
  config.UseMockProviders = true
  defer func() { config.UseMockProviders = original; config.ResetProviders() }()
  config.ResetProviders()
  ```

## Database

- **Schema**: `internal/db/migrations/*.sql` (goose)
- **Queries**: `internal/db/sql/*.sql` -> `internal/db/*.sql.go` (sqlc v1.30.0)
- **Regenerate**: `sqlc generate`
- Migrations auto-applied on `db.Connect()` via `goose.Up`

## UI Development

Read `internal/ui/AGENTS.md` before TUI work:
- No IO/expensive in `Update`: use `tea.Cmd`
- Components dumb: `Render(width) string`
- Use `charmbracelet/x/ansi` for ANSI
- etc.

## Committing

- Semantic prefixes: `feat:`, `fix:`, `chore:`, etc.
- One line + attribution

## Other

- Go: 1.26.0 (GOTOOLCHAIN: go1.25.0 for lint)
- Profiling: `task profile:cpu|heap|allocs`
- Deps: `task deps`
