package apikeys

import (
	"github.com/spf13/cobra"

	apiKeysCreateCmd "github.com/anyproto/anytype-cli/cmd/auth/apikeys/create"
	apiKeysListCmd "github.com/anyproto/anytype-cli/cmd/auth/apikeys/list"
	apiKeysRevokeCmd "github.com/anyproto/anytype-cli/cmd/auth/apikeys/revoke"
)

// NewApiKeysCmd creates the auth apikeys command
func NewApiKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apikeys <command>",
		Short: "Manage API keys",
		Long:  "Create, list, and revoke API keys for programmatic access",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(apiKeysCreateCmd.NewCreateCmd())
	cmd.AddCommand(apiKeysListCmd.NewListCmd())
	cmd.AddCommand(apiKeysRevokeCmd.NewRevokeCmd())

	return cmd
}
