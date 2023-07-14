package telegram

import (
	"context"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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

// sendToChat sends a message to a chat. It returns an error if the message could not be sent.
func (s *Service) sendToChat(ctx context.Context, chatID int64, conf SendConfig) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if conf.message == "" {
		return nil
	}

	if len(conf.attachments) == 0 {
		return s.sendTextMessage(chatID, conf)
	}

	return s.sendFileAttachments(ctx, chatID, conf)
}

// sendTextMessage sends a text message
func (s *Service) sendTextMessage(chatID int64, conf SendConfig) error {
	message := telegram.NewMessage(chatID, conf.message)
	message.ParseMode = conf.parseMode

	_, err := s.client.Send(message)
	return err
}

// sendFileAttachments sends file attachments
func (s *Service) sendFileAttachments(ctx context.Context, chatID int64, conf SendConfig) error {
	for idx, attachment := range conf.attachments {
		isFirst := idx == 0
		if err := s.sendFile(chatID, conf, isFirst, attachment); err != nil {
			return err
		}
	}

	return nil
}

// sendFile sends an individual file
func (s *Service) sendFile(chatID int64, conf SendConfig, isFirst bool, attachment notify.Attachment) error {
	document := telegram.NewDocument(chatID, telegram.FileReader{
		Reader: attachment,
		Name:   attachment.Name(),
	})

	// Set caption only for the first file
	if isFirst {
		document.Caption = conf.message
		document.ParseMode = conf.parseMode
	}

	_, err := s.client.Send(document)
	return err
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
