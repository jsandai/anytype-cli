package create

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/internal"
)

func NewCreateCmd() *cobra.Command {
	var name string
	var expiresIn string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new API key",
		Long:  "Create a new API key for programmatic access to Anytype",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Using the existing token creation logic temporarily
			if err := internal.CreateToken(); err != nil {
				return fmt.Errorf("✗ Failed to create API key: %w", err)
			}

			fmt.Println("✓ API key created successfully.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name for the API key")
	cmd.Flags().StringVarP(&expiresIn, "expires", "e", "", "Expiration duration (e.g., 30d, 1y)")
	cmd.Flags().String("mnemonic", "", "Provide mnemonic (12 words) for authentication")

	return cmd
}
