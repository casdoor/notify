package discord

import (
	"context"

	"github.com/nikoksr/notify/v2"
)

var _ notify.SendConfig = (*SendConfig)(nil)

// SendConfig is the configuration for sending a message to a channel or a webhook URL. It implements the
// notify.SendConfig interface.
type SendConfig struct {
	subject     string
	message     string
	attachments []notify.Attachment
	metadata    map[string]any
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

// SetAttachments adds attachments to the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetAttachments(attachments ...notify.Attachment) {
	c.attachments = attachments
}

// SetMetadata sets the metadata of the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetMetadata(metadata map[string]any) {
	c.metadata = metadata
}

// sendTo sends a message to a channel or a webhook URL. It returns an error if the message could not be sent.
func (s *Service) sendTo(ctx context.Context, receiver string, conf SendConfig) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := s.client.sendTo(receiver, conf); err != nil {
			return err
		}
	}

	return nil
}

// Send takes a message subject and a message body and sends them to all previously set chats.
func (s *Service) Send(ctx context.Context, subject, message string, opts ...notify.SendOption) error {
	if len(s.receivers) == 0 {
		return notify.ErrNoReceivers
	}

	conf := SendConfig{
		subject: subject,
		message: message,
	}

	for _, opt := range opts {
		opt(&conf)
	}

	conf.message = s.renderMessage(conf)

	for _, receiver := range s.receivers {
		if err := s.sendTo(ctx, receiver, conf); err != nil {
			return notify.NewErrSendNotification(receiver, err)
		}
	}

	return nil
}
