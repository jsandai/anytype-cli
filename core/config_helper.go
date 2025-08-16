package core

import (
	"fmt"

	"github.com/anyproto/anytype-cli/core/config"
)

// GetStoredAccountId retrieves the stored account Id from config
func GetStoredAccountId() (string, error) {
	configMgr := config.GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	cfg := configMgr.Get()
	if cfg.AccountId == "" {
		return "", fmt.Errorf("no account Id found in config")
	}

	return cfg.AccountId, nil
}

// GetStoredTechSpaceId retrieves the stored tech space Id from config
func GetStoredTechSpaceId() (string, error) {
	configMgr := config.GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	cfg := configMgr.Get()
	if cfg.TechSpaceId == "" {
		return "", fmt.Errorf("no tech space Id found in config")
	}

	return cfg.TechSpaceId, nil
}

// LoadStoredConfig loads and returns the entire config
func LoadStoredConfig() (*config.Config, error) {
	configMgr := config.GetConfigManager()
	if err := configMgr.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return configMgr.Get(), nil
}
