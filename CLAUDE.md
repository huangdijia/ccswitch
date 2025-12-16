# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ccswitch is a Go CLI tool for managing and switching between different Claude Code AI profiles and configurations. It allows users to:

- Switch between different Claude AI profiles/models
- Manage Claude API configurations
- Set environment variables for Claude Code
- Maintain multiple profiles for different use cases

## Build and Development Commands

```bash
# Build the main binary
make build
# or
go build -o ccswitch ./cmd/ccswitch

# Run tests
make test
# or
go test ./...

# Run tests with coverage
make test-coverage

# Install the binary and config files
make install
# or
go install ./cmd/ccswitch

# Format code
make fmt
# or
go fmt ./...

# Run go vet
make vet

# Clean build artifacts
make clean

# Download and tidy dependencies
make mod
```

## High-Level Architecture

The project follows a standard Go project layout with clean architecture principles:

```
ccswitch/
├── cmd/                    # CLI commands
│   ├── ccswitch/          # Main entry point
│   ├── init/              # Initialize configuration
│   ├── list/              # List available profiles
│   ├── set/               # Set specific profile
│   ├── switch/            # Interactive profile switching
│   └── use/               # Apply profile settings
├── internal/              # Private application code
│   ├── claude/            # Claude API client and settings
│   ├── config/            # Configuration management
│   ├── jsonutil/          # JSON utilities
│   └── osutil/            # OS utilities
└── docs/                  # Documentation and schemas
```

### Key Components

1. **Commands (`cmd/`)**: Each subdirectory represents a CLI command implementing Cobra's Commander interface
2. **Configuration (`internal/config/`)**:
   - `Config`: Main configuration structure
   - `ClaudeSettings`: Claude-specific settings with raw JSON support for extensibility
   - Profile and environment variable management
3. **Claude Integration (`internal/claude/`)**:
   - `ClaudeSettings`: Core Claude configuration model
   - `Write()` method persists settings to JSON file
   - Supports model, API key, and custom settings

## Key Patterns and Conventions

1. **Error Handling**: All functions return errors explicitly, use `fmt.Errorf` for error wrapping
2. **JSON Handling**: Use `json.RawMessage` for extensible JSON fields in configurations
3. **File Paths**: Configuration files stored in `~/.config/ccswitch/`
4. **Environment Variables**: Claude settings exported as `CLAUDE_*` environment variables
5. **Testing**: Each command package has corresponding `_test.go` files with unit tests

## Configuration File Location

The main configuration files are located at:

- Main config: `~/.ccswitch/ccs.json` (legacy) or `~/.config/ccswitch/config.json` (new)
- Claude settings: `~/.config/claude/settings.json`

## Important Implementation Details

1. **Profile Management**:
   - Each profile has a name and associated Claude settings
   - Profiles are stored in the main config file
   - The `use` command applies a profile by writing Claude settings

2. **Claude Settings Structure**:
   - Supports standard fields: `model`, `max_tokens`, `temperature`
   - Extensible via `Raw` field for custom JSON properties
   - Handles nested objects like `model_settings` and `tool_choice`

3. **Interactive Features**:
   - The `switch` command provides interactive profile selection
   - Uses terminal UI for user-friendly profile switching

## Testing Guidelines

- Write unit tests for all command packages
- Test error paths and edge cases
- Use table-driven tests for multiple test cases
- Mock external dependencies (file system, API calls) when necessary
- Current test count: 8 test files
- Use `make test-coverage` to generate coverage reports in `coverage.html`

## CI/CD

The project uses GitHub Actions for continuous integration:

- Tests run on every push and pull request with race detection
- Coverage reports are generated automatically
- Release workflow builds binaries for multiple platforms on version tags

## Development Workflow

1. Create new commands in `cmd/` directory
2. Implement Cobra command interface
3. Add corresponding tests in `[command]_test.go`
4. Update configuration structures if needed
5. Test with `make test` before committing
6. Format code with `make fmt`
7. Run static analysis with `make vet`

## Dependencies

Main Go modules include:

- github.com/spf13/cobra - CLI framework
- github.com/spf13/viper - Configuration management
- golang.org/x/term - Terminal utilities

## Pre-configured Profiles

The tool includes several pre-configured profiles for different Claude API providers:

- **default**: Official Anthropic API
- **anyrouter**: AnyRouter proxy service
- **glm**: Zhipu AI's GLM models
- **deepseek**: DeepSeek API
- **kimi-kfc**: Kimi Coding API
- **kimi-k2**: Kimi K2 API

## CCSwitch Codebase Summary

This is a Go-based CLI tool for managing Claude Code API profiles. It follows standard Go project structure with commands in `cmd/` and internal packages in `internal/`. The project uses Cobra for CLI framework, has comprehensive tests, and includes CI/CD via GitHub Actions.

The main binary entry point is `./cmd/ccswitch/main.go` (not `main.go` in root).
