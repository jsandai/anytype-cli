package serve

import (
	"fmt"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core/serviceprogram"
)

func NewServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the Anytype server",
		Long:  "Run the Anytype server. When running as a system service, it will be managed by the service manager. When running interactively, use Ctrl+C to stop.",
		RunE:  runServer,
	}

	return cmd
}

func runServer(cmd *cobra.Command, args []string) error {
	svcConfig := &service.Config{
		Name:        "anytype",
		DisplayName: "Anytype Server",
		Description: "Anytype personal knowledge management server",
	}

	prg := serviceprogram.New()

	s, err := service.New(prg, svcConfig)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	// Run will either:
	// 1. Run as a service if started by service manager
	// 2. Run in foreground if started interactively
	// It handles signals appropriately in both cases
	err = s.Run()
	if err != nil {
		return fmt.Errorf("service failed: %w", err)
	}

	return nil
}
