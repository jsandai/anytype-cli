package download

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download <object-id> [destination]",
		Short: "Download a file from Anytype",
		Long: `Download a file object from Anytype to a local path.

If destination is not specified, downloads to the current directory.
If destination is a directory, uses the original filename.

Examples:
  anytype file download bafyrei... 
  anytype file download bafyrei... /tmp/myfile.png
  anytype file download bafyrei... ~/Downloads/`,
		Args: cobra.RangeArgs(1, 2),
		RunE: runDownload,
	}

	return cmd
}

func runDownload(cmd *cobra.Command, args []string) error {
	objectId := args[0]
	
	destPath := "."
	if len(args) > 1 {
		destPath = args[1]
	}
	
	// Expand ~ 
	if strings.HasPrefix(destPath, "~/") {
		home, _ := os.UserHomeDir()
		destPath = filepath.Join(home, destPath[2:])
	}
	
	// Make absolute
	absPath, err := filepath.Abs(destPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	result, err := core.DownloadFile(objectId, absPath)
	if err != nil {
		return err
	}

	output.Info("Downloaded: %s", result.LocalPath)
	return nil
}
