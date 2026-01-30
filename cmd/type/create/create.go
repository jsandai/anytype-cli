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
		spaceID     string
		name        string
		pluralName  string
		description string
		icon        string
		jsonOut     bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a custom object type",
		Long: `Create a new custom object type in the specified space.

Examples:
  # Create a Memory type
  anytype type create --space <space-id> --name "Memory" --icon "🧠"

  # With description
  anytype type create --space <space-id> --name "Memory" \
    --plural "Memories" --description "Agent memory entries" --icon "🧠"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if spaceID == "" {
				return fmt.Errorf("--space is required")
			}
			if name == "" {
				return fmt.Errorf("--name is required")
			}

			req := core.CreateTypeRequest{
				SpaceID:     spaceID,
				Name:        name,
				PluralName:  pluralName,
				Description: description,
				IconEmoji:   icon,
			}

			result, err := core.CreateObjectType(req)
			if err != nil {
				return output.Error("Failed to create type: %w", err)
			}

			if jsonOut {
				out, _ := json.MarshalIndent(result, "", "  ")
				output.Print(string(out))
			} else {
				output.Success("Created type: %s", result.Name)
				output.Info("  ID:        %s", result.ID)
				output.Info("  UniqueKey: %s", result.UniqueKey)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&spaceID, "space", "", "Space ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "Type name (required)")
	cmd.Flags().StringVar(&pluralName, "plural", "", "Plural name (e.g., 'Memories')")
	cmd.Flags().StringVar(&description, "description", "", "Type description")
	cmd.Flags().StringVar(&icon, "icon", "", "Emoji icon (e.g., '🧠')")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")

	return cmd
}
