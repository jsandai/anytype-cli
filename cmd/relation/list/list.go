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
		Short: "List relations",
		Long: `List all relations available in a space.

Examples:
  anytype relation list --space <space-id>
  anytype relation list --space <space-id> --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if spaceID == "" {
				return fmt.Errorf("--space is required")
			}

			relations, err := core.ListRelations(spaceID)
			if err != nil {
				return output.Error("Failed to list relations: %w", err)
			}

			if len(relations) == 0 {
				output.Info("No relations found")
				return nil
			}

			if jsonOut {
				out, _ := json.MarshalIndent(relations, "", "  ")
				output.Print(string(out))
			} else {
				output.Info("Relations (%d):\n", len(relations))
				for _, r := range relations {
					output.Print("  %s (%s)", r.Name, r.Format)
					output.Print("    Key: %s", r.UniqueKey)
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
