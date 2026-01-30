package send

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

var (
	replyTo    string
	attachIds  []string
	filePaths  []string
	spaceId    string
)

func NewSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send <chat-id> [message]",
		Short: "Send a message to a chat",
		Long: `Send a text message to an Anytype chat object, optionally with attachments.

Examples:
  anytype chat send <chat-id> "Hello world"
  anytype chat send <chat-id> "Check this out" --file /path/to/image.png --space <space-id>
  anytype chat send <chat-id> --attach <file-object-id>
  anytype chat send <chat-id> "Multiple files" --file a.png --file b.jpg --space <space-id>`,
		Args: cobra.RangeArgs(1, 2),
		RunE: runSend,
	}

	cmd.Flags().StringVar(&replyTo, "reply-to", "", "Message ID to reply to")
	cmd.Flags().StringArrayVar(&attachIds, "attach", nil, "Object ID of existing file to attach (can be repeated)")
	cmd.Flags().StringArrayVar(&filePaths, "file", nil, "Local file path to upload and attach (can be repeated)")
	cmd.Flags().StringVar(&spaceId, "space", "", "Space ID (required when using --file)")

	return cmd
}

func runSend(cmd *cobra.Command, args []string) error {
	chatId := args[0]
	message := ""
	if len(args) > 1 {
		message = args[1]
	}

	// Validate: need either message or attachments
	if message == "" && len(attachIds) == 0 && len(filePaths) == 0 {
		return fmt.Errorf("must provide either a message or attachments")
	}

	// Validate: --file requires --space
	if len(filePaths) > 0 && spaceId == "" {
		return fmt.Errorf("--space is required when using --file")
	}

	var attachments []core.ChatAttachment

	// Process --attach flags (existing file object IDs)
	for _, id := range attachIds {
		attachments = append(attachments, core.ChatAttachment{
			ObjectId: id,
			Type:     model.ChatMessageAttachment_FILE, // Default to file, could be improved
		})
	}

	// Process --file flags (upload and attach)
	for _, path := range filePaths {
		// Expand ~ and make absolute
		if strings.HasPrefix(path, "~/") {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, path[2:])
		}
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to resolve path %s: %w", path, err)
		}

		// Check file exists
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", absPath)
		}

		// Detect file type
		fileType := core.DetectFileType(absPath)

		// Upload file
		output.Info("Uploading %s...", filepath.Base(absPath))
		result, err := core.UploadFile(spaceId, absPath, fileType)
		if err != nil {
			return fmt.Errorf("failed to upload %s: %w", path, err)
		}
		output.Info("Uploaded: %s", result.ObjectId)

		// Add as attachment
		attachments = append(attachments, core.ChatAttachment{
			ObjectId: result.ObjectId,
			Type:     core.DetectAttachmentType(fileType),
		})
	}

	// Send message
	var msgId string
	var err error

	if len(attachments) > 0 {
		msgId, err = core.SendChatMessageWithAttachments(chatId, message, replyTo, attachments)
	} else {
		msgId, err = core.SendChatMessage(chatId, message, replyTo)
	}

	if err != nil {
		return output.Error("Failed to send message: %w", err)
	}

	output.Info("Message sent successfully")
	output.Info("Message ID: %s", msgId)
	return nil
}
