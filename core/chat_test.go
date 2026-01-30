package core

import (
	"testing"
	"time"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/stretchr/testify/assert"
)

func TestParseChatMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    *model.ChatMessage
		expected ChatMessage
	}{
		{
			name: "basic message",
			input: &model.ChatMessage{
				Id:        "msg-123",
				OrderId:   "order-456",
				Creator:   "user-789",
				CreatedAt: 1704067200, // 2024-01-01 00:00:00 UTC
				Message: &model.ChatMessageMessageContent{
					Text: "Hello, world!",
				},
				Read: true,
			},
			expected: ChatMessage{
				ID:        "msg-123",
				OrderID:   "order-456",
				Creator:   "user-789",
				Text:      "Hello, world!",
				CreatedAt: time.Unix(1704067200, 0),
				Read:      true,
				Reactions: nil,
			},
		},
		{
			name: "message with reply",
			input: &model.ChatMessage{
				Id:               "msg-456",
				OrderId:          "order-789",
				Creator:          "user-abc",
				CreatedAt:        1704153600,
				ReplyToMessageId: "msg-123",
				Message: &model.ChatMessageMessageContent{
					Text: "This is a reply",
				},
			},
			expected: ChatMessage{
				ID:        "msg-456",
				OrderID:   "order-789",
				Creator:   "user-abc",
				Text:      "This is a reply",
				CreatedAt: time.Unix(1704153600, 0),
				ReplyTo:   "msg-123",
				Reactions: nil,
			},
		},
		{
			name: "message with reactions",
			input: &model.ChatMessage{
				Id:        "msg-789",
				OrderId:   "order-abc",
				Creator:   "user-def",
				CreatedAt: 1704240000,
				Message: &model.ChatMessageMessageContent{
					Text: "React to me!",
				},
				Reactions: &model.ChatMessageReactions{
					Reactions: map[string]*model.ChatMessageReactionsIdentityList{
						"👍": {Ids: []string{"user-1", "user-2"}},
						"❤️": {Ids: []string{"user-3"}},
					},
				},
			},
			expected: ChatMessage{
				ID:        "msg-789",
				OrderID:   "order-abc",
				Creator:   "user-def",
				Text:      "React to me!",
				CreatedAt: time.Unix(1704240000, 0),
				Reactions: map[string][]string{
					"👍": {"user-1", "user-2"},
					"❤️": {"user-3"},
				},
			},
		},
		{
			name: "message with nil content",
			input: &model.ChatMessage{
				Id:        "msg-nil",
				OrderId:   "order-nil",
				Creator:   "user-nil",
				CreatedAt: 1704326400,
				Message:   nil,
			},
			expected: ChatMessage{
				ID:        "msg-nil",
				OrderID:   "order-nil",
				Creator:   "user-nil",
				Text:      "",
				CreatedAt: time.Unix(1704326400, 0),
				Reactions: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseChatMessage(tt.input)

			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.OrderID, result.OrderID)
			assert.Equal(t, tt.expected.Creator, result.Creator)
			assert.Equal(t, tt.expected.Text, result.Text)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tt.expected.ReplyTo, result.ReplyTo)
			assert.Equal(t, tt.expected.Read, result.Read)

			if tt.expected.Reactions == nil {
				assert.Nil(t, result.Reactions)
			} else {
				assert.Equal(t, len(tt.expected.Reactions), len(result.Reactions))
				for emoji, users := range tt.expected.Reactions {
					assert.ElementsMatch(t, users, result.Reactions[emoji])
				}
			}
		})
	}
}

func TestContainsSubId(t *testing.T) {
	tests := []struct {
		name     string
		subIds   []string
		target   string
		expected bool
	}{
		{
			name:     "found in list",
			subIds:   []string{"sub-1", "sub-2", "sub-3"},
			target:   "sub-2",
			expected: true,
		},
		{
			name:     "not found",
			subIds:   []string{"sub-1", "sub-2", "sub-3"},
			target:   "sub-4",
			expected: false,
		},
		{
			name:     "empty list",
			subIds:   []string{},
			target:   "sub-1",
			expected: false,
		},
		{
			name:     "nil list",
			subIds:   nil,
			target:   "sub-1",
			expected: false,
		},
		{
			name:     "single element match",
			subIds:   []string{"sub-1"},
			target:   "sub-1",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsSubId(tt.subIds, tt.target)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseChatEvent(t *testing.T) {
	subId := "test-sub"
	otherSubId := "other-sub"

	tests := []struct {
		name     string
		msg      *pb.EventMessage
		subId    string
		expected *ChatEvent
	}{
		{
			name: "chat add event with matching subId",
			msg: &pb.EventMessage{
				Value: &pb.EventMessageValueOfChatAdd{
					ChatAdd: &pb.EventChatAdd{
						Id:     "msg-123",
						SubIds: []string{subId},
						Message: &model.ChatMessage{
							Id:      "msg-123",
							Creator: "user-1",
							Message: &model.ChatMessageMessageContent{
								Text: "Hello",
							},
						},
					},
				},
			},
			subId: subId,
			expected: &ChatEvent{
				Type:      ChatEventAdd,
				MessageID: "msg-123",
				Message: &ChatMessage{
					ID:      "msg-123",
					Creator: "user-1",
					Text:    "Hello",
				},
			},
		},
		{
			name: "chat add event with non-matching subId",
			msg: &pb.EventMessage{
				Value: &pb.EventMessageValueOfChatAdd{
					ChatAdd: &pb.EventChatAdd{
						Id:     "msg-123",
						SubIds: []string{otherSubId},
					},
				},
			},
			subId:    subId,
			expected: nil,
		},
		{
			name: "chat delete event",
			msg: &pb.EventMessage{
				Value: &pb.EventMessageValueOfChatDelete{
					ChatDelete: &pb.EventChatDelete{
						Id:     "msg-456",
						SubIds: []string{subId},
					},
				},
			},
			subId: subId,
			expected: &ChatEvent{
				Type:      ChatEventDelete,
				MessageID: "msg-456",
			},
		},
		{
			name: "chat update event",
			msg: &pb.EventMessage{
				Value: &pb.EventMessageValueOfChatUpdate{
					ChatUpdate: &pb.EventChatUpdate{
						Id:     "msg-789",
						SubIds: []string{subId},
						Message: &model.ChatMessage{
							Id:      "msg-789",
							Creator: "user-2",
							Message: &model.ChatMessageMessageContent{
								Text: "Updated text",
							},
						},
					},
				},
			},
			subId: subId,
			expected: &ChatEvent{
				Type:      ChatEventUpdate,
				MessageID: "msg-789",
				Message: &ChatMessage{
					ID:      "msg-789",
					Creator: "user-2",
					Text:    "Updated text",
				},
			},
		},
		{
			name: "non-chat event returns nil",
			msg: &pb.EventMessage{
				Value: &pb.EventMessageValueOfAccountShow{
					AccountShow: &pb.EventAccountShow{},
				},
			},
			subId:    subId,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseChatEvent(tt.msg, tt.subId)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.expected.Type, result.Type)
			assert.Equal(t, tt.expected.MessageID, result.MessageID)

			if tt.expected.Message != nil {
				assert.NotNil(t, result.Message)
				assert.Equal(t, tt.expected.Message.ID, result.Message.ID)
				assert.Equal(t, tt.expected.Message.Creator, result.Message.Creator)
				assert.Equal(t, tt.expected.Message.Text, result.Message.Text)
			}
		})
	}
}
