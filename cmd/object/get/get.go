package get

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewGetCmd() *cobra.Command {
	var (
		spaceID string
		jsonOut bool
	)

	cmd := &cobra.Command{
		Use:   "get <object-id>",
		Short: "Get object details",
		Long: `Retrieve full details of an Anytype object.

Examples:
  anytype object get <object-id> --space <space-id>
  anytype object get <object-id> --space <space-id> --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			objectID := args[0]

			if spaceID == "" {
				return fmt.Errorf("--space is required")
			}

			result, err := core.GetObject(spaceID, objectID)
			if err != nil {
				return output.Error("Failed to get object: %w", err)
			}

			if result == nil {
				return output.Error("Object not found")
			}

			if jsonOut {
				out, _ := json.MarshalIndent(result, "", "  ")
				output.Print(string(out))
			} else {
				output.Info("Object: %s", result.ID)
				output.Print("  Name:  %s", result.Name)
				output.Print("  Type:  %s", result.Type)
				output.Print("  Space: %s", result.SpaceID)
				
				if len(result.Details) > 0 {
					output.Print("\nDetails:")
					for key, val := range result.Details {
						// Skip internal/common fields already shown
						if key == "id" || key == "name" || key == "type" || key == "spaceId" {
							continue
						}
						output.Print("  %s: %v", key, val)
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&spaceID, "space", "", "Space ID (required)")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")

	return cmd
}
