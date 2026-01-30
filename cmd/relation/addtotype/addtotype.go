package addtotype

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewAddToTypeCmd() *cobra.Command {
	var (
		spaceID     string
		typeKey     string
		relationKey string
	)

	cmd := &cobra.Command{
		Use:   "add-to-type",
		Short: "Add a relation to an object type",
		Long: `Add a relation to an object type's recommended relations.
This makes the relation appear by default on objects of that type.

Examples:
  # Add a relation to a type by unique keys
  anytype relation add-to-type --space <space-id> --type ot-note --relation myCustomRelation

  # Add multiple relations (run multiple times)
  anytype relation add-to-type --space <space-id> --type ot-697cd4deb3ddd0b369d32e0a --relation memoryDate
  anytype relation add-to-type --space <space-id> --type ot-697cd4deb3ddd0b369d32e0a --relation memoryTags`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if spaceID == "" {
				return fmt.Errorf("--space is required")
			}
			if typeKey == "" {
				return fmt.Errorf("--type is required")
			}
			if relationKey == "" {
				return fmt.Errorf("--relation is required")
			}

			err := core.AddRelationToType(spaceID, typeKey, relationKey)
			if err != nil {
				return output.Error("Failed to add relation to type: %w", err)
			}

			output.Success("Added relation '%s' to type '%s'", relationKey, typeKey)
			return nil
		},
	}

	cmd.Flags().StringVar(&spaceID, "space", "", "Space ID (required)")
	cmd.Flags().StringVar(&typeKey, "type", "", "Object type unique key (required)")
	cmd.Flags().StringVar(&relationKey, "relation", "", "Relation unique key to add (required)")

	return cmd
}
