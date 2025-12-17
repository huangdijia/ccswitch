# Homebrew Formula for CCSwitch

This document describes how CCSwitch integrates with Homebrew.

## Automatic Updates

When a new release is published on GitHub, the following workflow automatically updates the Homebrew formula:

1. **Trigger**: The workflow is triggered when a new release is published (tags starting with `v`)
2. **Version Extraction**: The version is extracted from the release tag
3. **SHA256 Calculation**: The workflow downloads the source tarball and calculates its SHA256 checksum
4. **Formula Update**: Using `brew bump-formula-pr`, it creates a pull request in the homebrew-core repository with the updated version
5. **PR Review**: The pull request must be reviewed and merged by Homebrew maintainers (typically takes 2-7 days)
6. **Sync**: After merging, it may take additional time for the update to propagate to all users

⚠️ **Important**: The formula will not be immediately available after release. There is typically a delay of:

- **First-time submission**: 1-2 weeks (includes initial review)
- **Updates**: 2-7 days (PR review and merge)
- **Propagation**: Up to 24 hours after merge for full distribution

## Manual Formula Bumping

You can also manually trigger the formula update using the "Bump Homebrew Formula" workflow in the Actions tab:

1. Go to the [Actions page](https://github.com/huangdijia/ccswitch/actions)
2. Select "Bump Homebrew Formula" workflow
3. Click "Run workflow"
4. Enter the version number (e.g., `1.2.3`)
5. Click "Run workflow"

## Formula Structure

The Homebrew formula builds CCSwitch from source:

```ruby
class Ccswitch < Formula
  desc "CLI tool for managing and switching between different Claude Code AI profiles"
  homepage "https://github.com/huangdijia/ccswitch"
  url "https://github.com/huangdijia/ccswitch/archive/refs/tags/v#{version}.tar.gz"
  sha256 "..."
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version}"), "./cmd/ccswitch"
  end

  test do
    version_output = shell_output("#{bin}/ccswitch --version 2>&1")
    assert_match "ccswitch version #{version}", version_output
  end
end
```

## Installation

Once the formula is merged into homebrew-core:

```bash
# Install CCSwitch
brew install ccswitch

# Upgrade to the latest version
brew upgrade ccswitch

# Uninstall
brew uninstall ccswitch
```

## Testing the Formula Locally

To test the formula locally before submitting to homebrew-core:

```bash
# Create a local tap
brew tap-new username/tap

# Copy the formula to the tap
cp ccswitch.rb $(brew --repository)/Library/Taps/username/homebrew-tap/Formula/ccswitch.rb

# Install from the local tap
brew install username/tap/ccswitch
```

## First-Time Formula Submission

If CCSwitch is not yet in homebrew-core, you need to:

1. Fork the [homebrew-core](https://github.com/Homebrew/homebrew-core) repository
2. Create a new file in the `Formula` directory named `ccswitch.rb`
3. Copy the formula template from the workflow output
4. Submit a pull request with the following requirements:
   - Follow [Homebrew's formula style guide](https://docs.brew.sh/Formula-Cookbook)
   - Ensure all tests pass
   - Include a brief description of the software
   - Verify the license is compatible

## Requirements for Homebrew

- Software must be stable and not in early development
- Must have a working test in the formula
- Must follow Homebrew's naming conventions
- Must have a compatible open-source license
- Should not duplicate existing formulas

## Troubleshooting

### Formula not found

If the formula doesn't exist in homebrew-core yet:

- The workflow will generate a template for you
- You'll need to manually submit it as a new formula

### Build failures

If the build fails:

- Check if the build flags are compatible with the Go version
- Verify the source path is correct (`./cmd/ccswitch`)
- Ensure all dependencies are declared

### SHA256 mismatch

If the SHA256 doesn't match:

- Verify the correct URL is being used
- Check if the tarball has changed (e.g., different compression)
- Recalculate the SHA256 manually

## Related Links

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Contribution Guidelines](https://github.com/Homebrew/brew/blob/HEAD/docs/How-To-Open-a-Homebrew-Pull-Request.md)
- [brew bump-formula-pr documentation](https://man.he.net/man1/brew-bump-formula-pr)
