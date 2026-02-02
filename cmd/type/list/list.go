package list

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewListCmd() *cobra.Command {
	var (
		spaceID string
		jsonOut bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List object types",
		Long: `List all object types available in a space.

Examples:
  anytype type list --space <space-id>
  anytype type list --space <space-id> --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if spaceID == "" {
				return fmt.Errorf("--space is required")
			}

			types, err := core.ListObjectTypes(spaceID)
			if err != nil {
				return output.Error("Failed to list types: %w", err)
			}

			if len(types) == 0 {
				output.Info("No types found")
				return nil
			}

			if jsonOut {
				out, _ := json.MarshalIndent(types, "", "  ")
				output.Print(string(out))
			} else {
				output.Info("Object Types (%d):\n", len(types))
				for _, t := range types {
					icon := t.IconEmoji
					if icon == "" {
						icon = "📄"
					}
					output.Print("  %s %s", icon, t.Name)
					output.Print("    Key: %s", t.UniqueKey)
					if t.Description != "" {
						output.Print("    %s", t.Description)
					}
					output.Print("")
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&spaceID, "space", "", "Space ID (required)")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")

	return cmd
}
