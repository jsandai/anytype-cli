package service

import (
	"fmt"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core/serviceprogram"
)

// getService creates a service instance with our standard configuration
func getService() (service.Service, error) {
	svcConfig := &service.Config{
		Name:        "anytype",
		DisplayName: "Anytype Server",
		Description: "Anytype personal knowledge management server",
		Arguments:   []string{"serve"},
		Option: service.KeyValue{
			"UserService": true,
		},
	}

	prg := serviceprogram.New()
	return service.New(prg, svcConfig)
}

func NewServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage Anytype as a system service",
		Long:  "Install, uninstall, and control Anytype as a system service.",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "install",
			Short: "Install Anytype as a system service",
			RunE:  installService,
		},
		&cobra.Command{
			Use:   "uninstall",
			Short: "Uninstall the Anytype system service",
			RunE:  uninstallService,
		},
		&cobra.Command{
			Use:   "start",
			Short: "Start the Anytype service",
			RunE:  startService,
		},
		&cobra.Command{
			Use:   "stop",
			Short: "Stop the Anytype service",
			RunE:  stopService,
		},
		&cobra.Command{
			Use:   "restart",
			Short: "Restart the Anytype service",
			RunE:  restartService,
		},
		&cobra.Command{
			Use:   "status",
			Short: "Check the status of the Anytype service",
			RunE:  statusService,
		},
	)

	return cmd
}

func installService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	err = s.Install()
	if err != nil {
		return fmt.Errorf("failed to install service: %w", err)
	}

	fmt.Println("✓ Anytype service installed successfully")
	fmt.Println("\nTo manage the service:")
	fmt.Println("  Start:   anytype service start")
	fmt.Println("  Stop:    anytype service stop")
	fmt.Println("  Restart: anytype service restart")
	fmt.Println("  Status:  anytype service status")

	return nil
}

func uninstallService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	err = s.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall service: %w", err)
	}

	fmt.Println("✓ Anytype service uninstalled successfully")
	return nil
}

func startService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	err = s.Start()
	if err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	fmt.Println("✓ Anytype service started")
	return nil
}

func stopService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	err = s.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}

	fmt.Println("✓ Anytype service stopped")
	return nil
}

func restartService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	err = s.Restart()
	if err != nil {
		return fmt.Errorf("failed to restart service: %w", err)
	}

	fmt.Println("✓ Anytype service restarted")
	return nil
}

func statusService(cmd *cobra.Command, args []string) error {
	s, err := getService()
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	status, err := s.Status()
	if err != nil {
		if err == service.ErrNotInstalled {
			fmt.Println("✗ Anytype service is not installed")
			fmt.Println("  Run 'anytype service install' to install it")
			return nil
		}
		return fmt.Errorf("failed to get service status: %w", err)
	}

	switch status {
	case service.StatusRunning:
		fmt.Println("✓ Anytype service is running")
	case service.StatusStopped:
		fmt.Println("✗ Anytype service is stopped")
		fmt.Println("  Run 'anytype service start' to start it")
	default:
		fmt.Printf("? Anytype service status: %v\n", status)
	}

	return nil
}
