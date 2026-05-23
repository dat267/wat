package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	APIKeys  map[string]string `json:"apiKeys"`
	Settings map[string]string `json:"settings"`
}

type SchemaKey struct {
	Category    string
	IsMandatory bool
	IsSecret    bool
	Default     string
}
var ConfigSchema = map[string]SchemaKey{
	"openai":    {"apiKeys", false, true, ""},
	"anthropic": {"apiKeys", false, true, ""},
	"theme":     {"settings", true, false, "system"},
}
func IsKeySecret(category string, key string) bool {
	if s, found := ConfigSchema[key]; found {
		return s.IsSecret
	}
	return category == "apiKeys"
}
func IsKeyMandatory(key string) bool {
	if s, found := ConfigSchema[key]; found {
		return s.IsMandatory
	}
	return false
}
func GetDefaultValue(key string) string {
	if s, found := ConfigSchema[key]; found {
		return s.Default
	}
	return ""
}
func getConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "wat", "config.json"), nil
}
func loadConfig() (*AppConfig, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := &AppConfig{
				APIKeys:  make(map[string]string),
				Settings: make(map[string]string),
			}
			for k, s := range ConfigSchema {
				if s.IsMandatory && s.Default != "" {
					if s.Category == "apiKeys" {
						cfg.APIKeys[k] = s.Default
					} else {
						cfg.Settings[k] = s.Default
					}
				}
			}
			return cfg, nil
		}
		return nil, err
	}
	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.APIKeys == nil {
		cfg.APIKeys = make(map[string]string)
	}
	if cfg.Settings == nil {
		cfg.Settings = make(map[string]string)
	}
	for k, s := range ConfigSchema {
		if s.IsMandatory && s.Default != "" {
			if s.Category == "apiKeys" {
				if _, exists := cfg.APIKeys[k]; !exists {
					cfg.APIKeys[k] = s.Default
				}
			} else {
				if _, exists := cfg.Settings[k]; !exists {
					cfg.Settings[k] = s.Default
				}
			}
		}
	}
	return &cfg, nil
}
func saveConfig(cfg *AppConfig) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
func (a *App) GetConfig() (AppConfig, error) {
	cfg, err := loadConfig()
	if err != nil {
		return AppConfig{}, err
	}
	return *cfg, nil
}
func (a *App) SetConfigValue(category string, key string, val string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if category == "apiKeys" {
		cfg.APIKeys[key] = val
	} else {
		cfg.Settings[key] = val
	}
	return saveConfig(cfg)
}
func (a *App) DeleteConfigValue(category string, key string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if category == "apiKeys" {
		if IsKeyMandatory(key) {
			cfg.APIKeys[key] = GetDefaultValue(key)
		} else {
			delete(cfg.APIKeys, key)
		}
	} else {
		if IsKeyMandatory(key) {
			cfg.Settings[key] = GetDefaultValue(key)
		} else {
			delete(cfg.Settings, key)
		}
	}
	return saveConfig(cfg)
}
func (a *App) ResetConfig() error {
	cfg := &AppConfig{
		APIKeys:  make(map[string]string),
		Settings: make(map[string]string),
	}
	for k, s := range ConfigSchema {
		if s.IsMandatory && s.Default != "" {
			if s.Category == "apiKeys" {
				cfg.APIKeys[k] = s.Default
			} else {
				cfg.Settings[k] = s.Default
			}
		}
	}
	return saveConfig(cfg)
}
