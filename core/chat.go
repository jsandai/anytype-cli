package core

import (
	"context"
	"fmt"
	"time"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
)

// ChatMessage represents a chat message for display
type ChatMessage struct {
	ID        string
	OrderID   string
	Creator   string
	Text      string
	CreatedAt time.Time
	Reactions map[string][]string // emoji -> user IDs
	ReplyTo   string
	Read      bool
}

// ChatInfo represents a chat object found in a space
type ChatInfo struct {
	ChatID   string
	Name     string
	ObjectID string
}

// parseChatMessage converts a protobuf ChatMessage to our ChatMessage type
func parseChatMessage(m *model.ChatMessage) ChatMessage {
	msg := ChatMessage{
		ID:        m.Id,
		OrderID:   m.OrderId,
		Creator:   m.Creator,
		CreatedAt: time.Unix(m.CreatedAt, 0),
		ReplyTo:   m.ReplyToMessageId,
		Read:      m.Read,
	}
	if m.Message != nil {
		msg.Text = m.Message.Text
	}
	if m.Reactions != nil {
		msg.Reactions = make(map[string][]string)
		for emoji, identList := range m.Reactions.Reactions {
			msg.Reactions[emoji] = identList.Ids
		}
	}
	return msg
}

// FindChats searches for chat objects in a space
func FindChats(spaceId string) ([]ChatInfo, error) {
	var chats []ChatInfo
	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcObjectSearchRequest{
			SpaceId: spaceId,
			Keys:    []string{"id", "name", "chatId", "hasChat", "layout"},
			Limit:   100,
		}
		resp, err := client.ObjectSearch(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to search objects: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcObjectSearchResponseError_NULL {
			return fmt.Errorf("search error: %s", resp.Error.Description)
		}

		for _, record := range resp.Records {
			chatId := pbtypes.GetString(record, "chatId")
			hasChat := pbtypes.GetBool(record, "hasChat")

			if chatId != "" || hasChat {
				chats = append(chats, ChatInfo{
					ChatID:   chatId,
					Name:     pbtypes.GetString(record, "name"),
					ObjectID: pbtypes.GetString(record, "id"),
				})
			}
		}
		return nil
	})
	return chats, err
}

// SendChatMessage sends a message to a chat object
func SendChatMessage(chatObjectId string, text string, replyToMsgId string) (string, error) {
	var msgId string
	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		msg := &model.ChatMessage{
			Message: &model.ChatMessageMessageContent{
				Text: text,
			},
		}
		if replyToMsgId != "" {
			msg.ReplyToMessageId = replyToMsgId
		}

		req := &pb.RpcChatAddMessageRequest{
			ChatObjectId: chatObjectId,
			Message:      msg,
		}
		resp, err := client.ChatAddMessage(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcChatAddMessageResponseError_NULL {
			return fmt.Errorf("send message error: %s", resp.Error.Description)
		}
		msgId = resp.MessageId
		return nil
	})
	return msgId, err
}

// GetChatMessages retrieves messages from a chat object
func GetChatMessages(chatObjectId string, limit int32, beforeOrderId string, afterOrderId string) ([]ChatMessage, error) {
	var messages []ChatMessage
	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcChatGetMessagesRequest{
			ChatObjectId:  chatObjectId,
			Limit:         limit,
			BeforeOrderId: beforeOrderId,
			AfterOrderId:  afterOrderId,
		}
		resp, err := client.ChatGetMessages(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get messages: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcChatGetMessagesResponseError_NULL {
			return fmt.Errorf("get messages error: %s", resp.Error.Description)
		}

		for _, m := range resp.Messages {
			messages = append(messages, parseChatMessage(m))
		}
		return nil
	})
	return messages, err
}

// EditChatMessage edits an existing message
func EditChatMessage(chatObjectId string, messageId string, newText string) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcChatEditMessageContentRequest{
			ChatObjectId: chatObjectId,
			MessageId:    messageId,
			EditedMessage: &model.ChatMessage{
				Message: &model.ChatMessageMessageContent{
					Text: newText,
				},
			},
		}
		resp, err := client.ChatEditMessageContent(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to edit message: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcChatEditMessageContentResponseError_NULL {
			return fmt.Errorf("edit message error: %s", resp.Error.Description)
		}
		return nil
	})
}

// DeleteChatMessage deletes a message from a chat
func DeleteChatMessage(chatObjectId string, messageId string) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcChatDeleteMessageRequest{
			ChatObjectId: chatObjectId,
			MessageId:    messageId,
		}
		resp, err := client.ChatDeleteMessage(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to delete message: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcChatDeleteMessageResponseError_NULL {
			return fmt.Errorf("delete message error: %s", resp.Error.Description)
		}
		return nil
	})
}

// ToggleChatReaction adds or removes a reaction from a message
func ToggleChatReaction(chatObjectId string, messageId string, emoji string) (bool, error) {
	var added bool
	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcChatToggleMessageReactionRequest{
			ChatObjectId: chatObjectId,
			MessageId:    messageId,
			Emoji:        emoji,
		}
		resp, err := client.ChatToggleMessageReaction(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to toggle reaction: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcChatToggleMessageReactionResponseError_NULL {
			return fmt.Errorf("toggle reaction error: %s", resp.Error.Description)
		}
		added = resp.Added
		return nil
	})
	return added, err
}

// MarkChatMessagesRead marks messages as read up to a certain point
func MarkChatMessagesRead(chatObjectId string, afterOrderId string, beforeOrderId string) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcChatReadMessagesRequest{
			ChatObjectId:  chatObjectId,
			AfterOrderId:  afterOrderId,
			BeforeOrderId: beforeOrderId,
		}
		resp, err := client.ChatReadMessages(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to mark messages read: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcChatReadMessagesResponseError_NULL {
			return fmt.Errorf("mark read error: %s", resp.Error.Description)
		}
		return nil
	})
}

// ChatEventType represents the type of chat event
type ChatEventType string

const (
	ChatEventAdd             ChatEventType = "add"
	ChatEventUpdate          ChatEventType = "update"
	ChatEventDelete          ChatEventType = "delete"
	ChatEventReaction        ChatEventType = "reaction"
	ChatEventReadStatus      ChatEventType = "read_status"
)

// ChatEvent represents a real-time chat event
type ChatEvent struct {
	Type      ChatEventType `json:"type"`
	MessageID string        `json:"message_id"`
	Message   *ChatMessage  `json:"message,omitempty"`
	IsRead    bool          `json:"is_read,omitempty"`
}

// ChatSubscription holds the state for an active chat subscription
type ChatSubscription struct {
	ChatObjectID string
	SubID        string
	Messages     []ChatMessage
}

// SubscribeToChatMessages subscribes to a chat and returns initial messages.
// The subscription is registered with the server using the provided subId.
// Use the returned subId with ListenForEvents to receive real-time updates,
// and call UnsubscribeFromChat when done.
func SubscribeToChatMessages(chatObjectId string, subId string, limit int32) (*ChatSubscription, error) {
	sub := &ChatSubscription{
		ChatObjectID: chatObjectId,
		SubID:        subId,
	}

	err := GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcChatSubscribeLastMessagesRequest{
			ChatObjectId: chatObjectId,
			Limit:        limit,
			SubId:        subId,
		}
		resp, err := client.ChatSubscribeLastMessages(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to subscribe to chat: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcChatSubscribeLastMessagesResponseError_NULL {
			return fmt.Errorf("subscribe error: %s", resp.Error.Description)
		}

		for _, m := range resp.Messages {
			sub.Messages = append(sub.Messages, parseChatMessage(m))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return sub, nil
}

// UnsubscribeFromChat removes a chat subscription
func UnsubscribeFromChat(subId string) error {
	return GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		req := &pb.RpcChatUnsubscribeRequest{
			SubId: subId,
		}
		resp, err := client.ChatUnsubscribe(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to unsubscribe: %w", err)
		}
		if resp.Error != nil && resp.Error.Code != pb.RpcChatUnsubscribeResponseError_NULL {
			return fmt.Errorf("unsubscribe error: %s", resp.Error.Description)
		}
		return nil
	})
}

// ParseChatEvent extracts a ChatEvent from an EventMessage if it's a chat event
// matching the given subscription ID. Returns nil if not a matching chat event.
func ParseChatEvent(msg *pb.EventMessage, subId string) *ChatEvent {
	// Check for Add event
	if add := msg.GetChatAdd(); add != nil {
		if !containsSubId(add.SubIds, subId) {
			return nil
		}
		var chatMsg *ChatMessage
		if add.Message != nil {
			parsed := parseChatMessage(add.Message)
			chatMsg = &parsed
		}
		return &ChatEvent{
			Type:      ChatEventAdd,
			MessageID: add.Id,
			Message:   chatMsg,
		}
	}

	// Check for Update event
	if upd := msg.GetChatUpdate(); upd != nil {
		if !containsSubId(upd.SubIds, subId) {
			return nil
		}
		var chatMsg *ChatMessage
		if upd.Message != nil {
			parsed := parseChatMessage(upd.Message)
			chatMsg = &parsed
		}
		return &ChatEvent{
			Type:      ChatEventUpdate,
			MessageID: upd.Id,
			Message:   chatMsg,
		}
	}

	// Check for Delete event
	if del := msg.GetChatDelete(); del != nil {
		if !containsSubId(del.SubIds, subId) {
			return nil
		}
		return &ChatEvent{
			Type:      ChatEventDelete,
			MessageID: del.Id,
		}
	}

	// Check for Reaction update event
	if react := msg.GetChatUpdateReactions(); react != nil {
		if !containsSubId(react.SubIds, subId) {
			return nil
		}
		// Build a partial message with just reactions
		var chatMsg *ChatMessage
		if react.Reactions != nil {
			chatMsg = &ChatMessage{
				ID:        react.Id,
				Reactions: make(map[string][]string),
			}
			for emoji, identList := range react.Reactions.Reactions {
				chatMsg.Reactions[emoji] = identList.Ids
			}
		}
		return &ChatEvent{
			Type:      ChatEventReaction,
			MessageID: react.Id,
			Message:   chatMsg,
		}
	}

	// Check for Read status update
	if read := msg.GetChatUpdateMessageReadStatus(); read != nil {
		if !containsSubId(read.SubIds, subId) {
			return nil
		}
		// Read status applies to multiple messages; return first ID
		msgId := ""
		if len(read.Ids) > 0 {
			msgId = read.Ids[0]
		}
		return &ChatEvent{
			Type:      ChatEventReadStatus,
			MessageID: msgId,
			IsRead:    read.IsRead,
		}
	}

	return nil
}

// containsSubId checks if the subscription ID is in the list
func containsSubId(subIds []string, target string) bool {
	for _, id := range subIds {
		if id == target {
			return true
		}
	}
	return false
}
