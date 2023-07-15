package discord

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/nikoksr/notify/v2"
)

func (c *authClient) sendTo(recipient string, conf SendConfig) error {
	// Convert notify.Attachment to discordgo.File.
	files := attachmentsToFiles(conf.attachments)

	// Send message and attachments.
	_, err := c.session.ChannelMessageSendComplex(recipient, &discordgo.MessageSend{
		Content: conf.message,
		Files:   files,
	})

	return err
}

func (c *webhookClient) sendTo(recipient string, conf SendConfig) error {
	// Parse the recipient string as a webhook URL.
	// The format is: https://discord.com/api/webhooks/<webhook_id>/<webhook_token>
	u, err := url.Parse(recipient)
	if err != nil {
		return fmt.Errorf("invalid webhook URL: %w", err)
	}

	// Get the webhook ID and token from the URL.
	// The webhook ID is the second to last path segment.
	// The webhook token is the last path segment.
	segments := strings.Split(u.Path, "/")
	webhookID := segments[len(segments)-2]
	webhookToken := segments[len(segments)-1]

	// Convert notify.Attachment to discordgo.File.
	files := attachmentsToFiles(conf.attachments)

	_, err = c.session.WebhookExecute(webhookID, webhookToken, false, &discordgo.WebhookParams{
		Content: conf.message,
		Files:   files,
	})

	return err
}

// sendTo sends a message to a channel or a webhook URL. It returns an error if the message could not be sent.
func (s *Service) sendTo(_ context.Context, recipient string, conf SendConfig) error {
	return s.client.sendTo(recipient, conf)
}

// Send takes a message subject and a message body and sends them to all previously set chats.
func (s *Service) Send(ctx context.Context, subject, message string, opts ...notify.SendOption) error {
	if len(s.recipients) == 0 {
		return notify.ErrNoRecipients
	}

	conf := SendConfig{
		subject: subject,
		message: message,
	}

	for _, opt := range opts {
		opt(&conf)
	}

	conf.message = s.renderMessage(conf)

	for _, recipient := range s.recipients {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := s.sendTo(ctx, recipient, conf); err != nil {
			return &notify.SendNotificationError{
				Recipient: recipient,
				Cause:     asNotifyError(err),
			}
		}
	}

	return nil
}
