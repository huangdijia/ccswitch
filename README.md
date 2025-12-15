# CCSwitch

English | [中文文档](README_CN.md)

A command-line tool for managing and switching between different Claude Code API profiles and configurations.

## Description

CCSwitch allows you to easily manage multiple Claude Code API configurations (profiles) and switch between them. This is useful when you need to use different API endpoints, models, or authentication tokens for different projects or environments.

## Installation

### Global Installation (Recommended)

Install CCSwitch globally using Composer:

```bash
composer global require huangdijia/ccswitch
```

After installation, make sure the global Composer bin directory is in your PATH. Add the following line to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
export PATH="$HOME/.composer/vendor/bin:$PATH"
```

Or for newer Composer versions:

```bash
export PATH="$HOME/.config/composer/vendor/bin:$PATH"
```

Now you can use the `ccswitch` command from anywhere:

```bash
ccswitch init
```

### Local Installation

Alternatively, you can clone the repository and install locally:

1. Clone the repository:

```bash
git clone https://github.com/huangdijia/ccswitch.git
cd ccswitch
```

2. Install dependencies:

```bash
composer install
```

3. Run the command:

```bash
./bin/ccswitch init
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

### Code Analysis

Run PHPStan for static analysis:

```bash
composer analyse
```

### Code Style

Fix code style issues:

```bash
composer cs-fix
```

## Requirements

- PHP 8.1 or higher
- Composer

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or have questions, please file an issue on the [GitHub repository](https://github.com/huangdijia/ccswitch/issues).
