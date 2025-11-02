package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/service"

	serviceCmd "github.com/anyproto/anytype-cli/cmd/service"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
)

type updateCheck struct {
	LastCheck   time.Time `json:"lastCheck"`
	LastVersion string    `json:"lastVersion"`
}

func CheckAndUpdate() {
	go func() {
		if err := performUpdateCheck(); err != nil {
			logUpdateError(err)
		}
	}()
}

func logUpdateError(err error) {
	logPath := config.GetUpdateLogFilePath()
	logDir := filepath.Dir(logPath)

	if mkdirErr := os.MkdirAll(logDir, 0755); mkdirErr != nil {
		return
	}

	f, openErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		return
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf("Autoupdate error: %v\n", err)
}

func performUpdateCheck() error {
	configDir := config.GetConfigDir()
	if configDir == "" {
		return fmt.Errorf("could not determine config directory")
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	lockPath := config.GetUpdateLockFilePath()
	checkPath := config.GetUpdateCheckFilePath()

	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}
	defer func() {
		lockFile.Close()
		os.Remove(lockPath)
	}()

	if !shouldCheckForUpdate(checkPath) {
		return nil
	}

	latest, err := GetLatestVersion()
	if err != nil {
		return err
	}

	current := GetCurrentVersion()

	check := updateCheck{
		LastCheck:   time.Now(),
		LastVersion: latest,
	}
	saveCheckInfo(checkPath, check)

	if !NeedsUpdate(current, latest) {
		return nil
	}

	if !CanUpdateBinary() {
		return nil
	}

	if err := DownloadAndInstall(latest); err != nil {
		return err
	}

	output.Success("Anytype CLI has been automatically updated from %s to %s", current, latest)

	// Check if service is running and restart it automatically
	if err := restartServiceIfRunning(); err != nil {
		output.Warning("Failed to restart service: %v", err)
		output.Info("Restart manually with: anytype service restart")
	}

	return nil
}

// restartServiceIfRunning checks if the service is running and restarts it
func restartServiceIfRunning() error {
	s, err := serviceCmd.GetService()
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	// Check status
	status, err := s.Status()
	if err != nil {
		if errors.Is(err, service.ErrNotInstalled) {
			// Service not installed, nothing to restart
			return nil
		}
		return fmt.Errorf("failed to get service status: %w", err)
	}

	// Only restart if running
	if status == service.StatusRunning {
		output.Info("Restarting service with new binary...")
		if err := s.Restart(); err != nil {
			return fmt.Errorf("failed to restart service: %w", err)
		}
		output.Info("Service restarted successfully")
	}

	return nil
}

func shouldCheckForUpdate(checkPath string) bool {
	data, err := os.ReadFile(checkPath)
	if err != nil {
		return true
	}

	var check updateCheck
	if err := json.Unmarshal(data, &check); err != nil {
		return true
	}

	return time.Since(check.LastCheck) > config.UpdateCheckInterval
}

func saveCheckInfo(checkPath string, check updateCheck) error {
	data, err := json.Marshal(check)
	if err != nil {
		return fmt.Errorf("failed to marshal check info: %w", err)
	}

	if err := os.WriteFile(checkPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write check info: %w", err)
	}

	return nil
}
