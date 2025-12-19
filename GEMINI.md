# GEMINI.md

## Project Overview

This is a Go command-line tool named `ccswitch` for managing and switching between different Claude Code API profiles and configurations. It's built using the Cobra CLI framework.

The tool allows users to manage multiple Claude Code API configurations (profiles) and switch between them. This is useful when you need to use different API endpoints, models, or authentication tokens for different projects or environments.

The configuration is stored in `~/.ccswitch/ccs.json`.

## Building and Running

The project uses a `Makefile` for common development tasks.

### Building

To build the project, run:

```bash
make build
```

This will create a binary named `ccswitch` in the root directory.

### Running Tests

To run the tests, use:

```bash
make test
```

### Running the application

After building, you can run the application directly:

```bash
./ccswitch --help
```

## Development Conventions

### Code Style

The project follows standard Go conventions. Use `gofmt` to format the code before committing.

```bash
make fmt
```

### Linting

The project uses `go vet` for static analysis.

```bash
make vet
```

### Dependencies

Dependencies are managed with Go modules. To download and tidy dependencies, run:

```bash
make mod
```

### Commits

The project does not seem to enforce any strict commit message conventions, but the commit history shows a preference for conventional commit messages.
