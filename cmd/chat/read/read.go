package read

import (
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewReadCmd() *cobra.Command {
	var (
		afterOrderId  string
		beforeOrderId string
	)

	cmd := &cobra.Command{
		Use:   "read <chat-id>",
		Short: "Mark chat messages as read",
		Long:  "Mark messages in a chat as read within an optional range",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			chatId := args[0]

			err := core.MarkChatMessagesRead(chatId, afterOrderId, beforeOrderId)
			if err != nil {
				return output.Error("Failed to mark messages as read: %w", err)
			}

			output.Info("Messages marked as read")
			return nil
		},
	}

	cmd.Flags().StringVar(&afterOrderId, "after", "", "Mark messages after this order ID")
	cmd.Flags().StringVar(&beforeOrderId, "before", "", "Mark messages before this order ID")

	return cmd
}
