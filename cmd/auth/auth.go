package auth

import (
	"github.com/spf13/cobra"

	authApiKeysCmd "github.com/anyproto/anytype-cli/cmd/auth/apikeys"
	authCreateCmd "github.com/anyproto/anytype-cli/cmd/auth/create"
	authLoginCmd "github.com/anyproto/anytype-cli/cmd/auth/login"
	authLogoutCmd "github.com/anyproto/anytype-cli/cmd/auth/logout"
)

func NewAuthCmd() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate with Anytype",
	}

	authCmd.AddCommand(authLoginCmd.NewLoginCmd())
	authCmd.AddCommand(authLogoutCmd.NewLogoutCmd())
	authCmd.AddCommand(authCreateCmd.NewCreateCmd())
	authCmd.AddCommand(authApiKeysCmd.NewApiKeysCmd())

	return authCmd
}
