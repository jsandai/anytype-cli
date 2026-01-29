package join

import (
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewJoinCmd() *cobra.Command {
	var (
		networkId     string
		inviteCid     string
		inviteFileKey string
	)

	cmd := &cobra.Command{
		Use:   "join [invite-link]",
		Short: "Join a space",
		Long:  "Join a space using an invite link (any URL with /{cid}#{key} format) or flags (--invite-cid and --invite-key)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var spaceId string

			if networkId == "" {
				networkId = config.AnytypeNetworkAddress
			}

			// If flags are provided directly, use them (self-hosted support)
			if inviteCid != "" && inviteFileKey != "" {
				info, err := core.ViewSpaceInvite(inviteCid, inviteFileKey)
				if err != nil {
					return output.Error("Failed to view invite: %w", err)
				}

				output.Info("Joining space '%s' created by %s...", info.SpaceName, info.CreatorName)
				spaceId = info.SpaceId
			} else if len(args) > 0 {
				// Parse invite link (supports any URL with /{cid}#{key} format)
				input := args[0]
				u, err := url.Parse(input)
				if err != nil {
					return output.Error("invalid invite link: %w", err)
				}

				path := strings.TrimPrefix(u.Path, "/")
				if path == "" {
					return output.Error("invite link missing Cid")
				}
				inviteCid = path

				inviteFileKey = u.Fragment
				if inviteFileKey == "" {
					return output.Error("invite link missing key (should be after #)")
				}

				info, err := core.ViewSpaceInvite(inviteCid, inviteFileKey)
				if err != nil {
					return output.Error("Failed to view invite: %w", err)
				}

				output.Info("Joining space '%s' created by %s...", info.SpaceName, info.CreatorName)
				spaceId = info.SpaceId
			} else {
				return output.Error("provide an invite link or use --invite-cid and --invite-key flags")
			}

			if err := core.JoinSpace(networkId, spaceId, inviteCid, inviteFileKey); err != nil {
				return output.Error("Failed to join space: %w", err)
			}

			output.Success("Successfully sent join request to space '%s'", spaceId)
			return nil
		},
	}

	cmd.Flags().StringVar(&networkId, "network", "", "Network `id` to join")
	cmd.Flags().StringVar(&inviteCid, "invite-cid", "", "Invite `cid` (extracted from invite link if not provided)")
	cmd.Flags().StringVar(&inviteFileKey, "invite-key", "", "Invite file `key` (extracted from invite link if not provided)")

	return cmd
}
