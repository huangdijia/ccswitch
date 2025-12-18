# Claude Code Switch

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
- **交互式配置切换**：不指定配置文件时提供交互式选择
- **自动更新**：内置更新命令，保持工具为最新版本（支持 GitHub 发布）
- **预配置提供商**：开箱即用地支持各种 Claude API 提供商
- **自定义配置创建**：带有引导提示的交互式配置创建
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
ccswitch list
# 或
ccswitch ls
# 或
ccswitch profiles

# 切换到某个配置文件（交互式选择）
ccswitch use
# 或指定配置文件名
ccswitch use glm

# 添加新的自定义配置文件
ccswitch add myprofile --api-key "sk-..." --model "opus"

# 显示当前设置
ccswitch show

# 显示特定配置文件详情
ccswitch show glm
```

## 使用方法

### 初始化配置

```bash
ccswitch init
```

**选项：**

- `--full`: 使用包含所有可用提供商的完整配置
- `--force, -f`: 强制覆盖现有配置

此命令初始化 CCSwitch 配置。它将：

- 如果不存在则创建配置目录（`~/.ccswitch/`）
- 下载/复制默认配置文件
- 设置您的 Claude 设置路径

### 添加配置文件

```bash
ccswitch add <配置文件名> [flags]
```

**标志：**

- `--api-key, -k`: Anthropic API 密钥
- `--base-url, -u`: Anthropic 基础 URL (默认: <https://api.anthropic.com>)
- `--model, -m`: Anthropic 模型 (默认: opus)
- `--description, -d`: 配置文件描述
- `--force, -f`: 强制覆盖现有配置文件

此命令支持交互式或非交互式配置文件创建。没有标志时，会提示输入。

### 列出可用的配置文件

```bash
ccswitch list
# 别名: ls, profiles
```

这会以格式化的表格显示所有可用配置文件，包括配置文件名称、描述、URL、模型和状态（默认）。

### 显示当前配置

```bash
ccswitch show
ccswitch show --current
```

这将显示当前激活的 Claude 设置。

### 显示特定配置文件

```bash
ccswitch show <配置文件名>
```

显示特定配置文件的配置，而不切换到该配置文件。显示配置文件的描述和所有环境变量（敏感值会被掩码）。

### 切换到某个配置文件

```bash
ccswitch use [配置文件名]
```

**交互模式**：当未指定配置文件名时，打开键盘选择器（↑/↓ 选择，回车确认，q 取消）。

**直接模式**：当提供配置文件名时，直接切换到该配置文件。

通过使用配置文件的配置文件环境变量更新您的 Claude 设置来切换配置文件。

### 重置为默认配置

```bash
ccswitch reset
```

将您的 Claude 设置重置为空状态（移除所有特定于配置文件的设置）。

### 更新到最新版本

```bash
ccswitch update
# 别名: up
```

从 GitHub 发布更新 ccswitch 到最新版本。此命令将：

- 检查 GitHub 上的最新版本
- 自动下载并安装新版本
- 创建当前版本的备份（成功后删除）

更新命令支持跨平台更新（Linux、macOS、Windows），支持自动架构检测（x86_64、arm64、armv7）。

**更新选项：**

```bash
# 更新到最新版本
ccswitch update

# 更新到特定版本
ccswitch update --version v1.0.0

# 即使已经是最新的版本也强制更新
ccswitch update --force
```

## 配置

配置文件存储在 `~/.ccswitch/ccs.json` 中。配置文件具有以下结构：

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

### 配置文件设置

每个配置文件可以包含以下设置：

- `ANTHROPIC_API_KEY` 或 `ANTHROPIC_AUTH_TOKEN`: 您的身份验证令牌
- `ANTHROPIC_BASE_URL`: API 基础 URL
- `ANTHROPIC_MODEL`: 要使用的主要模型
- `ANTHROPIC_SMALL_FAST_MODEL`: 用于快速任务的较小、更快的模型
- `ANTHROPIC_DEFAULT_SONNET_MODEL`: 默认 Sonnet 模型版本
- `ANTHROPIC_DEFAULT_OPUS_MODEL`: 默认 Opus 模型版本
- `ANTHROPIC_DEFAULT_HAIKU_MODEL`: 默认 Haiku 模型版本
- `API_TIMEOUT_MS`: API 超时时间（毫秒）

### 描述

配置文件还支持 `descriptions` 字段来存储每个配置文件的人类可读描述，这些描述会显示在列表命令中。

## 预配置的配置文件

该工具预配置了几个针对不同 Claude API 提供商的配置文件：

- **default**: 官方 Anthropic API
- **anyrouter**: AnyRouter 代理服务
- **glm**: 智谱 AI 的 GLM 模型
- **deepseek**: DeepSeek API
- **kimi-kfc**: Kimi Coding API
- **kimi-k2**: Kimi K2 API
- **modelscope**: ModelScope 的 API
- **minimaxi-m2**: MiniMax 的 Anthropic API
- **xiaomi-mimo**: Xiaomi Mimo 的 Anthropic API

## 安全注意事项

- 您的 API 令牌以明文形式存储在配置文件中
- 确保您的配置文件具有适当的权限（仅您可读）
- 切勿将您的配置文件提交到版本控制
- 考虑使用环境变量以增加安全性
- `ccswitch update` 命令在可用时会使用校验和验证下载

## 高级用法

### 创建自定义配置文件

您可以使用 `add` 命令或直接编辑配置文件来创建自定义配置文件：

**使用 add 命令：**

```bash
# 交互式模式
ccswitch add my-profile

# 非交互式模式
ccswitch add my-profile \
    --api-key "sk-your-token" \
    --base-url "https://api.example.com" \
    --model "custom-model" \
    --description "My custom profile"

# 强制覆盖现有配置
ccswitch add my-profile --force
```

**手动配置编辑：**

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
},
"descriptions": {
    "my-custom-profile": "My custom profile description"
}
```

### 环境变量

CCSwitch 尊重以下环境变量：

- `CCSWITCH_PROFILES_PATH`: 覆盖默认的配置文件路径
- `CCSWITCH_SETTINGS_PATH`: 覆盖默认的 Claude 设置文件路径

### 配置文件位置

- **配置文件 (旧版本)**: `~/.ccswitch/ccs.json`
- **Claude 设置**: `~/.claude/settings.json` (默认)

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
│   ├── list.go            # 列出配置文件命令
│   ├── add.go             # 添加新配置文件命令
│   ├── show.go            # 显示配置命令
│   ├── use.go             # 使用配置文件命令
│   ├── reset.go           # 重置为默认命令
│   ├── update.go          # 更新工具命令
│   └── root.go            # 根命令和设置
├── internal/              # 私有应用程序代码
│   ├── cmdutil/           # 命令工具函数
│   ├── output/            # 输出格式化工具
│   ├── pathutil/          # 路径工具
│   ├── profiles/          # 配置文件管理
│   ├── settings/          # Claude 设置管理
│   └── httputil/          # HTTP 工具
├── config/                # 默认配置
│   ├── ccs.json           # 基本配置文件
│   └── ccs-full.json      # 完整配置文件
├── install.sh             # 安装脚本
├── Makefile               # 构建自动化
├── main.go                # 应用程序入口点
└── cmd/*_test.go          # 每个命令的单元测试
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
   - 使用 `echo $PATH` 检查 `~/.local/bin` 或安装目录是否包含在路径中

2. **"permission denied"**
   - 使二进制文件可执行：`chmod +x ~/.local/bin/ccswitch`
   - 检查目录权限：`ls -la ~/.local/bin/ccswitch`

3. **"configuration file not found"**
   - 运行 `ccswitch init` 创建初始配置
   - 检查两个位置：`~/.ccswitch/ccs.json`

4. **"update failed"**
   - 检查您的网络连接
   - 尝试使用 `--force` 标志绕过版本检查
   - 手动从 [GitHub 发布页面](https://github.com/huangdijia/ccswitch/releases) 下载

5. **"no profiles available"**
   - 确保您使用 `ccswitch init --full` 进行初始化以获取预配置配置文件
   - 检查您的配置文件是否具有有效的 JSON 结构
   - 使用 `ccswitch list` 验证配置文件是否存在

6. **"profile not found"** 在切换时
   - 使用 `ccswitch list` 列出可用配置文件
   - 检查配置文件名称是否有拼写错误
   - 如果可用，使用 Tab 补全

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

### 当前功能

- **交互式配置切换**：使用 `ccswitch use` 不带参数进行交互式选择
- **配置文件创建**：使用 `ccswitch add` 交互式添加配置文件
- **智能默认值**：缺少的模型字段自动从 ANTHROPIC_MODEL 填充
- **描述支持**：存储和显示配置文件描述
- **敏感值掩码**：API 密钥在输出中被掩码
- **增强的列表视图**：包含所有配置文件详细信息的全面表格视图
- **跨平台自动更新**：单个命令从 GitHub 发布进行更新
- **多提供商配置文件**：9 个不同提供商的预配置配置文件

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
