package set

import (
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
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
				return output.Error("failed to load config: %w", err)
			}

			switch key {
			case "accountId":
				if err := configMgr.SetAccountId(value); err != nil {
					return output.Error("failed to set account Id: %w", err)
				}
			case "techSpaceId":
				if err := configMgr.SetTechSpaceId(value); err != nil {
					return output.Error("failed to set tech space Id: %w", err)
				}
			default:
				return output.Error("unknown config key: %s", key)
			}

			output.Success("Set %s = %s", key, value)
			return nil
		},
	}
}
