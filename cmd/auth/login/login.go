package login

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to your bot account",
		Long:  "Authenticate using your bot account key to access your Anytype bot account and stored data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			accountKey, _ := cmd.Flags().GetString("account-key")
			rootPath, _ := cmd.Flags().GetString("path")
			apiAddr, _ := cmd.Flags().GetString("api-addr")

			if err := core.LoginBot(accountKey, rootPath, apiAddr); err != nil {
				return output.Error("failed to log in: %w", err)
			}
			output.Success("Successfully logged in")
			return nil

		},
	}

	cmd.Flags().String("account-key", "", "Provide bot account key for authentication")
	cmd.Flags().String("path", "", "Provide custom root path for wallet recovery")
	cmd.Flags().String("api-addr", "", fmt.Sprintf("API listen address (default: %s)", config.DefaultAPIAddress))

	return cmd
}
