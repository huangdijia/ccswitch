# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Core Principles (Priority Order)

1. **Repository Standards First**: Existing repository conventions take precedence. When conflicts arise, follow `README`, `AGENTS.md`, or other established documentation and record decisions.
2. **Documentation as Source of Truth**: Requirements, interactions, and interfaces must be documented. Changes should be traceable to specific tasks or issues.
3. **Atomic Tasks**: Handle one atomic task at a time. All changes must be traceable and include verification methods (tests or manual steps).
4. **Documentation Feedback Loop**: After implementation, update relevant documentation including task status, change logs, and specifications when needed.

## Project Overview

ccswitch is a Go CLI tool for managing and switching between different Claude Code AI profiles and configurations. It allows users to:

- Switch between different Claude AI profiles/models
- Manage Claude API configurations
- Set environment variables for Claude Code
- Maintain multiple profiles for different use cases

## Build, Test, and Development Commands

Use the repository's provided build tools. All commands should be run from the project root.

### Build
```bash
# Build the main binary (outputs ./ccswitch with version metadata)
make build

# Build without make
go build -ldflags="-X main.version=... -s -w" -o ccswitch .
```

### Test
```bash
# Run all tests
make test
# or
go test -v ./...

# Run tests with coverage report (generates coverage.html)
make test-coverage
```

### Code Quality
```bash
# Format code (required before commits)
make fmt
# or
go fmt ./...

# Run static analysis
make vet
```

### Installation & Cleanup
```bash
# Install binary and config files to system
make install

# Remove build artifacts
make clean

# Update dependencies
make mod
```

### Available Targets
Run `make help` to see all available targets and their descriptions.

## Repository Structure & Documentation

### Code Organization

The project follows a standard Go project layout with clean architecture principles:

```
ccswitch/
├── main.go                 # Entry point for the CLI
├── cmd/                    # Cobra commands (add, init, list, use, reset, show, update)
│   └── *_test.go          # Unit tests for each command
├── internal/              # Private application code
│   ├── profiles/          # Profile management
│   ├── settings/          # Claude settings management
│   ├── pathutil/          # Path utilities
│   ├── httputil/          # HTTP utilities
│   ├── output/            # Output formatting helpers
│   └── cmdutil/           # Command utility functions
├── config/                # Default/preset JSON files
│   ├── ccs.json           # Basic configuration
│   ├── ccs-full.json      # Complete configuration with all providers
│   └── preset.json        # Preset profiles for install command
└── install.sh             # Release installer script
```

**Key Organization Principles:**
- `main.go` is the entry point for the `ccswitch` CLI
- `cmd/` contains Cobra commands, each with matching `*_test.go` test files
- `internal/` holds reusable packages (not importable by external projects)
- `config/` stores default/preset JSON files copied during install
- The built binary is `./ccswitch`

### Documentation Structure

- **README.md**: User-facing documentation (installation, usage, examples)
- **README_CN.md**: Chinese version of user documentation
- **AGENTS.md**: Repository-specific development guidelines and workflows
- **CLAUDE.md**: This file - guidance for Claude Code AI agent
- **GEMINI.md**: Guidance for Gemini Code AI agent

### Configuration Files

Runtime configuration defaults:
- Main config: `~/.ccswitch/ccs.json` (override via `--profiles` or `CCSWITCH_PROFILES_PATH`)
- Claude settings: `~/.claude/settings.json` (override via `--settings` or `CCSWITCH_SETTINGS_PATH`)

## Development Workflow & Task Management

### Task and Issue Tracking

When working on tasks or fixing issues:

1. **Task Identification**: Each task should be atomic, completable in one work session, with observable output and independent verification
2. **Implementation**: Make focused, minimal changes that address the specific task
3. **Verification**: 
   - Add or update tests for code changes
   - Run `make test` to verify existing functionality isn't broken
   - For UI/CLI changes, provide manual verification steps
4. **Documentation**: Update relevant documentation (README, AGENTS.md, CLAUDE.md) when behavior changes

### Issue Management

When encountering bugs or issues:

1. **Reproduce**: Document environment, reproduction steps, expected vs actual behavior
2. **Investigate**: Identify root cause and affected components
3. **Fix**: Implement minimal, focused fix
4. **Verify**: Add regression tests, verify fix resolves the issue
5. **Document**: Record resolution details and verification results

See AGENTS.md for detailed issue tracking guidelines.

## Coding Standards & Best Practices

### Go Style Conventions

- **Formatting**: Use `gofmt` (tabs for indentation, enforced by formatter)
- **Naming**: 
  - Package names: lowercase, single word
  - Exported identifiers: `CamelCase`
  - Unexported identifiers: `camelCase`
  - Constants: `CamelCase` or `UPPER_SNAKE_CASE` for package-level constants
- **Line Length**: Aim for ≤120 characters when reasonable
- **Test Files**: Named `*_test.go`, live next to the package under test

### Project-Specific Conventions

- **Profile Fields**: Use environment-style keys (e.g., `ANTHROPIC_API_KEY`)
- **JSON Formatting**: Keep consistent with existing config files
- **Error Handling**: All functions return errors explicitly, use `fmt.Errorf` for error wrapping
- **JSON Handling**: Use `json.RawMessage` for extensible JSON fields
- **Environment Variables**: Claude settings use `ANTHROPIC_*` prefix

### Code Quality Requirements

1. **Minimal Changes**: Only modify what's necessary to accomplish the task
2. **No Mass Refactoring**: Unless the task specifically requires it, avoid reformatting or reorganizing existing code
3. **Tests Required**: Add tests for new features, update tests for modified behavior
4. **Table-Driven Tests**: Use `t.Run` subtests with table-driven patterns for multiple test cases
5. **Filesystem Isolation**: Use `t.TempDir()` for test isolation
6. **Error Testing**: Test error paths and edge cases

## Testing Guidelines

### Test Structure

- **Location**: Test files live next to the package under test, named `*_test.go`
- **Framework**: Use Go's built-in `testing` package
- **Organization**: Use `t.Run` for subtests, favor table-driven tests for variants
- **Isolation**: Use `t.TempDir()` for filesystem isolation in tests
- **Coverage**: No explicit threshold enforced, but use `make test-coverage` to generate reports

### Testing Workflow

```bash
# Run all tests
make test

# Generate coverage report (creates coverage.html)
make test-coverage

# Run specific tests
go test -v ./cmd/... -run TestSpecificFunction
```

### Test Requirements

- Test error paths and edge cases
- Mock external dependencies (filesystem, network calls) when necessary
- Ensure tests are deterministic and don't rely on external state
- Write meaningful test names that describe what's being tested
- Keep tests focused and independent

### Current Test Status

- 8 test files across command packages
- Each command has corresponding unit tests
- Coverage reports available via `make test-coverage`

## Key Implementation Details

### Profile Management

- **Storage**: Profiles stored in `~/.ccswitch/ccs.json`
- **Structure**: Each profile contains environment variables (keys like `ANTHROPIC_API_KEY`, `ANTHROPIC_BASE_URL`, etc.)
- **Application**: The `use` command applies a profile by writing its environment variables to Claude settings
- **Descriptions**: Optional descriptions stored in the `descriptions` field of the config

### Claude Settings Integration

- **Target File**: `~/.claude/settings.json` (configurable)
- **Format**: JSON format with environment variable mappings
- **Fields**: Standard fields include `ANTHROPIC_API_KEY`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_MODEL`, plus model variants (`ANTHROPIC_DEFAULT_HAIKU_MODEL`, etc.)
- **Extensibility**: Settings can include custom fields beyond the standard set

### Interactive Features

- **Interactive Selection**: Commands like `use` provide keyboard-driven selection (↑/↓ and Enter) when run without arguments
- **Terminal UI**: Uses `golang.org/x/term` for terminal utilities
- **User Feedback**: Masked sensitive values in output for security

## Commit & Pull Request Guidelines

### Commit Standards

- **Format**: Use imperative mood, short subject line (e.g., "Add feature", "Fix bug", "Update docs")
- **Convention**: Follow Conventional Commits when applicable (`feat:`, `fix:`, `docs:`, `test:`, `chore:`)
- **Scope**: One commit should focus on a single logical change or task
- **Mixed Languages**: Recent commits show both English and Chinese subjects are acceptable
- **PR References**: May include PR numbers as `(#N)` in commit message

### Pull Request Requirements

When submitting PRs:

1. **Description**: 
   - Explain the motivation and context
   - List behavior changes
   - Describe verification methods (tests run, manual testing performed)
   - Include risk assessment and rollback plan if applicable

2. **Testing**:
   - Run `make test` to ensure all tests pass
   - Run `make fmt` and `make vet` before committing
   - Add new tests for new features or bug fixes
   - Update existing tests if behavior changes

3. **Documentation**:
   - Update README.md if user-facing behavior changes
   - Update AGENTS.md if development workflow changes
   - Include CLI output examples for command changes
   - Attach screenshots/GIFs for UI changes

4. **Code Review**:
   - Keep changes focused and reviewable
   - Respond to feedback constructively
   - Update documentation based on review comments

## Security Considerations

### API Key Management

- **Never commit real API keys**: Use placeholders in examples and tests (e.g., `sk-...`, `****`)
- **Config File Security**: Configuration files contain API keys in plain text
  - Ensure proper file permissions (readable only by user)
  - Never commit config files to version control
  - Add config paths to `.gitignore`

### Code Security

- **Input Validation**: Validate all user inputs, especially file paths and URLs
- **Path Traversal Protection**: Sanitize paths to prevent directory traversal attacks
- **Dependency Security**: Keep dependencies updated, review security advisories
- **Network Security**: Verify checksums when downloading files (e.g., in update command)

### Best Practices

- Use environment variables for sensitive configuration when possible
- Log security-relevant events (authentication, configuration changes)
- Provide clear error messages without exposing sensitive information
- Document security considerations in code comments where relevant

## CI/CD Pipeline

The project uses GitHub Actions for continuous integration and deployment:

### Test Workflow
- **Trigger**: Runs on every push and pull request
- **Execution**: 
  - Runs all unit tests with race detection (`go test -race`)
  - Generates coverage reports
  - Builds the binary to verify compilation
- **Requirements**: All tests must pass before merging

### Release Workflow
- **Trigger**: Runs on version tags (e.g., `v1.0.0`)
- **Process**:
  - Runs full test suite
  - Builds binaries for multiple platforms (Linux, macOS, Windows)
  - Supports multiple architectures (amd64, arm64, arm)
  - Creates GitHub releases with binaries and checksums
- **Artifacts**: Cross-platform binaries with SHA256 checksums

## Dependencies

### Core Dependencies

- **github.com/spf13/cobra**: CLI framework for command structure
- **github.com/spf13/pflag**: POSIX/GNU-style command-line flags (used by Cobra)
- **golang.org/x/term**: Terminal utilities for interactive features

### Go Version

- **Minimum**: Go 1.21 or higher
- **Target**: Specified in `go.mod`

### Dependency Management

```bash
# Update dependencies
make mod

# Or manually
go mod download
go mod tidy
```

## Pre-configured Profiles

The tool includes several pre-configured profiles for different Claude API providers:

- **default**: Official Anthropic API
- **anyrouter**: AnyRouter proxy service
- **glm**: Zhipu AI's GLM models
- **deepseek**: DeepSeek API
- **kimi-kfc**: Kimi Coding API
- **kimi-k2**: Kimi K2 API
- **modelscope**: ModelScope API
- **minimaxi-m2**: MiniMax's Anthropic-compatible API
- **xiaomi-mimo**: Xiaomi Mimo's Anthropic-compatible API

These profiles are available in:
- `config/ccs-full.json`: Complete configuration with all providers
- `config/preset.json`: Preset configuration for the `install` command

## Common Development Tasks

### Adding a New Command

1. Create command file in `cmd/` directory (e.g., `cmd/newcommand.go`)
2. Implement Cobra command interface:
   ```go
   var newcommandCmd = &cobra.Command{
       Use:   "newcommand",
       Short: "Short description",
       Long:  "Long description",
       RunE:  runNewCommand,
   }
   ```
3. Add command to root in `cmd/root.go`
4. Create corresponding test file `cmd/newcommand_test.go`
5. Update README.md with usage examples
6. Run tests: `make test`

### Modifying Configuration Structure

1. Update structures in `internal/profiles/` or `internal/settings/`
2. Ensure backward compatibility with existing config files
3. Update default config files in `config/`
4. Add migration code if needed
5. Update tests to cover new fields
6. Document changes in README.md

### Adding a New Profile

1. Edit `config/ccs-full.json` to add the profile
2. Update `config/preset.json` if it should be available via `install`
3. Test with: `ccswitch init --full && ccswitch list`
4. Verify profile can be activated: `ccswitch use <profile-name>`
5. Update README.md pre-configured profiles section

## Troubleshooting Common Issues

### Build Issues

```bash
# Clean and rebuild
make clean
make build

# Update dependencies
make mod

# Check Go version
go version  # Should be 1.21+
```

### Test Failures

```bash
# Run specific test
go test -v ./cmd/... -run TestSpecificFunction

# Run with race detection
go test -race ./...

# Check test coverage
make test-coverage
# Open coverage.html in browser
```

### Configuration Issues

```bash
# Verify config file location
ls -la ~/.ccswitch/ccs.json

# Verify Claude settings location  
ls -la ~/.claude/settings.json

# Reinitialize configuration
ccswitch init --force
```

## Quick Reference

### File Locations
- Binary: `./ccswitch`
- Config: `~/.ccswitch/ccs.json`
- Claude Settings: `~/.claude/settings.json`
- Default Configs: `config/*.json`

### Essential Commands
```bash
make build          # Build binary
make test           # Run tests
make test-coverage  # Generate coverage report
make fmt            # Format code
make vet            # Static analysis
make clean          # Clean artifacts
```

### Key Packages
- `cmd/`: CLI commands
- `internal/profiles/`: Profile management
- `internal/settings/`: Claude settings handling
- `internal/pathutil/`: Path utilities
- `internal/output/`: Output formatting

---

## CCSwitch Summary

ccswitch is a Go-based CLI tool for managing Claude Code API profiles. It follows standard Go project structure with commands in `cmd/` and internal packages in `internal/`. The project uses Cobra for the CLI framework, has comprehensive tests, and includes CI/CD via GitHub Actions.

**Entry Point**: `main.go` (not in a subdirectory)
**Binary Output**: `./ccswitch`

For detailed user documentation, see [README.md](README.md) or [README_CN.md](README_CN.md).
For development workflow guidelines, see [AGENTS.md](AGENTS.md).
