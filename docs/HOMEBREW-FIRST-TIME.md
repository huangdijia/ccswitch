# 首次提交 Homebrew Formula 指南

本文档详细说明如何首次将 ccswitch 提交到 homebrew-core。

## 前置条件

1. GitHub 账户
2. 已安装 Git
3. 对 Go 和 Homebrew 有基本了解

## 步骤 1: Fork homebrew-core

1. 访问 https://github.com/Homebrew/homebrew-core
2. 点击右上角的 "Fork" 按钮
3. 等待 fork 完成

## 步骤 2: 克隆你的 fork

```bash
# 克隆你 fork 的仓库
git clone https://github.com/YOUR_USERNAME/homebrew-core.git
cd homebrew-core

# 添加上游仓库
git remote add upstream https://github.com/Homebrew/homebrew-core.git

# 创建新分支
git checkout -b ccswitch
```

## 步骤 3: 创建 Formula 文件

在 `Formula/c` 目录下创建 `ccswitch.rb` 文件：

```bash
# 在 Formula/c 目录创建文件
vim Formula/c/ccswitch.rb
```

## 步骤 4: 获取正确的 SHA256

对于最新版本（例如 v0.4.0），需要计算 SHA256：

```bash
# 下载源码
curl -L https://github.com/huangdijia/ccswitch/archive/refs/tags/v0.4.0.tar.gz -o source.tar.gz

# 计算 SHA256
shasum -a 256 source.tar.gz
```

将计算出的 SHA256 替换 formula 中的 `PLACEHOLDER_SHA256`。

## 步骤 5: Formula 内容

使用以下内容（记得替换 SHA256）：

```ruby
class Ccswitch < Formula
  desc "Command-line tool for managing Claude Code API profiles"
  homepage "https://github.com/huangdijia/ccswitch"
  url "https://github.com/huangdijia/ccswitch/archive/refs/tags/v0.4.0.tar.gz"
  sha256 "dc035ee1de978240cdad7948e029a5470e2ebe1b8afa5d5bbf40387b0e373406"  # 替换为实际 SHA256
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version}"), "./cmd/ccswitch"
  end

  test do
    version_output = shell_output("#{bin}/ccswitch --version 2>&1")
    assert_match "ccswitch version #{version}", version_output

    # Test help command
    help_output = shell_output("#{bin}/ccswitch --help 2>&1")
    assert_match "A command-line tool for managing", help_output

    # Test that config directory can be created
    (testpath/".config/ccswitch").mkdir
    assert (testpath/".config/ccswitch").directory?
  end
end
```

## 步骤 6: 本地测试

在提交前，确保 formula 能正常工作：

```bash
# 安装 formula
brew install ./Formula/c/ccswitch.rb

# 测试安装
ccswitch --version

# 测试功能
ccswitch --help

# 卸载测试版本
brew uninstall ccswitch
```

## 步骤 7: 提交更改

```bash
# 添加文件
git add Formula/c/ccswitch.rb

# 提交
git commit -m "ccswitch: add version v0.4.0

- CLI tool for managing Claude Code API profiles
- Supports multiple API providers (Anthropic, AnyRouter, GLM, etc.)
- Cross-platform support (Linux, macOS, Windows)
- MIT licensed"

# 推送到你的 fork
git push origin ccswitch
```

## 步骤 8: 创建 Pull Request

1. 访问你的 fork 页面：https://github.com/YOUR_USERNAME/homebrew-core
2. 点击 "New pull request" 按钮
3. 确保分支是从 `ccswitch` 到 `homebrew-core:master`
4. 填写 PR 描述：

**PR 标题**: `ccswitch v0.4.0 (new formula)`

**PR 描述**:

```
Created with `brew create https://github.com/huangdijia/ccswitch/archive/refs/tags/v0.4.0.tar.gz`.

---

### `ccswitch`

Command-line tool for managing and switching between different Claude Code API profiles and configurations.

**Features**:
- Multi-profile management for different Claude API configurations
- Support for various API providers (Anthropic, AnyRouter, GLM, DeepSeek, Kimi)
- Cross-platform compatibility (Linux, macOS, Windows)
- Automatic update functionality
- Simple CLI interface

**Homepage**: https://github.com/huangdijia/ccswitch

**License**: MIT

- [x] Have you followed the guidelines in our [Contributing document](https://github.com/Homebrew/brew/blob/HEAD/docs/How-To-Open-a-Homebrew-Pull-Request.md)?
- [x] Have you checked that there aren't other open [pull requests](https://github.com/Homebrew/homebrew-core/pulls) for the same formula update/change?
- [x] Have you built your formula locally with `brew install --build-from-source ./Formula/ccswitch.rb`?
- [x] Does your formula build with `brew test ./Formula/ccswitch.rb`?
```

## 步骤 9: 等待审核

提交后，Homebrew 维护者会审核你的 PR：

1. **初始审核**: 通常 1-3 天
2. **可能的反馈**: 可能需要修改
3. **测试**: Homebrew CI 会自动测试
4. **合并**: 通过审核后会被合并

## 常见问题和解决方案

### 1. SHA256 不匹配

确保使用正确的版本标签下载文件计算 SHA256。

### 2. 测试失败

检查测试用例是否匹配实际输出。

### 3. 命名冲突

确保 `ccswitch` 名称未被占用。

### 4. 依赖问题

确保声明的依赖是正确的。

## 审核标准

Homebrew 会检查：

- ✅ 软件稳定且不是早期开发版本
- ✅ 有有效的测试
- ✅ 遵循命名约定
- ✅ 兼容的开源许可证
- ✅ 不与现有 formula 重复
- ✅ 符合 Homebrew 的标准

## 后续更新

一旦 formula 被合并，未来的更新将通过自动化的 GitHub Actions 完成，无需手动干预。

## 有用链接

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [贡献指南](https://github.com/Homebrew/brew/blob/HEAD/docs/How-To-Open-a-Homebrew-Pull-Request.md)
- [审核标准](https://docs.brew.sh/Acceptable-Formulae)
- [Homebrew Discourse](https://discourse.brew.sh/)