# CCSwitch

English | [中文文档](README_CN.md)

A command-line tool for managing and switching between different Claude Code API profiles and configurations.

## Description

CCSwitch allows you to easily manage multiple Claude Code API configurations (profiles) and switch between them. This is useful when you need to use different API endpoints, models, or authentication tokens for different projects or environments.

## Installation

### Quick Install (Recommended)

The easiest way to install CCSwitch is with our installation script:

```bash
curl -sSL https://github.com/huangdijia/ccswitch/install.sh | bash
```

This will:

- Detect your platform and architecture
- Download the latest release binary
- Install it to `~/.local/bin`
- Add helpful instructions if the directory isn't in your PATH

**Installation options:**

```bash
# Install to a custom directory
curl -sSL https://github.com/huangdijia/ccswitch/install.sh | bash -s -- -d /usr/local/bin

# Install a specific version
curl -sSL https://github.com/huangdijia/ccswitch/install.sh | bash -s -- -v v1.0.0
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

### Using Homebrew (macOS/Linux)

Coming soon...

### Binary Releases

Download pre-built binaries from the [releases page](https://github.com/huangdijia/ccswitch/releases).

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

## Configuration

The profiles are stored in `~/.ccswitch/ccs.json`. The configuration file has the following structure:

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

### Available Make Targets

Run `make help` to see all available targets:

```bash
make help
```

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

- Go 1.21 or higher

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or have questions, please file an issue on the [GitHub repository](https://github.com/huangdijia/ccswitch/issues).
