package core

import (
	"fmt"

	"github.com/anyproto/anytype-cli/core/config"
)

// GetStoredAccountID retrieves the stored account ID from config
func GetStoredAccountID() (string, error) {
	configMgr := config.GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	cfg := configMgr.Get()
	if cfg.AccountID == "" {
		return "", fmt.Errorf("no account ID found in config")
	}

	return cfg.AccountID, nil
}

// GetStoredTechSpaceID retrieves the stored tech space ID from config
func GetStoredTechSpaceID() (string, error) {
	configMgr := config.GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	cfg := configMgr.Get()
	if cfg.TechSpaceID == "" {
		return "", fmt.Errorf("no tech space ID found in config")
	}

	return cfg.TechSpaceID, nil
}

// LoadStoredConfig loads and returns the entire config
func LoadStoredConfig() (*config.Config, error) {
	configMgr := config.GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return configMgr.Get(), nil
}
