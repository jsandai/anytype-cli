package subscribe

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewSubscribeCmd() *cobra.Command {
	var (
		jsonOutput   bool
		initialLimit int32
	)

	cmd := &cobra.Command{
		Use:   "subscribe <chat-id>",
		Short: "Subscribe to real-time chat updates",
		Long: `Subscribe to a chat and stream events as they occur.

Events are streamed to stdout until interrupted (Ctrl+C).
Use --json for machine-readable JSONL output suitable for piping to other tools.

Event types:
  add         - New message added
  update      - Message content updated
  delete      - Message deleted
  reaction    - Reactions changed
  read_status - Read status changed`,
		Example: `  # Stream events in human-readable format
  anytype chat subscribe <chat-id>

  # Stream as JSONL for automation
  anytype chat subscribe <chat-id> --json

  # Pipe to a processing script
  anytype chat subscribe <chat-id> --json | while read event; do
    echo "$event" | jq '.type'
  done`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			chatId := args[0]

			// Generate unique subscription ID
			subId := fmt.Sprintf("cli-sub-%s", uuid.New().String()[:8])

			// Get session token for event streaming
			token, found, err := core.GetStoredSessionToken()
			if err != nil {
				return output.Error("Failed to get session token: %w", err)
			}
			if !found || token == "" {
				return output.Error("Not authenticated. Run 'anytype auth login' first")
			}

			// Subscribe and get initial messages
			sub, err := core.SubscribeToChatMessages(chatId, subId, initialLimit)
			if err != nil {
				return output.Error("Failed to subscribe: %w", err)
			}

			// Ensure cleanup on exit
			defer func() {
				if err := core.UnsubscribeFromChat(subId); err != nil {
					output.Warning("Failed to unsubscribe: %v", err)
				}
			}()

			// Start event listener
			eventReceiver, err := core.ListenForEvents(token)
			if err != nil {
				return output.Error("Failed to start event stream: %w", err)
			}

			// Setup signal handling for graceful shutdown
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				<-sigChan
				if !jsonOutput {
					output.Info("\nUnsubscribing...")
				}
				cancel()
			}()

			// Print initial messages if any
			if len(sub.Messages) > 0 && !jsonOutput {
				output.Info("=== Recent messages ===")
				for _, msg := range sub.Messages {
					printMessage(msg)
				}
				output.Info("=== Listening for new events (Ctrl+C to stop) ===")
				output.Info("")
			} else if !jsonOutput {
				output.Info("Listening for events on chat %s (Ctrl+C to stop)...", chatId)
				output.Info("")
			}

			// Main event loop
			for {
				select {
				case <-ctx.Done():
					return nil
				default:
					// Wait for next event with timeout to allow checking context
					waitCtx, waitCancel := context.WithTimeout(ctx, 1*time.Second)
					msg, err := eventReceiver.WaitOne(waitCtx)
					waitCancel()

					if err != nil {
						if ctx.Err() != nil {
							// Context cancelled, clean exit
							return nil
						}
						// Timeout, continue loop
						continue
					}

					// Parse and filter chat events for our subscription
					event := core.ParseChatEvent(msg, subId)
					if event == nil {
						continue
					}

					if jsonOutput {
						printEventJSON(event)
					} else {
						printEventFormatted(event)
					}
				}
			}
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output events as JSONL (one JSON object per line)")
	cmd.Flags().Int32VarP(&initialLimit, "initial", "n", 10, "Number of recent messages to fetch initially")

	return cmd
}

// printMessage prints a chat message in human-readable format
func printMessage(msg core.ChatMessage) {
	timestamp := msg.CreatedAt.Format("15:04:05")
	output.Info("[%s] %s: %s", timestamp, msg.Creator, msg.Text)
}

// printEventJSON outputs an event as JSON
func printEventJSON(event *core.ChatEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		output.Warning("Failed to marshal event: %v", err)
		return
	}
	fmt.Println(string(data))
}

// printEventFormatted outputs an event in human-readable format
func printEventFormatted(event *core.ChatEvent) {
	timestamp := time.Now().Format("15:04:05")

	switch event.Type {
	case core.ChatEventAdd:
		if event.Message != nil {
			output.Info("[%s] + %s: %s", timestamp, event.Message.Creator, event.Message.Text)
		} else {
			output.Info("[%s] + New message: %s", timestamp, event.MessageID)
		}

	case core.ChatEventUpdate:
		if event.Message != nil {
			output.Info("[%s] ~ %s edited: %s", timestamp, event.Message.Creator, event.Message.Text)
		} else {
			output.Info("[%s] ~ Message edited: %s", timestamp, event.MessageID)
		}

	case core.ChatEventDelete:
		output.Info("[%s] - Message deleted: %s", timestamp, event.MessageID)

	case core.ChatEventReaction:
		if event.Message != nil && len(event.Message.Reactions) > 0 {
			reactionStr := ""
			for emoji, users := range event.Message.Reactions {
				reactionStr += fmt.Sprintf(" %s(%d)", emoji, len(users))
			}
			output.Info("[%s] ♥ Reactions updated:%s", timestamp, reactionStr)
		} else {
			output.Info("[%s] ♥ Reactions updated: %s", timestamp, event.MessageID)
		}

	case core.ChatEventReadStatus:
		status := "unread"
		if event.IsRead {
			status = "read"
		}
		output.Info("[%s] ● Messages marked %s", timestamp, status)
	}
}
