package profiles

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the profiles configuration
type Config struct {
	SettingsPath string                       `json:"settingsPath"`
	Default      string                       `json:"default"`
	Profiles     map[string]map[string]string `json:"profiles"`
	Descriptions map[string]string            `json:"descriptions,omitempty"`
}

// Profiles manages profile configurations
type Profiles struct {
	Path string
	Data *Config
}

// New creates a new Profiles instance
func New(path string) (*Profiles, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("profiles file not found: %s", path)
	}

	p := &Profiles{
		Path: path,
	}

	if err := p.load(); err != nil {
		return nil, err
	}

	return p, nil
}

// load reads the profiles configuration from file
func (p *Profiles) load() error {
	data, err := os.ReadFile(p.Path)
	if err != nil {
		return err
	}

	p.Data = &Config{
		Profiles:     make(map[string]map[string]string),
		Descriptions: make(map[string]string),
	}

	if err := json.Unmarshal(data, p.Data); err != nil {
		return err
	}

	return nil
}

// GetSettingsPath returns the settings path from config
func (p *Profiles) GetSettingsPath() string {
	return p.Data.SettingsPath
}

// Has checks if a profile exists
func (p *Profiles) Has(name string) bool {
	_, exists := p.Data.Profiles[name]
	return exists
}

// Default returns the default profile name
func (p *Profiles) Default() string {
	if p.Data.Default != "" {
		return p.Data.Default
	}
	return "default"
}

// Get returns the profile configuration with missing fields filled
func (p *Profiles) Get(name string) map[string]string {
	env, exists := p.Data.Profiles[name]
	if !exists {
		return make(map[string]string)
	}

	// Make a copy to avoid modifying the original
	result := make(map[string]string)
	for k, v := range env {
		result[k] = v
	}

	// Fill in missing model fields with ANTHROPIC_MODEL if present
	if model, ok := result["ANTHROPIC_MODEL"]; ok {
		missing := []string{
			"ANTHROPIC_DEFAULT_HAIKU_MODEL",
			"ANTHROPIC_DEFAULT_OPUS_MODEL",
			"ANTHROPIC_DEFAULT_SONNET_MODEL",
			"ANTHROPIC_SMALL_FAST_MODEL",
		}

		for _, key := range missing {
			if _, exists := result[key]; !exists {
				result[key] = model
			}
		}
	}

	return result
}

// GetAll returns all profile names
func (p *Profiles) GetAll() []string {
	names := make([]string, 0, len(p.Data.Profiles))
	for name := range p.Data.Profiles {
		names = append(names, name)
	}
	return names
}
