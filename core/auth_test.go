package core

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetDefaultDataPath(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		wantPath string
	}{
		{
			name:     "with DATA_PATH env",
			envValue: "/custom/data/path",
			wantPath: "/custom/data/path",
		},
		{
			name:     "without DATA_PATH env",
			envValue: "",
			wantPath: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalEnv := os.Getenv("DATA_PATH")
			defer func() {
				os.Setenv("DATA_PATH", originalEnv)
			}()

			if tt.envValue != "" {
				os.Setenv("DATA_PATH", tt.envValue)
			} else {
				os.Unsetenv("DATA_PATH")
			}

			got := getDefaultDataPath()

			if tt.wantPath != "" {
				if got != tt.wantPath {
					t.Errorf("getDefaultDataPath() = %v, want %v", got, tt.wantPath)
				}
			} else {
				if got == "" {
					t.Error("getDefaultDataPath() returned empty path")
				}

				if !strings.HasSuffix(got, "data") {
					t.Errorf("getDefaultDataPath() = %v, expected to end with 'data'", got)
				}
			}
		})
	}
}

func TestGetDefaultWorkDir(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		goos     string
		expected string
	}{
		{
			name:     "macOS",
			goos:     "darwin",
			expected: filepath.Join(homeDir, "Library", "Application Support", "anytype"),
		},
		{
			name:     "Windows",
			goos:     "windows",
			expected: filepath.Join(homeDir, "AppData", "Roaming", "anytype"),
		},
		{
			name:     "Linux",
			goos:     "linux",
			expected: filepath.Join(homeDir, ".config", "anytype"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if runtime.GOOS != tt.goos {
				t.Skipf("Skipping test for %s on %s", tt.goos, runtime.GOOS)
			}

			got := getDefaultWorkDir()
			if got != tt.expected {
				t.Errorf("getDefaultWorkDir() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateAccountKey(t *testing.T) {
	tests := []struct {
		name        string
		accountKey  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid bot account key",
			accountKey: "somevalidbase64encodedkeythatislongenough",
			wantErr:    false,
		},
		{
			name:        "empty account key",
			accountKey:  "",
			wantErr:     true,
			errContains: "bot account key cannot be empty",
		},
		{
			name:        "too short account key",
			accountKey:  "shortkey",
			wantErr:     true,
			errContains: "invalid bot account key format",
		},
		{
			name:       "minimum valid length key",
			accountKey: "12345678901234567890", // exactly 20 chars
			wantErr:    false,
		},
		{
			name:        "just under minimum length",
			accountKey:  "1234567890123456789", // 19 chars
			wantErr:     true,
			errContains: "invalid bot account key format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAccountKey(tt.accountKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAccountKey() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateAccountKey() error = %q, want to contain %q", err.Error(), tt.errContains)
				}
			}
		})
	}
}
