package list

import (
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available spaces",
		Long:  "List all spaces available in your account",
		RunE: func(cmd *cobra.Command, args []string) error {
			spaces, err := core.ListSpaces()
			if err != nil {
				return output.Error("failed to list spaces: %w", err)
			}

			if len(spaces) == 0 {
				output.Info("No spaces found")
				return nil
			}

			output.Info("%-75s %-30s %s", "SPACE ID", "NAME", "STATUS")
			output.Info("%-75s %-30s %s", "────────", "────", "──────")

			for _, space := range spaces {
				status := "Active"
				if space.Status == 0 {
					status = "Unknown"
				}

				name := space.Name
				if len(name) > 28 {
					name = name[:25] + "..."
				}

				output.Info("%-75s %-30s %s", space.SpaceId, name, status)
			}

			return nil
		},
	}

	return cmd
}
