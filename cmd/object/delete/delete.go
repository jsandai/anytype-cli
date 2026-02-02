package delete

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewDeleteCmd() *cobra.Command {
	var (
		force bool
	)

	cmd := &cobra.Command{
		Use:   "delete <object-id> [object-id...]",
		Short: "Delete objects",
		Long: `Delete one or more Anytype objects.

Examples:
  # Delete a single object
  anytype object delete <object-id>

  # Delete multiple objects
  anytype object delete <id1> <id2> <id3>

  # Skip confirmation
  anytype object delete <object-id> --force`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			objectIDs := args

			if !force {
				output.Info("About to delete %d object(s):", len(objectIDs))
				for _, id := range objectIDs {
					output.Print("  - %s", id)
				}
				output.Print("")

				// Simple confirmation
				var confirm string
				fmt.Print("Type 'yes' to confirm: ")
				fmt.Scanln(&confirm)
				if confirm != "yes" {
					output.Info("Cancelled")
					return nil
				}
			}

			err := core.DeleteObjects(objectIDs)
			if err != nil {
				return output.Error("Failed to delete: %w", err)
			}

			output.Success("Deleted %d object(s)", len(objectIDs))
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}
