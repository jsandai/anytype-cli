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
				if cfg.AccountId != "" {
					output.Info("accountId: %s", cfg.AccountId)
				}
				if cfg.TechSpaceId != "" {
					output.Info("techSpaceId: %s", cfg.TechSpaceId)
				}
				return nil
			}

			key := args[0]
			switch key {
			case "accountId":
				if cfg.AccountId != "" {
					output.Info(cfg.AccountId)
				}
			case "techSpaceId":
				if cfg.TechSpaceId != "" {
					output.Info(cfg.TechSpaceId)
				}
			default:
				return output.Error("unknown config key: %s", key)
			}

			return nil
		},
	}
}
