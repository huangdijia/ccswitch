package pathutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandHome(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "tilde only",
			path:    "~",
			wantErr: false,
		},
		{
			name:    "tilde with path",
			path:    "~/.config/test",
			wantErr: false,
		},
		{
			name:    "absolute path",
			path:    "/absolute/path",
			wantErr: false,
		},
		{
			name:    "relative path",
			path:    "relative/path",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExpandHome(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandHome() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.path == "~" {
					home, _ := os.UserHomeDir()
					if result != home {
						t.Errorf("ExpandHome() = %v, want %v", result, home)
					}
				} else if filepath.IsAbs(tt.path) || !filepath.IsAbs(result) && tt.path[0] != '~' {
					if result != tt.path {
						t.Errorf("ExpandHome() = %v, want %v", result, tt.path)
					}
				}
			}
		})
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "create new directory",
			path:    filepath.Join(tmpDir, "newdir"),
			wantErr: false,
		},
		{
			name:    "create nested directory",
			path:    filepath.Join(tmpDir, "parent", "child"),
			wantErr: false,
		},
		{
			name:    "existing directory",
			path:    tmpDir,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnsureDir(tt.path, 0755)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !FileExists(tt.path) {
				t.Errorf("EnsureDir() directory was not created: %s", tt.path)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "existing.txt")
	os.WriteFile(existingFile, []byte("test"), 0644)

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing file",
			path: existingFile,
			want: true,
		},
		{
			name: "non-existing file",
			path: filepath.Join(tmpDir, "nonexistent.txt"),
			want: false,
		},
		{
			name: "existing directory",
			path: tmpDir,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.path); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultSettingsPath(t *testing.T) {
	result := DefaultSettingsPath()
	if result == "" {
		t.Error("DefaultSettingsPath() returned empty string")
	}
}

func TestDefaultProfilesPath(t *testing.T) {
	result := DefaultProfilesPath()
	if result == "" {
		t.Error("DefaultProfilesPath() returned empty string")
	}
}
