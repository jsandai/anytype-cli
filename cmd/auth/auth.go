package auth

import (
	"github.com/spf13/cobra"

	authApiKeyCmd "github.com/anyproto/anytype-cli/cmd/auth/apikey"
	authCreateCmd "github.com/anyproto/anytype-cli/cmd/auth/create"
	authLoginCmd "github.com/anyproto/anytype-cli/cmd/auth/login"
	authLogoutCmd "github.com/anyproto/anytype-cli/cmd/auth/logout"
	authStatusCmd "github.com/anyproto/anytype-cli/cmd/auth/status"
)

func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Manage authentication and accounts",
		Long:  "Create accounts, log in/out, check authentication status, and manage API keys for programmatic access.",
	}

	cmd.AddCommand(authLoginCmd.NewLoginCmd())
	cmd.AddCommand(authLogoutCmd.NewLogoutCmd())
	cmd.AddCommand(authStatusCmd.NewStatusCmd())
	cmd.AddCommand(authCreateCmd.NewCreateCmd())
	cmd.AddCommand(authApiKeyCmd.NewApiKeyCmd())

	return cmd
}
