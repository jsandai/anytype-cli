package get

import (
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Get a configuration value",
		Long:  `Get a specific configuration value or all values if no key is specified`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configMgr := config.GetConfigManager()
			if err := configMgr.Load(); err != nil {
				return output.Error("failed to load config: %w", err)
			}

			cfg := configMgr.Get()

			if len(args) == 0 {
				if cfg.AccountID != "" {
					output.Info("accountId: %s", cfg.AccountID)
				}
				if cfg.TechSpaceID != "" {
					output.Info("techSpaceId: %s", cfg.TechSpaceID)
				}
				return nil
			}

			key := args[0]
			switch key {
			case "accountId", "accountID":
				if cfg.AccountID != "" {
					output.Info(cfg.AccountID)
				}
			case "techSpaceId", "techSpaceID":
				if cfg.TechSpaceID != "" {
					output.Info(cfg.TechSpaceID)
				}
			default:
				return output.Error("unknown config key: %s", key)
			}

			return nil
		},
	}
}
