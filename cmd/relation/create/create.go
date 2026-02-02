package create

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewCreateCmd() *cobra.Command {
	var (
		spaceID     string
		name        string
		format      string
		description string
		jsonOut     bool
	)

	// Build format list for help text
	var formats []string
	for f := range core.RelationFormats {
		formats = append(formats, f)
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a custom relation",
		Long: fmt.Sprintf(`Create a new custom relation (field/property) in the specified space.

Available formats: %s

Examples:
  # Create a text relation
  anytype relation create --space <space-id> --name "Notes" --format text

  # Create a select relation for memory types
  anytype relation create --space <space-id> --name "memoryKind" --format select \
    --description "Type of memory entry"

  # Create a date relation
  anytype relation create --space <space-id> --name "recordedAt" --format date`, strings.Join(formats, ", ")),
		RunE: func(cmd *cobra.Command, args []string) error {
			if spaceID == "" {
				return fmt.Errorf("--space is required")
			}
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if format == "" {
				return fmt.Errorf("--format is required")
			}

			req := core.CreateRelationRequest{
				SpaceID:     spaceID,
				Name:        name,
				Format:      format,
				Description: description,
			}

			result, err := core.CreateRelation(req)
			if err != nil {
				return output.Error("Failed to create relation: %w", err)
			}

			if jsonOut {
				out, _ := json.MarshalIndent(result, "", "  ")
				output.Print(string(out))
			} else {
				output.Success("Created relation: %s", result.Name)
				output.Info("  ID:     %s", result.ID)
				output.Info("  Key:    %s", result.UniqueKey)
				output.Info("  Format: %s", result.Format)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&spaceID, "space", "", "Space ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "Relation name (required)")
	cmd.Flags().StringVar(&format, "format", "", "Relation format (required): text, number, select, multi-select, date, checkbox, url, email, phone, object, file")
	cmd.Flags().StringVar(&description, "description", "", "Relation description")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")

	return cmd
}
