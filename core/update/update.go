package update

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/hashicorp/go-version"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Minute,
}

func GetLatestVersion() (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", config.GitHubAPIBaseURL, config.GitHubOwner, config.GitHubRepo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse release: %w", err)
	}

	return release.TagName, nil
}

func NeedsUpdate(current, latest string) bool {
	currentVer, err := version.NewVersion(current)
	if err != nil {
		return false
	}

	latestVer, err := version.NewVersion(latest)
	if err != nil {
		return false
	}

	return currentVer.LessThan(latestVer)
}

func GetCurrentVersion() string {
	return core.GetVersion()
}

func DownloadAndInstall(version string) error {
	tempDir, err := os.MkdirTemp("", "anytype-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	archivePath := filepath.Join(tempDir, getArchiveName(version))
	if err := downloadRelease(version, archivePath); err != nil {
		return err
	}

	if err := extractArchive(archivePath, tempDir); err != nil {
		return fmt.Errorf("failed to extract: %w", err)
	}

	binaryName := "anytype"
	if runtime.GOOS == "windows" {
		binaryName = "anytype.exe"
	}

	newBinary := filepath.Join(tempDir, binaryName)
	if _, err := os.Stat(newBinary); err != nil {
		return fmt.Errorf("binary not found in archive (expected %s)", binaryName)
	}

	if err := replaceBinary(newBinary); err != nil {
		return fmt.Errorf("failed to install: %w", err)
	}

	return nil
}

func CanUpdateBinary() bool {
	execPath, err := os.Executable()
	if err != nil {
		return false
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return false
	}

	file, err := os.OpenFile(execPath, os.O_RDWR, 0)
	if err != nil {
		return false
	}
	_ = file.Close()
	return true
}

func getArchiveName(version string) string {
	base := fmt.Sprintf("anytype-cli-%s-%s-%s", version, runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		return base + ".zip"
	}
	return base + ".tar.gz"
}

func downloadRelease(version, destination string) error {
	archiveName := filepath.Base(destination)

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return downloadViaAPI(version, archiveName, destination)
	}

	url := fmt.Sprintf("%s/%s/%s/releases/download/%s/%s",
		config.GitHubWebBaseURL, config.GitHubOwner, config.GitHubRepo, version, archiveName)

	return downloadFile(url, destination, "")
}

func downloadViaAPI(version, filename, destination string) error {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/tags/%s",
		config.GitHubAPIBaseURL, config.GitHubOwner, config.GitHubRepo, version)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var release struct {
		Assets []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to parse release: %w", err)
	}

	var assetURL string
	for _, asset := range release.Assets {
		if asset.Name == filename {
			assetURL = asset.URL
			break
		}
	}
	if assetURL == "" {
		return fmt.Errorf("release asset %s not found", filename)
	}

	return downloadFile(assetURL, destination, os.Getenv("GITHUB_TOKEN"))
}

func downloadFile(url, destination, token string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if token != "" {
		req.Header.Set("Authorization", "token "+token)
		req.Header.Set("Accept", "application/octet-stream")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractArchive(archivePath, destDir string) error {
	if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, destDir)
	}
	return extractTarGz(archivePath, destDir)
}

// isValidExtractionPath validates that the target path is within the destination directory to prevent path traversal attacks
func isValidExtractionPath(target, destDir string) bool {
	cleanTarget := filepath.Clean(target)
	cleanDest := filepath.Clean(destDir)

	// Check if target starts with destDir
	rel, err := filepath.Rel(cleanDest, cleanTarget)
	if err != nil {
		return false
	}

	// If the relative path starts with "..", it's trying to escape the destDir
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != ".."
}

func extractTarGz(archivePath, destDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, header.Name)

		if !isValidExtractionPath(target, destDir) {
			return fmt.Errorf("illegal file path in archive: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := writeFile(target, tr, header.FileInfo().Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func extractZip(archivePath, destDir string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		target := filepath.Join(destDir, f.Name)

		if !isValidExtractionPath(target, destDir) {
			return fmt.Errorf("illegal file path in archive: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(target, f.Mode())
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		if err := writeFile(target, rc, f.Mode()); err != nil {
			rc.Close()
			return err
		}
		rc.Close()
	}
	return nil
}

func writeFile(path string, r io.Reader, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, r); err != nil {
		return err
	}

	return os.Chmod(path, mode)
}

func replaceBinary(newBinary string) error {
	if err := os.Chmod(newBinary, 0755); err != nil {
		return err
	}

	currentBinary, err := os.Executable()
	if err != nil {
		return err
	}
	currentBinary, err = filepath.EvalSymlinks(currentBinary)
	if err != nil {
		return err
	}

	if err := os.Rename(newBinary, currentBinary); err != nil {
		return fmt.Errorf("failed to replace binary (insufficient permissions): %w", err)
	}

	return nil
}
