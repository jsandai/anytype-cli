package revoke

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRevokeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke <key-id>",
		Short: "Revoke an API key",
		Long:  "Revoke an API key by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyID := args[0]
			// TODO: Implement API key revocation
			// 1. Ensure user is authenticated
			// 2. Revoke key on server
			// 3. Confirm revocation
			_ = keyID
			return fmt.Errorf("API key revocation not yet implemented")
		},
	}

	return cmd
}
