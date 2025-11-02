package update

import (
	"github.com/anyproto/anytype-cli/core/output"
	"github.com/anyproto/anytype-cli/core/update"
	"github.com/spf13/cobra"
)

func NewUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update to the latest version",
		Long:  "Download and install the latest version of the Anytype CLI from GitHub releases.",
		RunE: func(cmd *cobra.Command, args []string) error {
			output.Info("Checking for updates...")

			latest, err := update.GetLatestVersion()
			if err != nil {
				return output.Error("Failed to check latest version: %w", err)
			}

			current := update.GetCurrentVersion()

			if !update.NeedsUpdate(current, latest) {
				output.Info("Already up to date (%s)", current)
				return nil
			}

			output.Info("Updating from %s to %s...", current, latest)

			if err := update.DownloadAndInstall(latest); err != nil {
				return output.Error("Update failed: %w", err)
			}

			output.Success("Successfully updated to %s", latest)
			output.Info("If the service is installed, restart it with: anytype service restart")
			output.Info("Otherwise, restart your terminal or run 'anytype' to use the new version")
			return nil
		},
	}
}
