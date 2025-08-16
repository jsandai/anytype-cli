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

func TestValidateMnemonic(t *testing.T) {
	tests := []struct {
		name        string
		mnemonic    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid 12 word mnemonic",
			mnemonic: "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11 word12",
			wantErr:  false,
		},
		{
			name:     "valid 12 words with extra spaces",
			mnemonic: "word1  word2   word3 word4 word5 word6 word7 word8 word9 word10 word11 word12",
			wantErr:  false,
		},
		{
			name:        "invalid 11 word mnemonic",
			mnemonic:    "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11",
			wantErr:     true,
			errContains: "must be exactly 12 words, got 11",
		},
		{
			name:        "invalid 13 word mnemonic",
			mnemonic:    "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11 word12 word13",
			wantErr:     true,
			errContains: "must be exactly 12 words, got 13",
		},
		{
			name:        "empty mnemonic",
			mnemonic:    "",
			wantErr:     true,
			errContains: "cannot be empty",
		},
		{
			name:        "whitespace only",
			mnemonic:    "   ",
			wantErr:     true,
			errContains: "must be exactly 12 words, got 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMnemonic(tt.mnemonic)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMnemonic() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateMnemonic() error = %q, want to contain %q", err.Error(), tt.errContains)
				}
			}
		})
	}
}
