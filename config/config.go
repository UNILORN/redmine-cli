package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Profile struct {
	Name       string `yaml:"name"`
	RedmineURL string `yaml:"redmine_url"`
	APIKey     string `yaml:"api_key"`
}

type Config struct {
	DefaultProfile string             `yaml:"default_profile"`
	Profiles       map[string]Profile `yaml:"profiles"`
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".redminecli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(configDir, "config"), nil
}

func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return &Config{
			Profiles: make(map[string]Profile),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if config.Profiles == nil {
		config.Profiles = make(map[string]Profile)
	}

	return &config, nil
}

func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (c *Config) GetCurrentProfile() (*Profile, error) {
	if c.DefaultProfile == "" && len(c.Profiles) == 0 {
		return nil, fmt.Errorf("no profiles configured")
	}

	profileName := c.DefaultProfile
	if profileName == "" {
		// If no default profile, use the first available profile
		for name := range c.Profiles {
			profileName = name
			break
		}
	}

	profile, exists := c.Profiles[profileName]
	if !exists {
		return nil, fmt.Errorf("profile '%s' not found", profileName)
	}

	return &profile, nil
}

func (c *Config) AddProfile(name, url, apiKey string) error {
	if c.Profiles == nil {
		c.Profiles = make(map[string]Profile)
	}

	c.Profiles[name] = Profile{
		Name:       name,
		RedmineURL: url,
		APIKey:     apiKey,
	}

	// Set as default if it's the first profile
	if c.DefaultProfile == "" {
		c.DefaultProfile = name
	}

	return nil
}

func (c *Config) SetDefaultProfile(name string) error {
	if _, exists := c.Profiles[name]; !exists {
		return fmt.Errorf("profile '%s' does not exist", name)
	}
	c.DefaultProfile = name
	return nil
}

func (c *Config) RemoveProfile(name string) error {
	if _, exists := c.Profiles[name]; !exists {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	delete(c.Profiles, name)

	// If removing the default profile, set a new default
	if c.DefaultProfile == name {
		c.DefaultProfile = ""
		for profileName := range c.Profiles {
			c.DefaultProfile = profileName
			break
		}
	}

	return nil
}