package status

import (
	"context"
	"fmt"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/config"
)

func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			hasMnemonic := false
			if _, err := core.GetStoredMnemonic(); err == nil {
				hasMnemonic = true
			}

			hasToken := false
			token := ""
			if t, err := core.GetStoredToken(); err == nil {
				hasToken = true
				token = t
			}

			// Get account info from config
			configMgr := config.GetConfigManager()
			_ = configMgr.Load()
			cfg := configMgr.Get()
			accountID := cfg.AccountID

			// First check if server is running
			serverRunning := false
			err := core.GRPCCallNoAuth(func(ctx context.Context, client service.ClientCommandsClient) error {
				_, err := client.AppGetVersion(ctx, &pb.RpcAppGetVersionRequest{})
				return err
			})
			serverRunning = err == nil

			// If server is running and we have a token, we're logged in
			// (server auto-logs in on restart using stored mnemonic)
			isLoggedIn := serverRunning && hasToken

			// Display status based on priority: server -> credentials -> login
			if !serverRunning {
				fmt.Println("Server is not running. Run 'anytype serve' to start the server.")
				if hasMnemonic || hasToken || accountID != "" {
					fmt.Println("Credentials are stored in keychain.")
				}
				return nil
			}

			if !hasMnemonic && !hasToken && accountID == "" {
				fmt.Println("Not authenticated. Run 'anytype auth login' to authenticate.")
				return nil
			}

			fmt.Println("Anytype")

			if isLoggedIn && accountID != "" {
				fmt.Printf("  ✓ Logged in to account %s (keychain)\n", accountID)
			} else if hasToken || hasMnemonic {
				fmt.Println("  ✗ Not logged in (credentials stored in keychain)")
				if !isLoggedIn && hasToken {
					fmt.Println("    Note: Server is not running or session expired. Run 'anytype serve' to start server.")
				}
			} else {
				fmt.Println("  ✗ Not logged in")
			}

			fmt.Printf("  - Active session: %v\n", isLoggedIn)

			if hasMnemonic {
				fmt.Println("  - Mnemonic: stored")
			}

			if hasToken {
				if len(token) > 8 {
					fmt.Printf("  - Token: %s****\n", token[:8])
				} else {
					fmt.Println("  - Token: stored")
				}
			}

			return nil
		},
	}

	return cmd
}
