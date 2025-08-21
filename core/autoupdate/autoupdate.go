package autoupdate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/anyproto/anytype-cli/core/update"
)

const (
	checkInterval   = 24 * time.Hour
	updateCheckFile = ".update-check"
	updateLockFile  = ".update-lock"
)

type updateCheck struct {
	LastCheck   time.Time `json:"lastCheck"`
	LastVersion string    `json:"lastVersion"`
}

func CheckAndUpdate() {
	go func() {
		if err := performUpdateCheck(); err != nil {
			return
		}
	}()
}

func performUpdateCheck() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	anytypeDir := filepath.Join(home, ".anytype")
	if err := os.MkdirAll(anytypeDir, 0755); err != nil {
		return err
	}

	lockPath := filepath.Join(anytypeDir, updateLockFile)
	checkPath := filepath.Join(anytypeDir, updateCheckFile)

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

	latest, err := update.GetLatestVersion()
	if err != nil {
		return err
	}

	current := update.GetCurrentVersion()

	check := updateCheck{
		LastCheck:   time.Now(),
		LastVersion: latest,
	}
	saveCheckInfo(checkPath, check)

	if !update.NeedsUpdate(current, latest) {
		return nil
	}

	if !update.CanUpdateBinary() {
		return nil
	}

	tempDir, err := os.MkdirTemp("", "anytype-autoupdate-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if err := update.DownloadAndInstall(latest); err != nil {
		return err
	}

	fmt.Printf("\n✓ Anytype CLI has been automatically updated from %s to %s\n\n", current, latest)
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

	return time.Since(check.LastCheck) > checkInterval
}

func saveCheckInfo(checkPath string, check updateCheck) {
	data, _ := json.Marshal(check)
	_ = os.WriteFile(checkPath, data, 0644)
}
