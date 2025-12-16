package cmd

import (
	"runtime"
	"testing"
)

func TestDetectPlatform(t *testing.T) {
	platform := detectPlatform()

	// Should contain OS and architecture
	if platform == "" {
		t.Error("detectPlatform returned empty string")
	}

	// Platform should be in expected format
	expectedOS := map[string]string{
		"darwin":  "Darwin",
		"linux":   "Linux",
		"windows": "Windows",
	}

	expectedArch := map[string]string{
		"amd64": "x86_64",
		"arm64": "arm64",
		"arm":   "armv7",
	}

	osName := runtime.GOOS
	arch := runtime.GOARCH

	expectedOSName, ok := expectedOS[osName]
	if !ok {
		t.Logf("Warning: OS %s not in expected list", osName)
		return
	}

	expectedArchName, ok := expectedArch[arch]
	if !ok {
		t.Logf("Warning: Architecture %s not in expected list", arch)
		return
	}

	expectedPlatform := expectedOSName + "_" + expectedArchName
	if platform != expectedPlatform {
		t.Errorf("detectPlatform() = %s, want %s", platform, expectedPlatform)
	}
}

func TestIsVersionUpToDate(t *testing.T) {
	tests := []struct {
		name    string
		current string
		target  string
		want    bool
	}{
		{
			name:    "same version",
			current: "1.0.0",
			target:  "1.0.0",
			want:    true,
		},
		{
			name:    "same version with v prefix",
			current: "v1.0.0",
			target:  "v1.0.0",
			want:    true,
		},
		{
			name:    "different prefix same version",
			current: "1.0.0",
			target:  "v1.0.0",
			want:    true,
		},
		{
			name:    "different versions",
			current: "1.0.0",
			target:  "1.0.1",
			want:    false,
		},
		{
			name:    "different major versions",
			current: "1.0.0",
			target:  "2.0.0",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isVersionUpToDate(tt.current, tt.target); got != tt.want {
				t.Errorf("isVersionUpToDate(%q, %q) = %v, want %v", tt.current, tt.target, got, tt.want)
			}
		})
	}
}

func TestFindAsset(t *testing.T) {
	tests := []struct {
		name     string
		release  *GitHubRelease
		platform string
		want     string
	}{
		{
			name: "matching asset",
			release: &GitHubRelease{
				TagName: "v1.0.0",
				Assets: []struct {
					Name               string `json:"name"`
					BrowserDownloadURL string `json:"browser_download_url"`
				}{
					{
						Name:               "ccswitch_1.0.0_Linux_x86_64.tar.gz",
						BrowserDownloadURL: "https://example.com/ccswitch_1.0.0_Linux_x86_64.tar.gz",
					},
					{
						Name:               "ccswitch_1.0.0_Darwin_arm64.tar.gz",
						BrowserDownloadURL: "https://example.com/ccswitch_1.0.0_Darwin_arm64.tar.gz",
					},
				},
			},
			platform: "Linux_x86_64",
			want:     "https://example.com/ccswitch_1.0.0_Linux_x86_64.tar.gz",
		},
		{
			name: "no matching asset",
			release: &GitHubRelease{
				TagName: "v1.0.0",
				Assets: []struct {
					Name               string `json:"name"`
					BrowserDownloadURL string `json:"browser_download_url"`
				}{
					{
						Name:               "ccswitch_1.0.0_Linux_x86_64.tar.gz",
						BrowserDownloadURL: "https://example.com/ccswitch_1.0.0_Linux_x86_64.tar.gz",
					},
				},
			},
			platform: "Windows_x86_64",
			want:     "",
		},
		{
			name: "empty assets",
			release: &GitHubRelease{
				TagName: "v1.0.0",
				Assets: []struct {
					Name               string `json:"name"`
					BrowserDownloadURL string `json:"browser_download_url"`
				}{},
			},
			platform: "Linux_x86_64",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findAsset(tt.release, tt.platform); got != tt.want {
				t.Errorf("findAsset() = %v, want %v", got, tt.want)
			}
		})
	}
}
