package create

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
	"github.com/spf13/cobra"
)

// NewCreateCmd creates the auth create command
func NewCreateCmd() *cobra.Command {
	var name string
	var rootPath string
	var apiAddr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new account",
		Long:  "Create a new Anytype account with a generated mnemonic phrase. The mnemonic is your master key for account recovery.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				output.Print("Enter account name: ")
				reader := bufio.NewReader(os.Stdin)
				name, _ = reader.ReadString('\n')
				name = strings.TrimSpace(name)

				if name == "" {
					return output.Error("account name is required")
				}
			}

			mnemonic, accountID, err := core.CreateWallet(name, rootPath, apiAddr)
			if err != nil {
				return output.Error("failed to create account: %w", err)
			}

			output.Success("Account created successfully!")

			output.Warning("IMPORTANT: Save your mnemonic phrase in a secure location.")
			output.Info("   This is the ONLY way to recover your account if you lose access.")

			words := strings.Split(mnemonic, " ")
			output.Print("")
			output.Print("╔════════════════════════════════════════════════════════╗")
			output.Print("║                    MNEMONIC PHRASE                     ║")
			output.Print("╠════════════════════════════════════════════════════════╣")
			output.Print("║  %-52s  ║", strings.Join(words[0:6], " "))
			output.Print("║  %-52s  ║", strings.Join(words[6:12], " "))
			output.Print("╚════════════════════════════════════════════════════════╝")

			output.Print("")
			output.Print("📋 Account Details:")
			output.Print("   Name: %s", name)
			output.Print("   ID: %s", accountID)

			output.Print("")
			output.Success("You are now logged in to your new account.")
			output.Success("Mnemonic saved to keychain.")

			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Account name")
	cmd.Flags().StringVar(&rootPath, "root-path", "", "Custom root path for storing account data")
	cmd.Flags().StringVar(&apiAddr, "api-addr", "", fmt.Sprintf("Custom API address (default: %s)", config.DefaultAPIAddress))

	return cmd
}
