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
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/anyproto/anytype-cli/core"
	"github.com/spf13/cobra"
)

const (
	githubOwner = "anyproto"
	githubRepo  = "anytype-cli"
)

func NewUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update anytype CLI to the latest version",
		RunE:  runUpdate,
	}
}

func runUpdate(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking for updates...")

	latest, err := getLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to check latest version: %w", err)
	}

	current := core.GetVersion()

	currentBase := current
	if idx := strings.Index(current, "-"); idx != -1 {
		currentBase = current[:idx]
	}

	if currentBase >= latest {
		fmt.Printf("Already up to date (%s)\n", current)
		return nil
	}

	fmt.Printf("Updating from %s to %s...\n", current, latest)

	if err := downloadAndInstall(latest); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Printf("Successfully updated to %s\n", latest)
	fmt.Println("Restart your terminal or run 'anytype' to use the new version")
	return nil
}

func getLatestVersion() (string, error) {
	resp, err := githubAPI("GET", "/releases/latest", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", handleAPIError(resp)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse release: %w", err)
	}

	return release.TagName, nil
}

func downloadAndInstall(version string) error {
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

func getArchiveName(version string) string {
	base := fmt.Sprintf("anytype-cli-%s-%s-%s", version, runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		return base + ".zip"
	}
	return base + ".tar.gz"
}

func downloadRelease(version, destination string) error {
	archiveName := filepath.Base(destination)
	fmt.Printf("Downloading %s...\n", archiveName)

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return downloadViaAPI(version, archiveName, destination)
	}

	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
		githubOwner, githubRepo, version, archiveName)

	return downloadFile(url, destination, "")
}

func downloadViaAPI(version, filename, destination string) error {
	resp, err := githubAPI("GET", fmt.Sprintf("/releases/tags/%s", version), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleAPIError(resp)
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

	resp, err := http.DefaultClient.Do(req)
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
		if runtime.GOOS != "windows" {
			cmd := exec.Command("sudo", "mv", newBinary, currentBinary)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to replace binary: %w", err)
			}
		} else {
			return fmt.Errorf("failed to replace binary: %w", err)
		}
	}

	return nil
}

func githubAPI(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s%s", githubOwner, githubRepo, endpoint)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	return http.DefaultClient.Do(req)
}

func handleAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
}
