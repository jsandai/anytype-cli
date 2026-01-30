package setcontent

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewSetContentCmd() *cobra.Command {
	var (
		spaceID  string
		filePath string
		text     string
	)

	cmd := &cobra.Command{
		Use:   "set-content <object-id>",
		Short: "Set object content from markdown",
		Long: `Set the content of an object using markdown text or a markdown file.
The markdown is converted to Anytype blocks using the native anymark converter.

Examples:
  # Set content from a file
  anytype object set-content <object-id> --space <space-id> --file ./note.md

  # Set content from text
  anytype object set-content <object-id> --space <space-id> --text "# Hello\n\nContent here"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			objectID := args[0]

			if spaceID == "" {
				return fmt.Errorf("--space is required")
			}

			if filePath == "" && text == "" {
				return fmt.Errorf("either --file or --text is required")
			}

			if filePath != "" && text != "" {
				return fmt.Errorf("cannot use both --file and --text")
			}

			var markdown string

			if filePath != "" {
				content, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("failed to read file: %w", err)
				}
				markdown = string(content)
			} else {
				markdown = text
			}

			err := core.SetObjectContent(spaceID, objectID, markdown)
			if err != nil {
				return output.Error("Failed to set content: %w", err)
			}

			output.Success("Content set successfully for object %s", objectID)
			return nil
		},
	}

	cmd.Flags().StringVar(&spaceID, "space", "", "Space ID (required)")
	cmd.Flags().StringVar(&filePath, "file", "", "Path to markdown file")
	cmd.Flags().StringVar(&text, "text", "", "Markdown text content")

	return cmd
}
