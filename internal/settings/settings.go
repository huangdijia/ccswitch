package settings

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/huangdijia/ccswitch/internal/pathutil"
)

// ClaudeSettings represents the Claude settings.json file
type ClaudeSettings struct {
	Path  string                 `json:"-"`
	Model string                 `json:"model,omitempty"`
	Env   map[string]interface{} `json:"env,omitempty"`
	raw   map[string]interface{} `json:"-"` // Store all fields to preserve them
}

// New creates a new ClaudeSettings instance
func New(path string) (*ClaudeSettings, error) {
	// Expand home directory
	expandedPath, err := pathutil.ExpandHome(path)
	if err != nil {
		return nil, err
	}
	path = expandedPath

	settings := &ClaudeSettings{
		Path: path,
		Env:  make(map[string]interface{}),
		raw:  make(map[string]interface{}),
	}

	// Create file if it doesn't exist
	if !pathutil.FileExists(path) {
		dir := filepath.Dir(path)
		if err := pathutil.EnsureDir(dir, 0755); err != nil {
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

	// Store all fields to preserve them
	s.raw = raw

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
	// Start with existing fields
	var output map[string]interface{}
	if s.raw != nil {
		output = make(map[string]interface{})
		for k, v := range s.raw {
			output[k] = v
		}
	} else {
		output = make(map[string]interface{})
	}

	// Update model if set
	if s.Model != "" {
		output["model"] = s.Model
	} else if _, exists := output["model"]; exists {
		// Remove model if empty and it existed
		delete(output, "model")
	}

	// Update env if set
	if len(s.Env) > 0 {
		output["env"] = s.Env
	} else if _, exists := output["env"]; exists {
		// Remove env if empty and it existed
		delete(output, "env")
	}

	// Marshal to JSON with pretty print
	data, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(s.Path, data, 0644)
}
