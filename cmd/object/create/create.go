package create

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewCreateCmd() *cobra.Command {
	var (
		spaceID  string
		typeID   string
		name     string
		details  string
		jsonOut  bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new object",
		Long: `Create a new Anytype object in the specified space.

Examples:
  # Create a note
  anytype object create --space <space-id> --type ot-note --name "My Note"

  # Create with custom type and details
  anytype object create --space <space-id> --type <custom-type-id> --name "Memory" \
    --details '{"memoryKind":"daily","importance":"high"}'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if spaceID == "" {
				return fmt.Errorf("--space is required")
			}
			if typeID == "" {
				return fmt.Errorf("--type is required")
			}

			req := core.CreateObjectRequest{
				SpaceID: spaceID,
				TypeID:  typeID,
				Name:    name,
			}

			// Parse details JSON if provided
			if details != "" {
				var detailsMap map[string]interface{}
				if err := json.Unmarshal([]byte(details), &detailsMap); err != nil {
					return fmt.Errorf("invalid --details JSON: %w", err)
				}
				req.Details = detailsMap
			}

			result, err := core.CreateObject(req)
			if err != nil {
				return output.Error("Failed to create object: %w", err)
			}

			if jsonOut {
				out, _ := json.MarshalIndent(result, "", "  ")
				output.Print(string(out))
			} else {
				output.Success("Created object: %s", result.ID)
				if result.Name != "" {
					output.Info("Name: %s", result.Name)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&spaceID, "space", "", "Space ID (required)")
	cmd.Flags().StringVar(&typeID, "type", "", "Object type unique key (e.g., ot-note, ot-page, or custom type ID) (required)")
	cmd.Flags().StringVar(&name, "name", "", "Object name")
	cmd.Flags().StringVar(&details, "details", "", "Additional details as JSON object")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")

	return cmd
}
