package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	// Default addresses
	LocalhostIP = "127.0.0.1"

	// Port configuration
	GRPCPort    = "31007"
	GRPCWebPort = "31008"
	APIPort     = "31009"

	// Full addresses
	DefaultGRPCAddress    = LocalhostIP + ":" + GRPCPort
	DefaultGRPCWebAddress = LocalhostIP + ":" + GRPCWebPort
	DefaultAPIAddress     = LocalhostIP + ":" + APIPort

	// URLs
	GRPCDNSAddress = "dns:///" + DefaultGRPCAddress

	// GitHub repository
	GitHubOwner = "anyproto"
	GitHubRepo  = "anytype-cli"

	// GitHub domains
	GitHubAPIBaseURL = "https://api.github.com"
	GitHubWebBaseURL = "https://github.com"

	// External URLs
	GitHubBaseURL    = GitHubWebBaseURL + "/" + GitHubOwner + "/" + GitHubRepo
	GitHubCommitURL  = GitHubBaseURL + "/commit/"
	GitHubReleaseURL = GitHubBaseURL + "/releases/tag/"

	// Anytype network address
	AnytypeNetworkAddress = "N83gJpVd9MuNRZAuJLZ7LiMntTThhPc6DtzWWVjb1M3PouVU"

	// Directory and file names
	AnytypeDirName = ".anytype"
	ConfigFileName = "config.json"
	DataDirName    = "data"
	LogsDirName    = "logs"
	AnytypeName    = "anytype"

	// Update related files
	UpdateCheckFile = ".update-check"
	UpdateLockFile  = ".update-lock"
	UpdateLogFile   = "autoupdate.log"

	// Update check interval
	UpdateCheckInterval = 24 * time.Hour
)

func GetWorkDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", AnytypeName)
	case "windows":
		return filepath.Join(homeDir, "AppData", "Roaming", AnytypeName)
	default:
		return filepath.Join(homeDir, ".config", AnytypeName)
	}
}

func GetConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, AnytypeDirName)
}

func GetConfigFilePath() string {
	return filepath.Join(GetConfigDir(), ConfigFileName)
}

func GetDataDir() string {
	if dataPath := os.Getenv("DATA_PATH"); dataPath != "" {
		return dataPath
	}
	return filepath.Join(GetWorkDir(), DataDirName)
}

func GetLogsDir() string {
	return filepath.Join(GetConfigDir(), LogsDirName)
}

func GetUpdateCheckFilePath() string {
	return filepath.Join(GetConfigDir(), UpdateCheckFile)
}

func GetUpdateLockFilePath() string {
	return filepath.Join(GetConfigDir(), UpdateLockFile)
}

func GetUpdateLogFilePath() string {
	return filepath.Join(GetLogsDir(), UpdateLogFile)
}
