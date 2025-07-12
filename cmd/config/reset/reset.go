package reset

import (
	"fmt"

	"github.com/anyproto/anytype-cli/core/config"
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
				return fmt.Errorf("failed to reset config: %w", err)
			}

			fmt.Println("Configuration reset to defaults")
			return nil
		},
	}
}
