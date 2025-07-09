package create

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCreateCmd creates the auth create command
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Anytype account",
		Long:  "Create a new Anytype account with a generated mnemonic phrase",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement account creation logic
			// 1. Generate new mnemonic
			// 2. Create account on server
			// 3. Store credentials in keyring
			// 4. Display mnemonic to user (only time they'll see it)
			return fmt.Errorf("account creation not yet implemented")
		},
	}

	return cmd
}
