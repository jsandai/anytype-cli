package setcontent

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewSetContentCmd() *cobra.Command {
	var (
		spaceID  string
		file     string
		text     string
		jsonOut  bool
	)

	cmd := &cobra.Command{
		Use:   "set-content <object-id>",
		Short: "Set object body content from markdown",
		Long: `Set the body content of an object from markdown text or file.

The markdown is converted to Anytype blocks (headings, lists, code, etc.)
using Anytype's native converter.

Examples:
  # Set content from a markdown file
  anytype object set-content <object-id> --space <space-id> --file ./note.md

  # Set content from inline text
  anytype object set-content <object-id> --space <space-id> --text "# Hello\n\nThis is content."

  # Pipe markdown content
  cat note.md | anytype object set-content <object-id> --space <space-id> --file -`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			objectID := args[0]

			if spaceID == "" {
				return fmt.Errorf("--space is required")
			}

			if file == "" && text == "" {
				return fmt.Errorf("either --file or --text is required")
			}

			var markdown string

			if file != "" {
				// Read from file or stdin
				if file == "-" {
					data, err := io.ReadAll(os.Stdin)
					if err != nil {
						return fmt.Errorf("failed to read stdin: %w", err)
					}
					markdown = string(data)
				} else {
					data, err := os.ReadFile(file)
					if err != nil {
						return fmt.Errorf("failed to read file: %w", err)
					}
					markdown = string(data)
				}
			} else {
				markdown = text
			}

			if markdown == "" {
				return fmt.Errorf("content is empty")
			}

			req := core.SetObjectContentRequest{
				SpaceID:  spaceID,
				ObjectID: objectID,
				Markdown: markdown,
			}

			result, err := core.SetObjectContent(req)
			if err != nil {
				return output.Error("Failed to set content: %w", err)
			}

			if jsonOut {
				out, _ := json.MarshalIndent(result, "", "  ")
				output.Print(string(out))
			} else {
				output.Success("Set content on object: %s", objectID)
				output.Info("Created %d block(s)", len(result.BlockIDs))
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&spaceID, "space", "", "Space ID (required)")
	cmd.Flags().StringVar(&file, "file", "", "Path to markdown file (use - for stdin)")
	cmd.Flags().StringVar(&text, "text", "", "Markdown text content")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")

	return cmd
}
