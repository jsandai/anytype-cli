package core

import (
	"testing"
	"time"

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
