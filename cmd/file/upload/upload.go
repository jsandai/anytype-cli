package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func NewUploadCmd() *cobra.Command {
	var (
		fileType string
		fromURL  bool
	)

	cmd := &cobra.Command{
		Use:   "upload <space-id> <path-or-url>",
		Short: "Upload a file to a space",
		Long: `Upload a local file or fetch from URL and add it to an Anytype space.

Returns the object ID of the uploaded file, which can be used as an attachment
in chat messages.

Examples:
  anytype file upload bafyrei... /path/to/image.png
  anytype file upload bafyrei... https://example.com/image.png --url
  anytype file upload bafyrei... /path/to/file.pdf --type pdf`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			spaceId := args[0]
			pathOrURL := args[1]

			// Determine file type
			var fType model.BlockContentFileType
			switch fileType {
			case "image":
				fType = model.BlockContentFile_Image
			case "audio":
				fType = model.BlockContentFile_Audio
			case "video":
				fType = model.BlockContentFile_Video
			case "pdf":
				fType = model.BlockContentFile_PDF
			case "file", "":
				if fileType == "" {
					fType = core.DetectFileType(pathOrURL)
				} else {
					fType = model.BlockContentFile_File
				}
			default:
				return fmt.Errorf("unknown file type: %s (use: image, audio, video, pdf, file)", fileType)
			}

			var result *core.FileUploadResult
			var err error

			if fromURL {
				result, err = core.UploadFileFromURL(spaceId, pathOrURL, fType)
			} else {
				path := pathOrURL
				if len(path) >= 2 && strings.HasPrefix(path, "~/") {
					home, _ := os.UserHomeDir()
					path = filepath.Join(home, path[2:])
				}
				absPath, absErr := filepath.Abs(path)
				if absErr != nil {
					return fmt.Errorf("failed to resolve path: %w", absErr)
				}

				if _, statErr := os.Stat(absPath); os.IsNotExist(statErr) {
					return fmt.Errorf("file not found: %s", absPath)
				}

				result, err = core.UploadFile(spaceId, absPath, fType)
			}

			if err != nil {
				return err
			}

			fmt.Printf("Uploaded: %s\n", result.ObjectId)
			return nil
		},
	}

	cmd.Flags().StringVarP(&fileType, "type", "t", "", "File type: image, audio, video, pdf, file (auto-detected if not specified)")
	cmd.Flags().BoolVar(&fromURL, "url", false, "Treat the path as a URL to fetch from")

	return cmd
}
