package create

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/internal"
)

func NewCreateCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new API token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := internal.CreateToken(); err != nil {
				return fmt.Errorf("✗ Failed to create token: %w", err)
			}

			fmt.Println("✓ Token created successfully.")
			return nil
		},
	}

	createCmd.Flags().String("mnemonic", "", "Provide mnemonic (12 words) for authentication")

	return createCmd
}
