# Claude Code Switch

[![CI](https://github.com/huangdijia/ccswitch/workflows/Test/badge.svg)](https://github.com/huangdijia/ccswitch/actions)
[![Release](https://img.shields.io/github/release/huangdijia/ccswitch.svg)](https://github.com/huangdijia/ccswitch/releases)
[![Downloads](https://img.shields.io/github/downloads/huangdijia/ccswitch/total.svg)](https://github.com/huangdijia/ccswitch/releases)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

English | [中文文档](README_CN.md)

A powerful command-line tool for managing and switching between different Claude Code API profiles and configurations.

## Description

CCSwitch allows you to easily manage multiple Claude Code API configurations (profiles) and switch between them. This is useful when you need to use different API endpoints, models, or authentication tokens for different projects or environments.

### Key Features

- **Multi-profile Management**: Store and switch between multiple Claude API configurations
- **Interactive Profile Switching**: Interactive selection when no profile is specified
- **Auto-update**: Built-in update command to keep the tool current with GitHub releases
- **Pre-configured Providers**: Support for various Claude API providers out of the box
- **Custom Profile Creation**: Interactive profile creation with guided prompts
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

# Quick way: Install a pre-configured profile interactively
ccswitch install

# List available profiles
ccswitch list
# or
ccswitch ls
# or
ccswitch profiles

# Switch to a profile (interactive selection)
ccswitch use
# or specify profile name
ccswitch use glm

# Add a new custom profile
ccswitch add myprofile --api-key "sk-..." --model "opus"

# Show current settings
ccswitch show

# Show specific profile details
ccswitch show glm
```

## Usage

### Initialize configuration

```bash
ccswitch init
```

**Options:**

- `--full`: Use full configuration with all available providers
- `--force, -f`: Force overwrite existing configuration

This command initializes the CCSwitch configuration. It will:

- Create the config directory if it doesn't exist (`~/.ccswitch/`)
- Download/copy the default profile configuration
- Set up your Claude settings path

### Add a new profile

```bash
ccswitch add <profile-name> [flags]
```

**Flags:**

- `--api-key, -k`: Anthropic API key
- `--base-url, -u`: Anthropic base URL (default: <https://api.anthropic.com>)
- `--model, -m`: Anthropic model (default: opus)
- `--description, -d`: Profile description
- `--force, -f`: Force overwrite existing profile

This command allows interactive or non-interactive profile creation. Without flags, it will prompt for input interactively.

### Install a profile from preset configuration

```bash
ccswitch install [flags]
```

**Flags:**

- `--force, -f`: Force overwrite existing profile

This command provides an interactive way to install pre-configured profiles from the preset configuration (preset.json). It will:

1. Download the preset configuration from GitHub to a temporary directory
2. Present an interactive selection menu with all available profiles
3. Prompt for your authentication token (with masked display)
4. Save the selected profile to your local configuration
5. Automatically clean up temporary files

**Interactive workflow:**

```bash
ccswitch install
# Downloading preset configuration from GitHub...
# Select profile to install: (use ↑/↓, Enter to select, q to cancel)
# > default
#   anyrouter
#   glm
#   deepseek
#   kimi-kfc
#   ...
# Enter authentication token for profile 'glm': ****
# ✓ Profile 'glm' installed successfully!
```

This is the easiest way to quickly set up profiles from the extensive list of pre-configured API providers.

### List available profiles

```bash
ccswitch list
# Aliases: ls, profiles
```

This displays all available profiles in a nicely formatted table showing profile name, description, URL, model, and status (default).

### Show current configuration

```bash
ccswitch show
ccswitch show --current
```

This shows the currently active Claude settings.

### Show a specific profile

```bash
ccswitch show <profile-name>
```

Displays the configuration for a specific profile without switching to it. Shows the profile's description and all environment variables (with sensitive values masked).

### Switch to a profile

```bash
ccswitch use [profile-name]
```

**Interactive mode**: When no profile name is specified, opens a keyboard-driven selector (use ↑/↓ and Enter).

**Direct mode**: When a profile name is provided, directly switches to that profile.

Switches to the specified profile by updating your Claude settings file with the profile's environment variables.

### Reset to default

```bash
ccswitch reset
```

Resets your Claude settings to an empty state (removes all profile-specific settings).

### Update to latest version

```bash
ccswitch update
# Alias: up
```

Updates ccswitch to the latest version from GitHub releases. This command will:

- Check for the latest version on GitHub
- Download and install the new version automatically
- Create a backup of the current version (removed on success)

The update command supports cross-platform updates (Linux, macOS, Windows) with automatic architecture detection (x86_64, arm64, armv7).

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

The profiles are stored in `~/.ccswitch/ccs.json`. The configuration file has the following structure:

```json
{
    "settingsPath": "~/.claude/settings.json",
    "default": "default",
    "profiles": {
        "default": {
            "ANTHROPIC_API_KEY": "sk-XXXXXXXXXXXXXXXXXXXXXX",
            "ANTHROPIC_BASE_URL": "https://api.anthropic.com",
            "ANTHROPIC_MODEL": "opus",
            "ANTHROPIC_DEFAULT_HAIKU_MODEL": "haiku",
            "ANTHROPIC_DEFAULT_OPUS_MODEL": "opus",
            "ANTHROPIC_DEFAULT_SONNET_MODEL": "sonnet",
            "ANTHROPIC_SMALL_FAST_MODEL": "haiku"
        }
    },
    "descriptions": {
        "default": "Use default profile"
    }
}
```

### Profile Settings

Each profile can contain the following settings:

- `ANTHROPIC_API_KEY` or `ANTHROPIC_AUTH_TOKEN`: Your authentication token
- `ANTHROPIC_BASE_URL`: The API base URL
- `ANTHROPIC_MODEL`: The primary model to use
- `ANTHROPIC_SMALL_FAST_MODEL`: A smaller, faster model for quick tasks
- `ANTHROPIC_DEFAULT_SONNET_MODEL`: Default Sonnet model variant
- `ANTHROPIC_DEFAULT_OPUS_MODEL`: Default Opus model variant
- `ANTHROPIC_DEFAULT_HAIKU_MODEL`: Default Haiku model variant
- `API_TIMEOUT_MS`: API timeout in milliseconds

### Descriptions

The configuration also supports a `descriptions` field to store human-readable descriptions for each profile, which are displayed in the list command.

## Pre-configured Profiles

The tool comes with several pre-configured profiles for different Claude API providers:

- **default**: Official Anthropic API
- **anyrouter**: AnyRouter proxy service
- **glm**: Zhipu AI's GLM models
- **deepseek**: DeepSeek API
- **kimi-kfc**: Kimi Coding API
- **kimi-k2**: Kimi K2 API
- **modelscope**: ModelScope's API
- **minimaxi-m2**: MiniMax's Anthropic API
- **xiaomi-mimo**: Xiaomi Mimo's Anthropic API

## Security Considerations

- Your API tokens are stored in plain text in the configuration file
- Ensure your configuration file has appropriate permissions (readable only by you)
- Never commit your configuration file to version control
- Consider using environment variables for additional security
- The `ccswitch update` command verifies downloads using checksums when available

## Advanced Usage

### Custom Profile Creation

You can create custom profiles using the `add` command or by editing the configuration file directly:

**Using the add command:**

```bash
# Interactive mode
ccswitch add my-profile

# Non-interactive mode
ccswitch add my-profile \
    --api-key "sk-your-token" \
    --base-url "https://api.example.com" \
    --model "custom-model" \
    --description "My custom profile"

# Force overwrite existing
ccswitch add my-profile --force
```

**Manual configuration editing:**

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
},
"descriptions": {
    "my-custom-profile": "My custom profile description"
}
```

### Environment Variables

CCSwitch respects the following environment variables:

- `CCSWITCH_PROFILES_PATH`: Override the default profiles configuration file path
- `CCSWITCH_SETTINGS_PATH`: Override the default Claude settings file path

### Configuration File Locations

- **Profiles config (legacy)**: `~/.ccswitch/ccs.json`
- **Claude settings**: `~/.claude/settings.json` (default)

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
│   ├── list.go            # List profiles command
│   ├── add.go             # Add new profile command
│   ├── show.go            # Show configuration command
│   ├── use.go             # Use profile command
│   ├── reset.go           # Reset to default command
│   ├── update.go          # Update tool command
│   └── root.go            # Root command and setup
├── internal/              # Private application code
│   ├── cmdutil/           # Command utility functions
│   ├── output/            # Output formatting utilities
│   ├── pathutil/          # Path utilities
│   ├── profiles/          # Profile management
│   ├── settings/          # Claude settings management
│   └── httputil/          # HTTP utilities
├── config/                # Default configurations
│   ├── ccs.json           # Basic profile configuration
│   ├── ccs-full.json      # Complete profile configuration
│   └── preset.json        # Preset profile configuration for install command
├── install.sh             # Installation script
├── Makefile               # Build automation
├── main.go                # Application entry point
└── cmd/*_test.go          # Unit tests for each command
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
   - Check with `echo $PATH` that `~/.local/bin` or your install directory is included

2. **"permission denied"**
   - Make the binary executable: `chmod +x ~/.local/bin/ccswitch`
   - Check directory permissions: `ls -la ~/.local/bin/ccswitch`

3. **"configuration file not found"**
   - Run `ccswitch init` to create the initial configuration
   - Check both locations: `~/.ccswitch/ccs.json`

4. **"update failed"**
   - Check your internet connection
   - Try with `--force` flag to bypass version check
   - Manually download from [GitHub releases page](https://github.com/huangdijia/ccswitch/releases)

5. **"no profiles available"**
   - Ensure you initialized with `ccswitch init --full` for pre-configured profiles
   - Check your config file has valid JSON structure
   - Verify profiles exist with `ccswitch list`

6. **"profile not found"** when switching
   - List available profiles with `ccswitch list`
   - Check for typos in the profile name
   - Use tab completion if available

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

### Current Features

- **Interactive profile switching**: Use `ccswitch use` without arguments for interactive selection
- **Profile creation**: Add profiles interactively with `ccswitch add`
- **Smart defaults**: Missing model fields automatically filled from ANTHROPIC_MODEL
- **Description support**: Store and display profile descriptions
- **Masked sensitive values**: API keys are masked in output
- **Enhanced list view**: Comprehensive table view with all profile details
- **Cross-platform auto-update**: Single command updates from GitHub releases
- **Multiple provider profiles**: Pre-configured profiles for 9 different providers

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
