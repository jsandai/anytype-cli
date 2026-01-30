package core

import (
	"testing"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func TestDetectAttachmentType(t *testing.T) {
	tests := []struct {
		name     string
		fileType model.BlockContentFileType
		expected model.ChatMessageAttachmentAttachmentType
	}{
		{
			name:     "image file type maps to image attachment",
			fileType: model.BlockContentFile_Image,
			expected: model.ChatMessageAttachment_IMAGE,
		},
		{
			name:     "audio file type maps to file attachment",
			fileType: model.BlockContentFile_Audio,
			expected: model.ChatMessageAttachment_FILE,
		},
		{
			name:     "video file type maps to file attachment",
			fileType: model.BlockContentFile_Video,
			expected: model.ChatMessageAttachment_FILE,
		},
		{
			name:     "pdf file type maps to file attachment",
			fileType: model.BlockContentFile_PDF,
			expected: model.ChatMessageAttachment_FILE,
		},
		{
			name:     "generic file type maps to file attachment",
			fileType: model.BlockContentFile_File,
			expected: model.ChatMessageAttachment_FILE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectAttachmentType(tt.fileType)
			if result != tt.expected {
				t.Errorf("DetectAttachmentType(%v) = %v, want %v", tt.fileType, result, tt.expected)
			}
		})
	}
}
