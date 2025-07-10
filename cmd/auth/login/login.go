package login

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/daemon"
	"github.com/anyproto/anytype-cli/internal"
)

func NewLoginCmd() *cobra.Command {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to your Anytype vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			mnemonic, _ := cmd.Flags().GetString("mnemonic")
			rootPath, _ := cmd.Flags().GetString("path")
			apiAddr, _ := cmd.Flags().GetString("api-addr")

			statusResp, err := daemon.SendTaskStatus("server")
			if err != nil || statusResp.Status != "running" {
				return fmt.Errorf("server is not running")
			}

			if err := internal.Login(mnemonic, rootPath, apiAddr); err != nil {
				return fmt.Errorf("✗ Failed to log in: %w", err)
			}
			fmt.Println("✓ Successfully logged in")
			return nil

		},
	}

	loginCmd.Flags().String("mnemonic", "", "Provide mnemonic (12 words) for authentication")
	loginCmd.Flags().String("path", "", "Provide custom root path for wallet recovery")
	loginCmd.Flags().String("api-addr", "", "API listen address (default: 127.0.0.1:31009)")

	return loginCmd
}
