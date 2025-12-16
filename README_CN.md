# CCSwitch

[![CI](https://github.com/huangdijia/ccswitch/workflows/Test/badge.svg)](https://github.com/huangdijia/ccswitch/actions)
[![Release](https://img.shields.io/github/release/huangdijia/ccswitch.svg)](https://github.com/huangdijia/ccswitch/releases)
[![Downloads](https://img.shields.io/github/downloads/huangdijia/ccswitch/total.svg)](https://github.com/huangdijia/ccswitch/releases)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[English](README.md) | 中文文档

一个强大的命令行工具，用于管理和切换不同 Claude Code API 配置文件和设置。

## 描述

CCSwitch 允许您轻松管理多个 Claude Code API 配置（配置文件）并在它们之间切换。当您需要为不同的项目或环境使用不同的 API 端点、模型或身份验证令牌时，这将非常有用。

### 主要特性

- **多配置文件管理**：存储并在多个 Claude API 配置之间切换
- **自动更新**：内置更新命令，保持工具为最新版本
- **预配置提供商**：开箱即用地支持各种 Claude API 提供商
- **跨平台**：支持 Linux、macOS 和 Windows
- **简单的 CLI**：直观的命令，便于配置文件管理
- **配置持久化**：您的设置被安全存储并自动应用

## 安装

### 快速安装（推荐）

安装 CCSwitch 最简单的方法是使用我们的安装脚本：

```bash
curl -sSL https://raw.githubusercontent.com/huangdijia/ccswitch/main/install.sh | bash
```

这将：

- 自动检测您的平台和架构
- 下载最新版本的二进制文件
- 安装到 `~/.local/bin` 目录
- 如果目录不在 PATH 中，会提供有用的说明

**安装选项：**

```bash
# 安装到自定义目录
curl -sSL https://raw.githubusercontent.com/huangdijia/ccswitch/main/install.sh | bash -s -- -d /usr/local/bin

# 安装特定版本
curl -sSL https://raw.githubusercontent.com/huangdijia/ccswitch/main/install.sh | bash -s -- -v v1.0.0
```

### 使用 Go Install

如果您已安装 Go 1.21 或更高版本，可以直接安装 CCSwitch：

```bash
go install github.com/huangdijia/ccswitch@latest
```

确保您的 Go bin 目录在 PATH 中：

```bash
export PATH="$HOME/go/bin:$PATH"
```


### 从源码编译

1. 克隆仓库：

```bash
git clone https://github.com/huangdijia/ccswitch.git
cd ccswitch
```

2. 编译并安装：

```bash
make install
```

或仅编译：

```bash
make build
./ccswitch init
```

### 二进制发布版本

从 [releases 页面](https://github.com/huangdijia/ccswitch/releases) 下载预编译的二进制文件。

## 快速开始

```bash
# 安装 ccswitch（如果尚未安装）
curl -sSL https://raw.githubusercontent.com/huangdijia/ccswitch/main/install.sh | bash

# 初始化您的配置
ccswitch init

# 列出可用的配置文件
ccswitch profiles

# 切换到某个配置文件
ccswitch use glm
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

### 更新到最新版本

```bash
ccswitch update
```

从 GitHub releases 更新 ccswitch 到最新版本。此命令将：

- 检查 GitHub 上的最新版本
- 自动下载并安装新版本
- 创建当前版本的备份（成功后删除）

**更新选项：**

```bash
# 更新到最新版本
ccswitch update

# 更新到特定版本
ccswitch update --version v1.0.0

# 即使已经是最新版本也强制更新
ccswitch update --force
```

## 配置

配置文件存储在 `~/.ccswitch/ccs.json`（旧版）或 `~/.config/ccswitch/config.json`（新版）中。配置文件具有以下结构：

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

## 安全注意事项

- 您的 API 令牌以明文形式存储在配置文件中
- 确保您的配置文件具有适当的权限（仅您可读）
- 切勿将您的配置文件提交到版本控制
- 考虑使用环境变量以增加安全性
- `ccswitch update` 命令在可用时会使用校验和验证下载

## 高级用法

### 创建自定义配置文件

您可以通过编辑配置文件来创建自定义配置文件：

```bash
# 在编辑器中打开配置文件
nano ~/.ccswitch/ccs.json
```

在 `profiles` 部分添加新的配置文件：

```json
"my-custom-profile": {
    "ANTHROPIC_BASE_URL": "https://api.example.com",
    "ANTHROPIC_AUTH_TOKEN": "sk-your-token-here",
    "ANTHROPIC_MODEL": "your-model-name",
    "ANTHROPIC_SMALL_FAST_MODEL": "fast-model-name"
}
```

### 环境变量

CCSwitch 尊重以下环境变量：

- `CCSWITCH_PROFILES_PATH`: 覆盖默认的配置文件路径
- `CCSWITCH_SETTINGS_PATH`: 覆盖默认的 Claude 设置文件路径

## 开发

### 编译

编译二进制文件：

```bash
make build
```

### 测试

运行测试：

```bash
make test
```

运行测试并生成覆盖率报告：

```bash
make test-coverage
```

这将在 `coverage.html` 中生成覆盖率报告。

### 代码格式化

格式化代码：

```bash
make fmt
```

运行静态分析：

```bash
make vet
```

### 项目结构

```
ccswitch/
├── cmd/                    # CLI 命令
│   ├── init.go            # 初始化配置命令
│   ├── profiles.go        # 列出配置文件命令
│   ├── reset.go           # 重置为默认命令
│   ├── root.go            # 根命令和设置
│   ├── show.go            # 显示配置命令
│   ├── update.go          # 更新工具命令
│   └── use.go             # 使用配置文件命令
├── internal/              # 私有应用程序代码
│   ├── claude/            # Claude API 客户端和设置
│   ├── config/            # 配置管理
│   ├── jsonutil/          # JSON 工具
│   └── osutil/            # OS 工具
├── config/                # 默认配置
│   ├── ccs.json           # 基本配置文件
│   └── ccs-full.json      # 完整配置文件
├── install.sh             # 安装脚本
├── Makefile               # 构建自动化
├── main.go                # 应用程序入口点
└── tests/                 # 测试文件（如果有）
```

### 可用的 Make 命令

运行 `make help` 查看所有可用命令：

```bash
make help
```

### 代码风格

我们遵循标准的 Go 约定：

- 使用 `gofmt` 进行代码格式化
- 运行 `golangci-lint` 进行额外的代码检查（可选）
- 编写有意义的提交信息
- 为新功能添加单元测试

## 持续集成

项目使用 GitHub Actions 进行 CI/CD：

- **测试工作流**：在每次推送和拉取请求时运行
  - 执行所有单元测试并进行竞态检测
  - 生成覆盖率报告
  - 构建二进制文件并验证

- **发布工作流**：在版本标签时运行（例如 `v1.0.0`）
  - 运行所有测试
  - 为多个平台构建二进制文件（Linux、macOS、Windows）
  - 支持多种架构（amd64、arm64、arm）
  - 创建 GitHub 发布版本，包含二进制文件和校验和

## 系统要求

- Go 1.21 或更高版本（用于从源码编译）
- 对于二进制安装：任何现代操作系统（Linux、macOS、Windows）

## 故障排除

### 常见问题

1. **"command not found: ccswitch"**
   - 确保安装目录在您的 PATH 中
   - 尝试重启终端或运行 `source ~/.bashrc` 或 `source ~/.zshrc`

2. **"permission denied"**
   - 使二进制文件可执行：`chmod +x ~/.local/bin/ccswitch`
   - 检查目录权限

3. **"configuration file not found"**
   - 运行 `ccswitch init` 创建初始配置
   - 检查 `~/.ccswitch/ccs.json` 是否存在

4. **"update failed"**
   - 检查您的网络连接
   - 尝试使用 `--force` 标志绕过版本检查
   - 手动从发布页面下载

### 获取帮助

- 运行 `ccswitch --help` 获取命令行帮助
- 查看 [GitHub Issues](https://github.com/huangdijia/ccswitch/issues) 了解已知问题
- 如果您遇到错误，请创建新的 issue

## 许可证

本项目采用 MIT 许可证。

## 贡献

我们欢迎贡献！请遵循以下步骤：

1. Fork 仓库
2. 创建功能分支：`git checkout -b feature/amazing-feature`
3. 进行您的更改
4. 运行测试：`make test`
5. 格式化代码：`make fmt`
6. 提交您的更改：`git commit -m 'Add amazing feature'`
7. 推送到分支：`git push origin feature/amazing-feature`
8. 打开 Pull Request

### 开发指南

- 遵循现有的代码风格
- 为新功能添加测试
- 根据需要更新文档
- 确保在提交前所有测试都通过

## 更新日志

### v0.3.0-beta.2
- 添加了在线更新功能和 `ccswitch update` 命令
- 增强了安全性，增加了路径遍历保护
- 改进了版本比较和更新逻辑
- 添加了构建信息（版本、提交、构建日期）

### v0.3.0-beta.1
- 从 PHP 完全重写为 Go
- 添加了所有原始功能和改进
- 实现了全面的测试套件
- 添加了 GitHub Actions CI/CD

### 以前的版本
- 最初使用 PHP 和 Symfony Console 实现
- 基本的配置文件管理功能

## 支持

如果您遇到任何问题或有疑问，请在 [GitHub 仓库](https://github.com/huangdijia/ccswitch/issues) 上提交 issue。

## 致谢

- 感谢所有[贡献者](https://github.com/huangdijia/ccswitch/graphs/contributors)帮助改进这个项目
- 使用 [Cobra](https://github.com/spf13/cobra) CLI 框架构建
- 灵感来源于对无缝 Claude Code API 配置文件管理的需求
