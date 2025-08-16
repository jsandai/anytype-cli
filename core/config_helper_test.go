package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetStoredAccountId(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "anytype-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")
	testConfig := `{"accountId":"test-account-123"}`
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	accountId, err := GetStoredAccountId()
	if err == nil && accountId != "" {
		t.Logf("GetStoredAccountId() = %v", accountId)
	}
}

func TestGetStoredTechSpaceId(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "anytype-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, ".anytype", "config.json")
	os.MkdirAll(filepath.Dir(configPath), 0755)
	testConfig := `{"techSpaceId":"tech-space-789"}`
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	techSpaceId, err := GetStoredTechSpaceId()
	if err == nil && techSpaceId != "" {
		t.Logf("GetStoredTechSpaceId() = %v", techSpaceId)
	}
}

func TestLoadStoredConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "anytype-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, ".anytype", "config.json")
	os.MkdirAll(filepath.Dir(configPath), 0755)
	testConfig := `{
		"accountId":"test-account-123",
		"techSpaceId":"tech-space-789"
	}`
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	cfg, err := LoadStoredConfig()
	if err == nil && cfg != nil {
		if cfg.AccountId != "" || cfg.TechSpaceId != "" {
			t.Logf("LoadStoredConfig() loaded config with AccountId=%v, TechSpaceId=%v",
				cfg.AccountId, cfg.TechSpaceId)
		}
	}
}
