package discord

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/nikoksr/notify/v2"
)

func (c *authClient) sendTo(recipient string, conf *SendConfig) error {
	c.logger.Debug().Str("recipient", recipient).Msg("Sending message and attachments to channel")

	// Convert notify.Attachment to discordgo.File.
	files := attachmentsToFiles(conf.Attachments)

	// Send message and attachments.
	_, err := c.session.ChannelMessageSendComplex(recipient, &discordgo.MessageSend{
		Content: conf.Message,
		Files:   files,
	})
	if err != nil {
		return err
	}

	c.logger.Info().Str("recipient", recipient).Msg("Message and attachments sent to channel")

	return nil
}

func (c *webhookClient) sendTo(recipient string, conf *SendConfig) error {
	c.logger.Debug().Str("recipient", recipient).Msg("Sending message and attachments to webhook")

	// Parse the recipient string as a webhook URL.
	// The format is: https://discord.com/api/webhooks/<webhook_id>/<webhook_token>
	u, err := url.Parse(recipient)
	if err != nil {
		return fmt.Errorf("invalid webhook URL: %w", err)
	}

	// Validate the URL.
	if u.Scheme != "https" || u.Host != "discord.com" || !strings.HasPrefix(u.Path, "/api/webhooks/") {
		return fmt.Errorf("invalid webhook URL: %s", u.String())
	}

	// Sanity check to avoid panics.
	segments := strings.Split(u.Path, "/")
	if len(segments) < 3 {
		return fmt.Errorf("invalid webhook URL: %s", u.String())
	}

	// Get the webhook ID and token from the URL.
	// The webhook ID is the second to last path segment.
	// The webhook token is the last path segment.
	webhookID := segments[len(segments)-2]
	webhookToken := segments[len(segments)-1]

	// Convert notify.Attachment to discordgo.File.
	files := attachmentsToFiles(conf.Attachments)

	_, err = c.session.WebhookExecute(webhookID, webhookToken, false, &discordgo.WebhookParams{
		Content: conf.Message,
		Files:   files,
	})
	if err != nil {
		return err
	}

	c.logger.Info().Str("recipient", recipient).Msg("Message and attachments sent to webhook")

	return nil
}

// newSendConfig creates a new send config with default values.
func (s *Service) newSendConfig(subject, message string, opts ...notify.SendOption) *SendConfig {
	conf := &SendConfig{
		Subject:       subject,
		Message:       message,
		DryRun:        s.dryRun,
		ContinueOnErr: s.continueOnErr,
	}

	for _, opt := range opts {
		opt(conf)
	}

	conf.Message = s.renderMessage(conf)

	return conf
}

// The function 'send' is responsible for the process of sending a message to every recipient in the list.
//
// For each recipient, it checks if context was cancelled. If yes, it immediately returns the error from context. If
// not, it tries to send the message to the phone number.
//
// If the message sending process fails, it switches to the error handling routine 'handleError' that appends recipient
// and error into respective slices and logs the error. If the 'ContinueOnErr' option is set to false, the function
// returns the collected errors. If not, it continues to the next recipient.
func (s *Service) send(ctx context.Context, conf *SendConfig) error {
	s.logger.Debug().Msg("Sending message to all recipients")

	var failedRecipients []string
	var errorList []error

	handleError := func(recipient string, err error) {
		// Append error info and log
		failedRecipients = append(failedRecipients, recipient)
		errorList = append(errorList, asNotifyError(err))
		s.logger.Warn().Err(err).Str("recipient", recipient).Msg("Error sending message to recipient")
	}

	for _, recipient := range s.recipients {
		// If context is cancelled, return error immediately
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := s.client.sendTo(recipient, conf); err != nil {
			handleError(recipient, err) // Handle the error

			if !conf.ContinueOnErr {
				// Return collected errors
				return &notify.SendError{
					FailedRecipients: failedRecipients,
					Errors:           errorList,
				}
			}
		}
	}

	// If any errors occurred, return them
	if len(errorList) > 0 {
		return &notify.SendError{
			FailedRecipients: failedRecipients,
			Errors:           errorList,
		}
	}

	s.logger.Info().Msg("Message successfully sent to all recipients")

	return nil
}

// Send takes a message subject and a message body and sends them to all previously set chats.
func (s *Service) Send(ctx context.Context, subject, message string, opts ...notify.SendOption) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.recipients) == 0 {
		return notify.ErrNoRecipients
	}

	// Create new send config from service's default values and passed options
	conf := s.newSendConfig(subject, message, opts...)

	if conf.Message == "" && len(conf.Attachments) == 0 {
		s.logger.Warn().Msg("Message is empty and no attachments are present. Aborting send.")
		return nil
	}

	if conf.DryRun {
		s.logger.Info().Str("message", conf.Message).Msg("Dry run enabled - Message not sent.")
		return nil
	}

	// Send message to all recipients
	return s.send(ctx, conf)
}
