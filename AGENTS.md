# Repository Guidelines

## Project Structure & Module Organization

- `main.go` is the entry point for the `ccswitch` CLI.
- `cmd/` contains Cobra commands (`add`, `init`, `list`, `use`, `reset`, `show`, `update`). Each command has a matching `*_test.go` nearby.
- `internal/` holds reusable packages (profiles, settings, path/HTTP utils, terminal UI, output helpers).
- `config/` stores default/preset JSON files (`ccs.json`, `ccs-full.json`, `preset.json`) copied during install.
- `install.sh` provides the release installer; the built binary is `./ccswitch`.

## Build, Test, and Development Commands

- `make build` builds `./ccswitch` with version metadata.
- `make test` or `go test -v ./...` runs the unit tests.
- `make test-coverage` generates `coverage.html`.
- `make fmt` runs `go fmt ./...`; `make vet` runs `go vet ./...`.
- `make install` installs the binary and `config/` defaults; `make clean` removes build artifacts.

## Coding Style & Naming Conventions

- Follow standard Go style and `gofmt`; indentation is tabs (gofmt will enforce).
- Package names are lowercase; exported identifiers use `CamelCase`.
- Test files are named `*_test.go` and live next to the package under test.
- Profile fields use environment-style keys (e.g., `ANTHROPIC_API_KEY`). Keep JSON formatting consistent with existing files.

## Testing Guidelines

- Tests use Go’s `testing` package with `t.Run` subtests; favor table-driven cases for variants.
- Use `t.TempDir()` for filesystem isolation.
- No explicit coverage threshold is enforced; use `make test-coverage` for reports.

## Commit & Pull Request Guidelines

- Recent commit subjects are short, imperative, and often start with verbs like “Add”, “Fix”, “Update”, “Refactor” (some entries are in Chinese). PR numbers may be appended as `(#N)`.
- Keep commits focused; update docs and tests when behavior changes.
- PRs should describe the change, mention tests run, and include screenshots/CLI output when UX changes.

## Configuration & Security Notes

- Runtime config defaults to `~/.ccswitch/ccs.json` and `~/.claude/settings.json` (override via `--profiles`/`--settings`).
- Never commit real API keys; use placeholders in examples/tests.
