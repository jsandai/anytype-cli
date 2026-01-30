package find

import (
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewFindCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find <space-id>",
		Short: "Find chat objects in a space",
		Long:  "Search for objects with chat functionality in a space and display their chat IDs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			spaceId := args[0]

			chats, err := core.FindChats(spaceId)
			if err != nil {
				return output.Error("Failed to find chats: %w", err)
			}

			if len(chats) == 0 {
				output.Info("No chat objects found in this space")
				return nil
			}

			output.Info("%-40s %-20s %s", "CHAT ID", "NAME", "OBJECT ID")
			output.Info("%-40s %-20s %s", "───────", "────", "─────────")

			for _, chat := range chats {
				name := chat.Name
				if len(name) > 18 {
					name = name[:15] + "..."
				}
				chatId := chat.ChatID
				if chatId == "" {
					chatId = "(no chatId set)"
				}
				output.Info("%-40s %-20s %s", chatId, name, chat.ObjectID)
			}

			return nil
		},
	}

	return cmd
}
