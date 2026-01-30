package chat

import (
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/cmd/chat/delete"
	"github.com/anyproto/anytype-cli/cmd/chat/edit"
	"github.com/anyproto/anytype-cli/cmd/chat/find"
	"github.com/anyproto/anytype-cli/cmd/chat/list"
	"github.com/anyproto/anytype-cli/cmd/chat/react"
	"github.com/anyproto/anytype-cli/cmd/chat/read"
	"github.com/anyproto/anytype-cli/cmd/chat/send"
	"github.com/anyproto/anytype-cli/cmd/chat/subscribe"
)

func NewChatCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Chat operations",
		Long:  "Send, receive, and manage chat messages in Anytype spaces",
	}

	cmd.AddCommand(find.NewFindCmd())
	cmd.AddCommand(send.NewSendCmd())
	cmd.AddCommand(list.NewListCmd())
	cmd.AddCommand(edit.NewEditCmd())
	cmd.AddCommand(delete.NewDeleteCmd())
	cmd.AddCommand(react.NewReactCmd())
	cmd.AddCommand(read.NewReadCmd())
	cmd.AddCommand(subscribe.NewSubscribeCmd())

	return cmd
}
