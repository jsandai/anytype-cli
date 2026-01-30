package update

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewUpdateCmd() *cobra.Command {
	var (
		spaceID string
		name    string
		details string
	)

	cmd := &cobra.Command{
		Use:   "update <object-id>",
		Short: "Update object details",
		Long: `Update an object's name or other details/relations.

Examples:
  # Update name
  anytype object update <object-id> --space <space-id> --name "New Name"

  # Update custom fields
  anytype object update <object-id> --space <space-id> \
    --details '{"importance":"high","memoryKind":"insight"}'`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			objectID := args[0]

			if spaceID == "" {
				return fmt.Errorf("--space is required")
			}

			updateDetails := make(map[string]interface{})

			// Add name if provided
			if name != "" {
				updateDetails["name"] = name
			}

			// Parse and merge details JSON if provided
			if details != "" {
				var detailsMap map[string]interface{}
				if err := json.Unmarshal([]byte(details), &detailsMap); err != nil {
					return fmt.Errorf("invalid --details JSON: %w", err)
				}
				for k, v := range detailsMap {
					updateDetails[k] = v
				}
			}

			if len(updateDetails) == 0 {
				return fmt.Errorf("nothing to update: provide --name or --details")
			}

			err := core.UpdateObjectDetails(spaceID, objectID, updateDetails)
			if err != nil {
				return output.Error("Failed to update object: %w", err)
			}

			output.Success("Updated object: %s", objectID)
			return nil
		},
	}

	cmd.Flags().StringVar(&spaceID, "space", "", "Space ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "New object name")
	cmd.Flags().StringVar(&details, "details", "", "Details to update as JSON object")

	return cmd
}
