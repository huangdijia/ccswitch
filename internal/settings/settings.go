package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// ClaudeSettings represents the Claude settings.json file
type ClaudeSettings struct {
	Path  string                 `json:"-"`
	Model string                 `json:"model,omitempty"`
	Env   map[string]interface{} `json:"env,omitempty"`
}

// New creates a new ClaudeSettings instance
func New(path string) (*ClaudeSettings, error) {
	// Expand home directory
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, path[1:])
	}

	settings := &ClaudeSettings{
		Path: path,
		Env:  make(map[string]interface{}),
	}

	// Create file if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
		if err := settings.Write(); err != nil {
			return nil, err
		}
	}

	// Read existing file
	if err := settings.Read(); err != nil {
		return nil, err
	}

	return settings, nil
}

// Read reads the settings from the file
func (s *ClaudeSettings) Read() error {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return err
	}

	// Parse JSON
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Extract model if present
	if model, ok := raw["model"].(string); ok {
		s.Model = model
	}

	// Extract env if present
	if env, ok := raw["env"].(map[string]interface{}); ok {
		s.Env = env
	} else {
		s.Env = make(map[string]interface{})
	}

	return nil
}

// Write writes the settings to the file
func (s *ClaudeSettings) Write() error {
	// Build output structure
	output := make(map[string]interface{})

	if s.Model != "" {
		output["model"] = s.Model
	}

	if len(s.Env) > 0 {
		output["env"] = s.Env
	}

	// Marshal to JSON with pretty print
	data, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(s.Path, data, 0644)
}
