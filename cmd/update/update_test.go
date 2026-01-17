package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/minio/selfupdate"
)

func TestExtractBinaryFromTarGz(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.tar.gz")
	binaryContent := []byte("fake binary content")
	binaryName := "anytype"

	f, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("failed to create archive: %v", err)
	}

	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)

	hdr := &tar.Header{
		Name: binaryName,
		Mode: 0755,
		Size: int64(len(binaryContent)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}
	if _, err := tw.Write(binaryContent); err != nil {
		t.Fatalf("failed to write tar content: %v", err)
	}

	tw.Close()
	gw.Close()
	f.Close()

	reader, err := extractBinaryFromTarGz(archivePath, binaryName)
	if err != nil {
		t.Fatalf("extractBinaryFromTarGz failed: %v", err)
	}
	defer reader.Close()

	extracted, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read extracted binary: %v", err)
	}

	if !bytes.Equal(extracted, binaryContent) {
		t.Errorf("extracted content mismatch: got %q, want %q", extracted, binaryContent)
	}
}

func TestExtractBinaryFromZip(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")
	binaryContent := []byte("fake binary content")
	binaryName := "anytype.exe"

	f, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("failed to create archive: %v", err)
	}

	zw := zip.NewWriter(f)
	w, err := zw.Create(binaryName)
	if err != nil {
		t.Fatalf("failed to create zip entry: %v", err)
	}
	if _, err := w.Write(binaryContent); err != nil {
		t.Fatalf("failed to write zip content: %v", err)
	}

	zw.Close()
	f.Close()

	reader, err := extractBinaryFromZip(archivePath, binaryName)
	if err != nil {
		t.Fatalf("extractBinaryFromZip failed: %v", err)
	}
	defer reader.Close()

	extracted, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read extracted binary: %v", err)
	}

	if !bytes.Equal(extracted, binaryContent) {
		t.Errorf("extracted content mismatch: got %q, want %q", extracted, binaryContent)
	}
}

func TestExtractBinaryNotFound(t *testing.T) {
	tempDir := t.TempDir()

	tarPath := filepath.Join(tempDir, "test.tar.gz")
	f, _ := os.Create(tarPath)
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	tw.Close()
	gw.Close()
	f.Close()

	_, err := extractBinaryFromTarGz(tarPath, "nonexistent")
	if err == nil {
		t.Error("expected error for missing binary in tar.gz")
	}

	zipPath := filepath.Join(tempDir, "test.zip")
	f, _ = os.Create(zipPath)
	zw := zip.NewWriter(f)
	zw.Close()
	f.Close()

	_, err = extractBinaryFromZip(zipPath, "nonexistent")
	if err == nil {
		t.Error("expected error for missing binary in zip")
	}
}

func TestSelfUpdateApply(t *testing.T) {
	tempDir := t.TempDir()

	currentBinary := filepath.Join(tempDir, "current")
	if runtime.GOOS == "windows" {
		currentBinary += ".exe"
	}

	oldContent := []byte("old version")
	if err := os.WriteFile(currentBinary, oldContent, 0755); err != nil {
		t.Fatalf("failed to write current binary: %v", err)
	}

	newContent := []byte("new version")

	err := selfupdate.Apply(bytes.NewReader(newContent), selfupdate.Options{
		TargetPath: currentBinary,
	})
	if err != nil {
		t.Fatalf("selfupdate.Apply failed: %v", err)
	}

	updatedContent, err := os.ReadFile(currentBinary)
	if err != nil {
		t.Fatalf("failed to read updated binary: %v", err)
	}

	if !bytes.Equal(updatedContent, newContent) {
		t.Errorf("update failed: got %q, want %q", updatedContent, newContent)
	}

	if runtime.GOOS == "windows" {
		os.Remove(currentBinary + ".old")
	}
}

func TestGetArchiveName(t *testing.T) {
	version := "v1.0.0"
	name := getArchiveName(version)
	expectedBase := "anytype-cli-v1.0.0-" + runtime.GOOS + "-" + runtime.GOARCH

	if runtime.GOOS == "windows" {
		expected := expectedBase + ".zip"
		if name != expected {
			t.Errorf("getArchiveName() = %q, want %q", name, expected)
		}
	} else {
		expected := expectedBase + ".tar.gz"
		if name != expected {
			t.Errorf("getArchiveName() = %q, want %q", name, expected)
		}
	}
}
