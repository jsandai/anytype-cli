package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

// FileUploadResult contains the result of a file upload
type FileUploadResult struct {
	ObjectId string
	Details  map[string]interface{}
}

// UploadFile uploads a file to Anytype and returns the object ID
func UploadFile(spaceId string, localPath string, fileType model.BlockContentFileType) (*FileUploadResult, error) {
	var result FileUploadResult
	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcFileUploadRequest{
			SpaceId:   spaceId,
			LocalPath: localPath,
			Type:      fileType,
		}
		resp, err := client.FileUpload(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to upload file: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcFileUploadResponseError_NULL {
			return fmt.Errorf("upload error: %s", resp.Error.Description)
		}
		result.ObjectId = resp.ObjectId
		return nil
	})
	return &result, err
}

// UploadFileFromURL uploads a file from a URL to Anytype
func UploadFileFromURL(spaceId string, url string, fileType model.BlockContentFileType) (*FileUploadResult, error) {
	var result FileUploadResult
	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcFileUploadRequest{
			SpaceId: spaceId,
			Url:     url,
			Type:    fileType,
		}
		resp, err := client.FileUpload(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to upload file from URL: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcFileUploadResponseError_NULL {
			return fmt.Errorf("upload error: %s", resp.Error.Description)
		}
		result.ObjectId = resp.ObjectId
		return nil
	})
	return &result, err
}

// DetectFileType attempts to determine the file type from extension
func DetectFileType(path string) model.BlockContentFileType {
	// Simple extension-based detection
	switch {
	case hasImageExtension(path):
		return model.BlockContentFile_Image
	case hasAudioExtension(path):
		return model.BlockContentFile_Audio
	case hasVideoExtension(path):
		return model.BlockContentFile_Video
	case hasPDFExtension(path):
		return model.BlockContentFile_PDF
	default:
		return model.BlockContentFile_File
	}
}

func hasImageExtension(path string) bool {
	exts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp", ".ico"}
	return hasExtension(path, exts)
}

func hasAudioExtension(path string) bool {
	exts := []string{".mp3", ".wav", ".ogg", ".m4a", ".flac", ".aac"}
	return hasExtension(path, exts)
}

func hasVideoExtension(path string) bool {
	exts := []string{".mp4", ".mov", ".avi", ".mkv", ".webm"}
	return hasExtension(path, exts)
}

func hasPDFExtension(path string) bool {
	exts := []string{".pdf"}
	return hasExtension(path, exts)
}

func hasExtension(path string, exts []string) bool {
	pathLower := strings.ToLower(path)
	for _, ext := range exts {
		if strings.HasSuffix(pathLower, ext) {
			return true
		}
	}
	return false
}

// FileDownloadResult contains the result of a file download
type FileDownloadResult struct {
	LocalPath string
}

// DownloadFile downloads a file from Anytype to a local path
func DownloadFile(objectId string, destPath string) (*FileDownloadResult, error) {
	var result FileDownloadResult
	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcFileDownloadRequest{
			ObjectId: objectId,
			Path:     destPath,
		}
		resp, err := client.FileDownload(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to download file: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcFileDownloadResponseError_NULL {
			return fmt.Errorf("download error: %s", resp.Error.Description)
		}
		result.LocalPath = resp.LocalPath
		return nil
	})
	return &result, err
}
