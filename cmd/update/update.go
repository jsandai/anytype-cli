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

	"github.com/minio/selfupdate"
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update to the latest version",
		Long:  "Download and install the latest version of the Anytype CLI from GitHub releases.",
		RunE:  runUpdate,
	}
}

func runUpdate(cmd *cobra.Command, args []string) error {
	output.Info("Checking for updates...")

	latest, err := getLatestVersion()
	if err != nil {
		return output.Error("Failed to check latest version: %w", err)
	}

	current := core.GetVersion()

	currentBase := current
	if idx := strings.Index(current, "-"); idx != -1 {
		currentBase = current[:idx]
	}

	if currentBase >= latest {
		output.Info("Already up to date (%s)", current)
		return nil
	}

	output.Info("Updating from %s to %s...", current, latest)

	if err := downloadAndInstall(latest); err != nil {
		return output.Error("update failed: %w", err)
	}

	output.Success("Successfully updated to %s", latest)
	output.Info("If the service is installed, restart it with: anytype service restart")
	output.Info("Otherwise, restart your terminal or run 'anytype' to use the new version")
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
		return "", output.Error("Failed to parse release: %w", err)
	}

	return release.TagName, nil
}

func downloadAndInstall(version string) error {
	tempDir, err := os.MkdirTemp("", "anytype-update-*")
	if err != nil {
		return output.Error("Failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	archivePath := filepath.Join(tempDir, getArchiveName(version))
	if err := downloadRelease(version, archivePath); err != nil {
		return err
	}

	binaryName := "anytype"
	if runtime.GOOS == "windows" {
		binaryName = "anytype.exe"
	}

	binaryReader, err := extractBinary(archivePath, binaryName)
	if err != nil {
		return output.Error("Failed to extract: %w", err)
	}
	defer binaryReader.Close()

	if err := selfupdate.Apply(binaryReader, selfupdate.Options{}); err != nil {
		return output.Error("failed to apply update: %w", err)
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
	output.Info("Downloading %s...", archiveName)

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return downloadViaAPI(version, archiveName, destination)
	}

	url := fmt.Sprintf("%s/releases/download/%s/%s",
		config.GitHubBaseURL, version, archiveName)

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
		return output.Error("Failed to parse release: %w", err)
	}

	var assetURL string
	for _, asset := range release.Assets {
		if asset.Name == filename {
			assetURL = asset.URL
			break
		}
	}
	if assetURL == "" {
		return output.Error("release asset %s not found", filename)
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
		return output.Error("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractBinary(archivePath, binaryName string) (io.ReadCloser, error) {
	if strings.HasSuffix(archivePath, ".zip") {
		return extractBinaryFromZip(archivePath, binaryName)
	}
	return extractBinaryFromTarGz(archivePath, binaryName)
}

func extractBinaryFromTarGz(archivePath, binaryName string) (io.ReadCloser, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}

	gz, err := gzip.NewReader(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			file.Close()
			gz.Close()
			return nil, err
		}

		if header.Typeflag == tar.TypeReg && header.Name == binaryName {
			return &tarGzReader{Reader: tr, gz: gz, file: file}, nil
		}
	}

	file.Close()
	gz.Close()
	return nil, fmt.Errorf("binary %s not found in archive", binaryName)
}

type tarGzReader struct {
	io.Reader
	gz   *gzip.Reader
	file *os.File
}

func (r *tarGzReader) Close() error {
	r.gz.Close()
	return r.file.Close()
}

func extractBinaryFromZip(archivePath, binaryName string) (io.ReadCloser, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}

	for _, f := range r.File {
		if f.Name == binaryName {
			rc, err := f.Open()
			if err != nil {
				r.Close()
				return nil, err
			}
			return &zipReader{ReadCloser: rc, zipReader: r}, nil
		}
	}

	r.Close()
	return nil, fmt.Errorf("binary %s not found in archive", binaryName)
}

type zipReader struct {
	io.ReadCloser
	zipReader *zip.ReadCloser
}

func (r *zipReader) Close() error {
	r.ReadCloser.Close()
	return r.zipReader.Close()
}

func githubAPI(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := config.GitHubAPIURL + endpoint

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
	return output.Error("API error %d: %s", resp.StatusCode, string(body))
}
