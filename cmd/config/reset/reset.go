package reset

import (
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
	"github.com/spf13/cobra"
)

func NewResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to defaults",
		Long:  `Reset all configuration values to their default state`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configMgr := config.GetConfigManager()
			if err := configMgr.Reset(); err != nil {
				return output.Error("Failed to reset config: %w", err)
			}

			output.Success("Configuration reset to defaults")
			return nil
		},
	}
}
