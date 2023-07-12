package telegram

import (
	"context"
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/nikoksr/notify/v2"
)

var _ notify.SendConfig = (*SendConfig)(nil)

// SendConfig is the configuration for sending a message. It implements the
// notify.SendConfig interface.
type SendConfig struct {
	subject     string
	message     string
	attachments []notify.Attachment
	metadata    map[string]any

	// Custom fields
	parseMode string
}

// Common fields

// Subject returns the subject of the message.
func (c *SendConfig) Subject() string {
	return c.subject
}

// Message returns the message.
func (c *SendConfig) Message() string {
	return c.message
}

// Attachments returns the attachments.
func (c *SendConfig) Attachments() []notify.Attachment {
	return c.attachments
}

// Metadata returns the metadata.
func (c *SendConfig) Metadata() map[string]any {
	return c.metadata
}

// Telegram specific fields

// ParseMode returns the parse mode of the message. This is a Telegram specific option.
func (c *SendConfig) ParseMode() string {
	return c.parseMode
}

// notify.SendConfig implementation

// SetAttachments adds attachments to the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetAttachments(attachments ...notify.Attachment) {
	c.attachments = attachments
}

// SetMetadata sets the metadata of the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetMetadata(metadata map[string]any) {
	c.metadata = metadata
}

// Send options

// SendWithParseMode is a send option that sets the parse mode of the message.
func SendWithParseMode(parseMode string) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.parseMode = parseMode
		}
	}
}

// Send logic

// sendMessageToChat sends a message to a chat. It returns an error if the message could not be sent.
func (s *Service) sendMessageToChat(ctx context.Context, chatID int64, conf SendConfig) error {
	if conf.message == "" {
		return nil
	}

	message := telegram.NewMessage(chatID, conf.message)
	message.ParseMode = conf.parseMode

	_, err := s.client.Send(message)

	return err
}

// sendAttachmentsToChat sends attachments to a chat. It returns an error if the message could not be sent.
func (s *Service) sendAttachmentsToChat(ctx context.Context, chatID int64, conf SendConfig) error {
	for _, attachment := range conf.attachments {
		document := telegram.NewDocumentUpload(chatID, telegram.FileReader{
			Reader: attachment,
			Name:   attachment.Name(),
			Size:   -1,
		})
		if _, err := s.client.Send(document); err != nil {
			return err
		}
	}

	return nil
}

// sendToChat sends a message to a chat. It returns an error if the message could not be sent.
func (s *Service) sendToChat(ctx context.Context, chatID int64, conf SendConfig) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := s.sendMessageToChat(ctx, chatID, conf); err != nil {
			return fmt.Errorf("send message to chat: %w", err)
		}

		if err := s.sendAttachmentsToChat(ctx, chatID, conf); err != nil {
			return fmt.Errorf("send attachments to chat: %w", err)
		}
	}

	return nil
}

// Send sends a message to all chats that are configured to receive messages. It returns an error if the message could
// not be sent.
func (s *Service) Send(ctx context.Context, subject, message string, opts ...notify.SendOption) error {
	if len(s.chatIDs) == 0 {
		return notify.ErrNoRecipients
	}

	conf := SendConfig{
		parseMode: s.parseMode,
		subject:   subject,
		message:   message,
	}

	for _, opt := range opts {
		opt(&conf)
	}

	conf.message = s.renderMessage(conf)

	for _, chatID := range s.chatIDs {
		if err := s.sendToChat(ctx, chatID, conf); err != nil {
			return notify.NewErrSendNotification(chatID, err)
		}
	}

	return nil
}
