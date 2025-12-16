# CCSwitch

[![CI](https://github.com/huangdijia/ccswitch/workflows/CI/badge.svg)](https://github.com/huangdijia/ccswitch/actions)
[![Release](https://img.shields.io/github/release/huangdijia/ccswitch.svg)](https://github.com/huangdijia/ccswitch/releases)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

English | [中文文档](README_CN.md)

A powerful command-line tool for managing and switching between different Claude Code API profiles and configurations.

## Description

CCSwitch allows you to easily manage multiple Claude Code API configurations (profiles) and switch between them. This is useful when you need to use different API endpoints, models, or authentication tokens for different projects or environments.

### Key Features

- **Multi-profile Management**: Store and switch between multiple Claude API configurations
- **Auto-update**: Built-in update command to keep the tool current
- **Pre-configured Providers**: Support for various Claude API providers out of the box
- **Cross-platform**: Works on Linux, macOS, and Windows
- **Simple CLI**: Intuitive commands for easy profile management
- **Configuration Persistence**: Your settings are safely stored and applied automatically

## Installation

### Quick Install (Recommended)

The easiest way to install CCSwitch is with our installation script:

```bash
curl -sSL https://raw.githubusercontent.com/huangdijia/ccswitch/main/install.sh | bash
```

This will:

- Detect your platform and architecture
- Download the latest release binary
- Install it to `~/.local/bin`
- Add helpful instructions if the directory isn't in your PATH

**Installation options:**

```bash
# Install to a custom directory
curl -sSL https://raw.githubusercontent.com/huangdijia/ccswitch/main/install.sh | bash -s -- -d /usr/local/bin

# Install a specific version
curl -sSL https://raw.githubusercontent.com/huangdijia/ccswitch/main/install.sh | bash -s -- -v v1.0.0
```

### Using Go Install

If you have Go 1.21 or higher installed, you can install CCSwitch directly:

```bash
go install github.com/huangdijia/ccswitch@latest
```

Make sure your Go bin directory is in your PATH:

```bash
export PATH="$HOME/go/bin:$PATH"
```


### From Source

1. Clone the repository:

```bash
git clone https://github.com/huangdijia/ccswitch.git
cd ccswitch
```

2. Build and install:

```bash
make install
```

Or build only:

```bash
make build
./ccswitch init
```

### Binary Releases

Download pre-built binaries from the [releases page](https://github.com/huangdijia/ccswitch/releases).

## Quick Start

```bash
# Install ccswitch (if you haven't already)
curl -sSL https://raw.githubusercontent.com/huangdijia/ccswitch/main/install.sh | bash

# Initialize your configuration
ccswitch init

# List available profiles
ccswitch profiles

# Switch to a profile
ccswitch use glm
```

## Usage

### Initialize configuration

```bash
ccswitch init
```

This command initializes the CCSwitch configuration. It will:

- Create a `~/.ccswitch/` directory if it doesn't exist
- Copy the default profile configuration to `~/.ccswitch/ccs.json`
- Set up your Claude settings path

### List available profiles

```bash
ccswitch profiles
```

This displays all available profiles configured in your ccs.json file.

### Show current configuration

```bash
ccswitch show
```

This shows the currently active profile and its settings.

### Show a specific profile

```bash
ccswitch show <profile-name>
```

Displays the configuration for a specific profile without switching to it.

### Switch to a profile

```bash
ccswitch use <profile-name>
```

Switches to the specified profile by updating your Claude Code settings.

### Reset to default

```bash
ccswitch reset
```

Resets your Claude Code settings to the default profile.

### Update to latest version

```bash
ccswitch update
```

Updates ccswitch to the latest version from GitHub releases. This command will:

- Check for the latest version on GitHub
- Download and install the new version automatically
- Create a backup of the current version (removed on success)

**Update options:**

```bash
# Update to latest version
ccswitch update

# Update to a specific version
ccswitch update --version v1.0.0

# Force update even if already up-to-date
ccswitch update --force
```

## Configuration

The profiles are stored in `~/.ccswitch/ccs.json` (legacy) or `~/.config/ccswitch/config.json` (new). The configuration file has the following structure:

```json
{
    "settingsPath": "~/.claude/settings.json",
    "default": "default",
    "profiles": {
        "default": {
            "ANTHROPIC_BASE_URL": "https://api.anthropic.com",
            "ANTHROPIC_AUTH_TOKEN": "sk-your-token-here",
            "ANTHROPIC_MODEL": "claude-3-5-sonnet-20241022",
            "ANTHROPIC_SMALL_FAST_MODEL": "claude-3-haiku-20240307",
            "API_TIMEOUT_MS": "300000"
        },
        "custom-profile": {
            "ANTHROPIC_BASE_URL": "https://api.example.com",
            "ANTHROPIC_AUTH_TOKEN": "sk-your-custom-token",
            "ANTHROPIC_MODEL": "custom-model",
            "ANTHROPIC_SMALL_FAST_MODEL": "fast-model"
        }
    }
}
```

### Profile Settings

Each profile can contain the following settings:

- `ANTHROPIC_BASE_URL`: The API base URL
- `ANTHROPIC_AUTH_TOKEN`: Your authentication token
- `ANTHROPIC_MODEL`: The primary model to use
- `ANTHROPIC_SMALL_FAST_MODEL`: A smaller, faster model for quick tasks
- `ANTHROPIC_DEFAULT_SONNET_MODEL`: Default Sonnet model variant
- `ANTHROPIC_DEFAULT_OPUS_MODEL`: Default Opus model variant
- `ANTHROPIC_DEFAULT_HAIKU_MODEL`: Default Haiku model variant
- `API_TIMEOUT_MS`: API timeout in milliseconds

## Pre-configured Profiles

The tool comes with several pre-configured profiles for different Claude API providers:

- **default**: Official Anthropic API
- **anyrouter**: AnyRouter proxy service
- **glm**: Zhipu AI's GLM models
- **deepseek**: DeepSeek API
- **kimi-kfc**: Kimi Coding API
- **kimi-k2**: Kimi K2 API

## Security Considerations

- Your API tokens are stored in plain text in the configuration file
- Ensure your configuration file has appropriate permissions (readable only by you)
- Never commit your configuration file to version control
- Consider using environment variables for additional security
- The `ccswitch update` command verifies downloads using checksums when available

## Advanced Usage

### Custom Profile Creation

You can create custom profiles by editing the configuration file:

```bash
# Open the configuration file in your editor
nano ~/.ccswitch/ccs.json
```

Add a new profile to the `profiles` section:

```json
"my-custom-profile": {
    "ANTHROPIC_BASE_URL": "https://api.example.com",
    "ANTHROPIC_AUTH_TOKEN": "sk-your-token-here",
    "ANTHROPIC_MODEL": "your-model-name",
    "ANTHROPIC_SMALL_FAST_MODEL": "fast-model-name"
}
```

### Environment Variables

CCSwitch respects the following environment variables:

- `CCSWITCH_PROFILES_PATH`: Override the default profiles configuration file path
- `CCSWITCH_SETTINGS_PATH`: Override the default Claude settings file path

## Development

### Building

Build the binary:

```bash
make build
```

### Testing

Run tests:

```bash
make test
```

Run tests with coverage:

```bash
make test-coverage
```

This will generate a coverage report in `coverage.html`.

### Code Formatting

Format the code:

```bash
make fmt
```

Run static analysis:

```bash
make vet
```

### Project Structure

```
ccswitch/
├── cmd/                    # CLI commands
│   ├── init.go            # Initialize configuration command
│   ├── profiles.go        # List profiles command
│   ├── reset.go           # Reset to default command
│   ├── root.go            # Root command and setup
│   ├── show.go            # Show configuration command
│   ├── update.go          # Update tool command
│   └── use.go             # Use profile command
├── internal/              # Private application code
│   ├── claude/            # Claude API client and settings
│   ├── config/            # Configuration management
│   ├── jsonutil/          # JSON utilities
│   └── osutil/            # OS utilities
├── config/                # Default configurations
│   ├── ccs.json           # Basic profile configuration
│   └── ccs-full.json      # Complete profile configuration
├── install.sh             # Installation script
├── Makefile               # Build automation
├── main.go                # Application entry point
└── tests/                 # Test files (if any)
```

### Available Make Targets

Run `make help` to see all available targets:

```bash
make help
```

### Code Style

We follow the standard Go conventions:

- Use `gofmt` for code formatting
- Run `golangci-lint` for additional linting (optional)
- Write meaningful commit messages
- Add unit tests for new features

## Continuous Integration

The project uses GitHub Actions for CI/CD:

- **Test Workflow**: Runs on every push and pull request
  - Executes all unit tests with race detection
  - Generates coverage reports
  - Builds the binary and verifies it

- **Release Workflow**: Runs on version tags (e.g., `v1.0.0`)
  - Runs all tests
  - Builds binaries for multiple platforms (Linux, macOS, Windows)
  - Supports multiple architectures (amd64, arm64, arm)
  - Creates GitHub releases with binaries and checksums

## Requirements

- Go 1.21 or higher (for building from source)
- For binary installation: any modern operating system (Linux, macOS, Windows)

## Troubleshooting

### Common Issues

1. **"command not found: ccswitch"**
   - Ensure the installation directory is in your PATH
   - Try restarting your terminal or running `source ~/.bashrc` or `source ~/.zshrc`

2. **"permission denied"**
   - Make the binary executable: `chmod +x ~/.local/bin/ccswitch`
   - Check directory permissions

3. **"configuration file not found"**
   - Run `ccswitch init` to create the initial configuration
   - Check if `~/.ccswitch/ccs.json` exists

4. **"update failed"**
   - Check your internet connection
   - Try with `--force` flag to bypass version check
   - Manually download from releases page

### Getting Help

- Run `ccswitch --help` for command-line help
- Check the [GitHub Issues](https://github.com/huangdijia/ccswitch/issues) for known problems
- Create a new issue if you encounter bugs

## License

This project is licensed under the MIT License.

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests: `make test`
5. Format code: `make fmt`
6. Commit your changes: `git commit -m 'Add amazing feature'`
7. Push to the branch: `git push origin feature/amazing-feature`
8. Open a Pull Request

### Development Guidelines

- Follow the existing code style
- Add tests for new features
- Update documentation as needed
- Ensure all tests pass before submitting

## Changelog

### v0.3.0-beta.2
- Added online update functionality with `ccswitch update` command
- Enhanced security with path traversal protection
- Improved version comparison and update logic
- Added build information (version, commit, build date)

### v0.3.0-beta.1
- Complete rewrite from PHP to Go
- Added all original features and improvements
- Implemented comprehensive test suite
- Added CI/CD with GitHub Actions

### Previous Versions
- Originally implemented in PHP with Symfony Console
- Basic profile management functionality

## Support

If you encounter any issues or have questions, please file an issue on the [GitHub repository](https://github.com/huangdijia/ccswitch/issues).

## Acknowledgments

- Thanks to all [contributors](https://github.com/huangdijia/ccswitch/graphs/contributors) who have helped make this project better
- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Inspired by the need for seamless Claude Code API profile management
