# CCSwitch

[English](README.md) | 中文

一个用于管理和切换不同 Claude Code API 配置文件和设置的命令行工具。

## 项目简介

CCSwitch 允许您轻松管理多个 Claude Code API 配置（配置文件）并在它们之间切换。当您需要为不同的项目或环境使用不同的 API 端点、模型或身份验证令牌时，这将非常有用。

## 安装

### 全局安装（推荐）

使用 Composer 全局安装 CCSwitch：

```bash
composer global require huangdijia/ccswitch
```

安装后，确保全局 Composer bin 目录在您的 PATH 中。将以下行添加到您的 shell 配置文件（`~/.bashrc`、`~/.zshrc` 等）：

```bash
export PATH="$HOME/.composer/vendor/bin:$PATH"
```

或者对于较新版本的 Composer：

```bash
export PATH="$HOME/.config/composer/vendor/bin:$PATH"
```

现在您可以在任何地方使用 `ccswitch` 命令：

```bash
ccswitch init
```

### 本地安装

或者，您可以克隆仓库并在本地安装：

1. 克隆仓库：

```bash
git clone https://github.com/huangdijia/ccswitch.git
cd ccswitch
```

2. 安装依赖：

```bash
composer install
```

3. 运行命令：

```bash
./bin/ccswitch init
```

## 使用方法

### 初始化配置

```bash
ccswitch init
```

此命令初始化 CCSwitch 配置。它将：

- 如果不存在则创建 `~/.ccswitch/` 目录
- 将默认配置文件复制到 `~/.ccswitch/ccs.json`
- 设置您的 Claude 设置路径

### 列出可用的配置文件

```bash
ccswitch profiles
```

这将显示您的 ccs.json 文件中配置的所有可用配置文件。

### 显示当前配置

```bash
ccswitch show
```

这将显示当前活动的配置文件及其设置。

### 显示特定配置文件

```bash
ccswitch show <配置文件名>
```

显示特定配置文件的配置，而不切换到该配置文件。

### 切换到某个配置文件

```bash
ccswitch use <配置文件名>
```

通过更新您的 Claude Code 设置切换到指定的配置文件。

### 重置为默认配置

```bash
ccswitch reset
```

将您的 Claude Code 设置重置为默认配置文件。

## 配置

配置文件存储在 `~/.ccswitch/ccs.json` 中。配置文件具有以下结构：

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

### 配置文件设置

每个配置文件可以包含以下设置：

- `ANTHROPIC_BASE_URL`: API 基础 URL
- `ANTHROPIC_AUTH_TOKEN`: 您的身份验证令牌
- `ANTHROPIC_MODEL`: 要使用的主要模型
- `ANTHROPIC_SMALL_FAST_MODEL`: 用于快速任务的较小、更快的模型
- `ANTHROPIC_DEFAULT_SONNET_MODEL`: 默认 Sonnet 模型版本
- `ANTHROPIC_DEFAULT_OPUS_MODEL`: 默认 Opus 模型版本
- `ANTHROPIC_DEFAULT_HAIKU_MODEL`: 默认 Haiku 模型版本
- `API_TIMEOUT_MS`: API 超时时间（毫秒）

## 预配置的配置文件

该工具预配置了几个针对不同 Claude API 提供商的配置文件：

- **default**: 官方 Anthropic API
- **anyrouter**: AnyRouter 代理服务
- **glm**: 智谱 AI 的 GLM 模型
- **deepseek**: DeepSeek API
- **kimi-kfc**: Kimi Coding API
- **kimi-k2**: Kimi K2 API

## 开发

### 代码分析

运行 PHPStan 进行静态分析：

```bash
composer analyse
```

### 代码风格

修复代码风格问题：

```bash
composer cs-fix
```

## 系统要求

- PHP 8.1 或更高版本
- Composer

## 许可证

本项目采用 MIT 许可证。

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 支持

如果您遇到任何问题或有疑问，请在 [GitHub 仓库](https://github.com/huangdijia/ccswitch/issues) 上提交 issue。
