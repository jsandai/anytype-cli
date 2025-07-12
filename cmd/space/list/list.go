package list

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available spaces",
		Long:  "List all spaces available in your account",
		RunE: func(cmd *cobra.Command, args []string) error {
			spaces, err := core.ListSpaces()
			if err != nil {
				return fmt.Errorf("failed to list spaces: %w", err)
			}

			if len(spaces) == 0 {
				fmt.Println("No spaces found")
				return nil
			}

			fmt.Printf("%-75s %-30s %s\n", "SPACE ID", "NAME", "STATUS")
			fmt.Printf("%-75s %-30s %s\n", "────────", "────", "──────")

			for _, space := range spaces {
				status := "Active"
				if space.Status == 0 {
					status = "Unknown"
				}

				name := space.Name
				if len(name) > 28 {
					name = name[:25] + "..."
				}

				fmt.Printf("%-75s %-30s %s\n", space.SpaceID, name, status)
			}

			return nil
		},
	}

	return cmd
}
