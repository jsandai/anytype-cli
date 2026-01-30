package core

import (
	"testing"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func TestDetectFileType(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected model.BlockContentFileType
	}{
		// Images
		{name: "jpg image", path: "/path/to/image.jpg", expected: model.BlockContentFile_Image},
		{name: "jpeg image", path: "/path/to/image.jpeg", expected: model.BlockContentFile_Image},
		{name: "png image", path: "/path/to/image.png", expected: model.BlockContentFile_Image},
		{name: "gif image", path: "/path/to/image.gif", expected: model.BlockContentFile_Image},
		{name: "webp image", path: "/path/to/image.webp", expected: model.BlockContentFile_Image},
		{name: "svg image", path: "/path/to/image.svg", expected: model.BlockContentFile_Image},
		{name: "bmp image", path: "/path/to/image.bmp", expected: model.BlockContentFile_Image},
		{name: "ico image", path: "/path/to/favicon.ico", expected: model.BlockContentFile_Image},

		// Audio
		{name: "mp3 audio", path: "/path/to/song.mp3", expected: model.BlockContentFile_Audio},
		{name: "wav audio", path: "/path/to/sound.wav", expected: model.BlockContentFile_Audio},
		{name: "ogg audio", path: "/path/to/audio.ogg", expected: model.BlockContentFile_Audio},
		{name: "m4a audio", path: "/path/to/audio.m4a", expected: model.BlockContentFile_Audio},
		{name: "flac audio", path: "/path/to/audio.flac", expected: model.BlockContentFile_Audio},
		{name: "aac audio", path: "/path/to/audio.aac", expected: model.BlockContentFile_Audio},

		// Video
		{name: "mp4 video", path: "/path/to/video.mp4", expected: model.BlockContentFile_Video},
		{name: "mov video", path: "/path/to/video.mov", expected: model.BlockContentFile_Video},
		{name: "avi video", path: "/path/to/video.avi", expected: model.BlockContentFile_Video},
		{name: "mkv video", path: "/path/to/video.mkv", expected: model.BlockContentFile_Video},
		{name: "webm video", path: "/path/to/video.webm", expected: model.BlockContentFile_Video},

		// PDF
		{name: "pdf document", path: "/path/to/document.pdf", expected: model.BlockContentFile_PDF},

		// Generic files
		{name: "text file", path: "/path/to/file.txt", expected: model.BlockContentFile_File},
		{name: "zip archive", path: "/path/to/archive.zip", expected: model.BlockContentFile_File},
		{name: "json file", path: "/path/to/data.json", expected: model.BlockContentFile_File},
		{name: "no extension", path: "/path/to/README", expected: model.BlockContentFile_File},

		// Case insensitivity
		{name: "uppercase JPG", path: "/path/to/IMAGE.JPG", expected: model.BlockContentFile_Image},
		{name: "mixed case Png", path: "/path/to/Image.Png", expected: model.BlockContentFile_Image},
		{name: "uppercase PDF", path: "/path/to/DOC.PDF", expected: model.BlockContentFile_PDF},

		// Edge cases
		{name: "empty path", path: "", expected: model.BlockContentFile_File},
		{name: "dot file", path: ".gitignore", expected: model.BlockContentFile_File},
		{name: "double extension", path: "/path/to/file.tar.gz", expected: model.BlockContentFile_File},
		{name: "extension in path", path: "/path/with.jpg/file.txt", expected: model.BlockContentFile_File},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectFileType(tt.path)
			if result != tt.expected {
				t.Errorf("DetectFileType(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestHasImageExtension(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/path/to/image.jpg", true},
		{"/path/to/image.PNG", true},
		{"/path/to/file.txt", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := hasImageExtension(tt.path)
			if result != tt.expected {
				t.Errorf("hasImageExtension(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestHasAudioExtension(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/path/to/song.mp3", true},
		{"/path/to/audio.FLAC", true},
		{"/path/to/file.txt", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := hasAudioExtension(tt.path)
			if result != tt.expected {
				t.Errorf("hasAudioExtension(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestHasVideoExtension(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/path/to/video.mp4", true},
		{"/path/to/movie.MKV", true},
		{"/path/to/file.txt", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := hasVideoExtension(tt.path)
			if result != tt.expected {
				t.Errorf("hasVideoExtension(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestHasPDFExtension(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/path/to/doc.pdf", true},
		{"/path/to/DOC.PDF", true},
		{"/path/to/file.txt", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := hasPDFExtension(tt.path)
			if result != tt.expected {
				t.Errorf("hasPDFExtension(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}
