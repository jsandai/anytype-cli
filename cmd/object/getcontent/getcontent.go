package getcontent

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewGetContentCmd() *cobra.Command {
	var (
		spaceID  string
		outFile  string
	)

	cmd := &cobra.Command{
		Use:   "get-content <object-id>",
		Short: "Get object content as markdown",
		Long: `Get the content of an object and output it as markdown/text.
The blocks are converted back to a readable format.

Examples:
  # Print content to stdout
  anytype object get-content <object-id> --space <space-id>

  # Save content to a file
  anytype object get-content <object-id> --space <space-id> --output ./note.md`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			objectID := args[0]

			if spaceID == "" {
				return fmt.Errorf("--space is required")
			}

			content, err := core.GetObjectContent(spaceID, objectID)
			if err != nil {
				return output.Error("Failed to get content: %w", err)
			}

			if outFile != "" {
				err := os.WriteFile(outFile, []byte(content), 0644)
				if err != nil {
					return output.Error("Failed to write file: %w", err)
				}
				output.Success("Content saved to %s", outFile)
			} else {
				fmt.Println(content)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&spaceID, "space", "", "Space ID (required)")
	cmd.Flags().StringVar(&outFile, "output", "", "Output file path (optional, defaults to stdout)")
	cmd.Flags().StringVarP(&outFile, "o", "o", "", "Output file path (short form)")

	return cmd
}
