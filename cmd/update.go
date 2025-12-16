package cmd

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	updateForce   bool
	updateVersion string
)

const (
	repo = "huangdijia/ccswitch"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update ccswitch to the latest version",
	Long: `Update ccswitch to the latest version from GitHub releases.
This command will download and install the latest version of ccswitch from GitHub.
If you want to install a specific version, use the --version flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get executable path
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}
		exePath, err = filepath.EvalSymlinks(exePath)
		if err != nil {
			return fmt.Errorf("failed to resolve symlink: %w", err)
		}

		currentVer := rootCmd.Version
		fmt.Printf("Current version: %s\n", currentVer)
		fmt.Printf("Executable path: %s\n", exePath)

		// Determine target version
		targetVersion := updateVersion
		var release *GitHubRelease

		if targetVersion == "" {
			// Get latest release
			fmt.Println("Fetching latest release information...")
			release, err = getLatestRelease()
			if err != nil {
				return fmt.Errorf("failed to get latest release: %w", err)
			}
			targetVersion = release.TagName
		} else {
			// Get specific release
			fmt.Printf("Fetching release %s information...\n", targetVersion)
			release, err = getRelease(targetVersion)
			if err != nil {
				return fmt.Errorf("failed to get release %s: %w", targetVersion, err)
			}
		}

		fmt.Printf("Target version: %s\n", targetVersion)

		// Check if already up-to-date
		if !updateForce && isVersionUpToDate(currentVer, targetVersion) {
			fmt.Println("✓ Already up-to-date!")
			return nil
		}

		// Detect platform
		platform := detectPlatform()
		fmt.Printf("Platform: %s\n", platform)

		// Find matching asset
		assetURL := findAsset(release, platform)
		if assetURL == "" {
			return fmt.Errorf("no binary found for platform: %s", platform)
		}

		// Download and install
		fmt.Printf("Downloading %s...\n", assetURL)
		if err := downloadAndInstall(assetURL, exePath); err != nil {
			return fmt.Errorf("failed to update: %w", err)
		}

		fmt.Printf("✓ Successfully updated to version %s!\n", targetVersion)
		fmt.Println("Please restart ccswitch to use the new version.")

		return nil
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&updateForce, "force", "f", false, "Force update even if already up-to-date")
	updateCmd.Flags().StringVarP(&updateVersion, "version", "v", "", "Update to a specific version (e.g., v1.0.0)")
}

// getLatestRelease fetches the latest release from GitHub API
func getLatestRelease() (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	return fetchRelease(url)
}

// getRelease fetches a specific release from GitHub API
func getRelease(version string) (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/%s", repo, version)
	return fetchRelease(url)
}

// fetchRelease fetches release information from GitHub API
func fetchRelease(url string) (*GitHubRelease, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// detectPlatform detects the current platform and architecture
func detectPlatform() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Convert Go OS names to our naming convention
	switch osName {
	case "darwin":
		osName = "Darwin"
	case "linux":
		osName = "Linux"
	case "windows":
		osName = "Windows"
	}

	// Convert Go arch names to our naming convention
	switch arch {
	case "amd64":
		arch = "x86_64"
	case "arm64":
		arch = "arm64"
	case "arm":
		arch = "armv7"
	}

	return fmt.Sprintf("%s_%s", osName, arch)
}

// findAsset finds the matching asset URL for the platform
func findAsset(release *GitHubRelease, platform string) string {
	version := strings.TrimPrefix(release.TagName, "v")
	expectedName := fmt.Sprintf("ccswitch_%s_%s.tar.gz", version, platform)

	for _, asset := range release.Assets {
		if asset.Name == expectedName {
			return asset.BrowserDownloadURL
		}
	}

	return ""
}

// isVersionUpToDate checks if the current version is up-to-date
func isVersionUpToDate(current, target string) bool {
	// Remove 'v' prefix if present
	current = strings.TrimPrefix(current, "v")
	target = strings.TrimPrefix(target, "v")

	// Compare semantic versions
	return compareVersions(current, target) >= 0
}

// compareVersions compares two semantic version strings
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	// Split versions into parts
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Pad the shorter version with zeros
	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for len(parts1) < maxLen {
		parts1 = append(parts1, "0")
	}
	for len(parts2) < maxLen {
		parts2 = append(parts2, "0")
	}

	// Compare each part
	for i := 0; i < maxLen; i++ {
		// Parse as integers, ignoring non-numeric suffixes
		var n1, n2 int
		fmt.Sscanf(parts1[i], "%d", &n1)
		fmt.Sscanf(parts2[i], "%d", &n2)

		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}

	return 0
}

// downloadAndInstall downloads the binary and installs it
func downloadAndInstall(url, exePath string) error {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "ccswitch-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download archive
	archivePath := filepath.Join(tempDir, "ccswitch.tar.gz")
	if err := downloadFile(url, archivePath); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	// Extract archive
	binaryPath := filepath.Join(tempDir, "ccswitch")
	if err := extractTarGz(archivePath, tempDir); err != nil {
		return fmt.Errorf("failed to extract: %w", err)
	}

	// Verify binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("binary not found in archive")
	}

	// Make binary executable
	if err := os.Chmod(binaryPath, 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	// Create backup of current binary
	backupPath := exePath + ".backup"
	if err := copyFile(exePath, backupPath); err != nil {
		fmt.Printf("Warning: failed to create backup: %v\n", err)
	} else {
		defer func() {
			// Remove backup on success
			os.Remove(backupPath)
		}()
	}

	// Replace current binary
	if err := os.Rename(binaryPath, exePath); err != nil {
		// On Windows or if rename fails, try copy and delete
		if err := copyFile(binaryPath, exePath); err != nil {
			return fmt.Errorf("failed to replace binary: %w", err)
		}
	}

	return nil
}

// downloadFile downloads a file from a URL
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// extractTarGz extracts a tar.gz archive
func extractTarGz(archivePath, destDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Prevent path traversal attacks (zip slip)
		// Use filepath.Rel to ensure the path is relative and within destDir
		cleanName := filepath.Clean(header.Name)
		if strings.HasPrefix(cleanName, "..") || filepath.IsAbs(cleanName) {
			return fmt.Errorf("invalid file path in archive: %s", header.Name)
		}

		target := filepath.Join(destDir, cleanName)

		// Double-check: Ensure the target is within destDir using filepath.Rel
		rel, err := filepath.Rel(destDir, target)
		if err != nil || strings.HasPrefix(rel, "..") {
			return fmt.Errorf("invalid file path in archive: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// Ensure parent directory exists
			parentDir := filepath.Dir(target)
			if parentDir != "." {
				if err := os.MkdirAll(parentDir, 0755); err != nil {
					return err
				}
			}
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		case tar.TypeSymlink:
			// Ensure parent directory exists for symlinks
			parentDir := filepath.Dir(target)
			if parentDir != "." {
				if err := os.MkdirAll(parentDir, 0755); err != nil {
					return err
				}
			}
			if err := os.Symlink(header.Linkname, target); err != nil {
				return fmt.Errorf("failed to create symlink: %w", err)
			}
		}
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Copy permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, sourceInfo.Mode())
}
