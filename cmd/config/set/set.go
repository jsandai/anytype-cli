package set

import (
	"fmt"

	"github.com/anyproto/anytype-cli/core/config"
	"github.com/spf13/cobra"
)

func NewSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long:  `Set a specific configuration value`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			configMgr := config.GetConfigManager()
			if err := configMgr.Load(); err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			switch key {
			case "accountId", "accountID":
				if err := configMgr.SetAccountID(value); err != nil {
					return fmt.Errorf("failed to set account ID: %w", err)
				}
			case "techSpaceId", "techSpaceID":
				if err := configMgr.SetTechSpaceID(value); err != nil {
					return fmt.Errorf("failed to set tech space ID: %w", err)
				}
			default:
				return fmt.Errorf("unknown config key: %s", key)
			}

			fmt.Printf("Set %s = %s\n", key, value)
			return nil
		},
	}
}
