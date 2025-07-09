package list

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all API keys",
		Long:  "List all API keys associated with your account",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement API key listing
			// 1. Ensure user is authenticated
			// 2. Fetch API keys from server
			// 3. Display in table format
			return fmt.Errorf("API key listing not yet implemented")
		},
	}

	return cmd
}
