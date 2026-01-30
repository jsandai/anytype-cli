package file

import (
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/cmd/file/download"
	"github.com/anyproto/anytype-cli/cmd/file/upload"
)

func NewFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file",
		Short: "File operations",
		Long:  "Upload, download, and manage files in Anytype spaces",
	}

	cmd.AddCommand(download.NewDownloadCmd())
	cmd.AddCommand(upload.NewUploadCmd())

	return cmd
}
