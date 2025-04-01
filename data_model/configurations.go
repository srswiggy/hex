package data_model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Configurations struct {
	showMockServices    bool
	showOnScreenOptions bool
	selectedPods        []string
}

// GetConfigFilePath returns the path to the configuration file
func GetConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create directory for config file if it doesn't exist
	configDir := filepath.Join(homeDir, ".hex")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(configDir, "config.json"), nil
}

// LoadConfig reads the configuration from file
func LoadConfig() (*Configurations, error) {
	// Default configuration
	config := &Configurations{
		showMockServices:    true,
		showOnScreenOptions: true,
		selectedPods:        []string{},
	}

	configPath, err := GetConfigFilePath()
	if err != nil {
		return config, err
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config file
		if err := config.Save(); err != nil {
			return config, fmt.Errorf("failed to create default config: %w", err)
		}
		return config, nil
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal into temp struct to handle unexported fields
	var jsonData struct {
		ShowMockServices    bool     `json:"showMockServices"`
		ShowOnScreenOptions bool     `json:"showOnScreenOptions"`
		SelectedPods        []string `json:"selectedPods"`
	}

	if err := json.Unmarshal(data, &jsonData); err != nil {
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Map to our unexported fields
	config.showMockServices = jsonData.ShowMockServices
	config.showOnScreenOptions = jsonData.ShowOnScreenOptions
	config.selectedPods = jsonData.SelectedPods

	return config, nil
}

// Save writes the configuration to file
func (c *Configurations) Save() error {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	// Map unexported fields to exportable struct for JSON serialization
	jsonData := struct {
		ShowMockServices    bool     `json:"showMockServices"`
		ShowOnScreenOptions bool     `json:"showOnScreenOptions"`
		SelectedPods        []string `json:"selectedPods"`
	}{
		ShowMockServices:    c.showMockServices,
		ShowOnScreenOptions: c.showOnScreenOptions,
		SelectedPods:        c.selectedPods,
	}

	// Marshal with indentation for readability
	data, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Getter and setter methods for the unexported fields
func (c *Configurations) GetShowMockServices() bool {
	return c.showMockServices
}

func (c *Configurations) SetShowMockServices(value bool) {
	c.showMockServices = value
}

func (c *Configurations) GetShowOnScreenOptions() bool {
	return c.showOnScreenOptions
}

func (c *Configurations) SetShowOnScreenOptions(value bool) {
	c.showOnScreenOptions = value
}

func (c *Configurations) GetSelectedPods() []string {
	return c.selectedPods
}

func (c *Configurations) SetSelectedPods(pods []string) {
	c.selectedPods = pods
}
