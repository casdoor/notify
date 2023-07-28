package mail

import (
	"bytes"
	"context"

	mail "github.com/xhit/go-simple-mail/v2"

	"github.com/nikoksr/notify/v2"
)

// newSendConfig creates a new send config with default values.
func (s *Service) newSendConfig(subject, message string, opts ...notify.SendOption) *SendConfig {
	conf := &SendConfig{
		Subject: subject,
		Message: message,
		DryRun:  s.dryRun,
	}

	for _, opt := range opts {
		opt(conf)
	}

	conf.Message = s.renderMessage(conf)

	return conf
}

// send sends a message to all recipients. It returns an error if the message could not be sent.
func (s *Service) send(ctx context.Context, conf *SendConfig) error {
	s.logger.Debug().Msg("Sending message to all recipients")

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Create a new email message.
	email := mail.NewMSG().
		SetFrom(s.senderName).
		AddTo(s.recipients...).
		AddCc(s.ccRecipients...).
		AddBcc(s.bccRecipients...).
		SetPriority(mail.Priority(s.priority)).
		SetSubject(conf.Subject).
		SetBody(mail.ContentType(conf.ParseMode),
			conf.Message,
		)

	// Add attachments
	for _, attachment := range conf.Attachments {
		s.logger.Debug().Str("attachment", attachment.Name()).Msg("Adding attachment")

		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(attachment); err != nil {
			return err
		}

		email.Attach(&mail.File{
			Name:   attachment.Name(),
			Data:   buf.Bytes(),
			Inline: conf.InlineAttachments,
		})

		s.logger.Info().Str("attachment", attachment.Name()).Msg("Attachment added")
	}

	// Send the email to the SMTP server.
	if err := email.Send(s.client); err != nil {
		return &notify.SendNotificationError{
			Recipient: "unknown", // TODO: Not sure how to get the recipient from the error
			Cause:     err,       // TODO: convert to Notify error
		}
	}

	s.logger.Info().Msg("Message successfully sent to all recipients")

	return nil
}

// Send sends a notification to each Mail channel defined in Service. The sender is configured through SendOption and
// SendConfig. Returns an error upon failure to send the message, or if there are no recipients identified.
func (s *Service) Send(ctx context.Context, subject, message string, opts ...notify.SendOption) error {
	s.mu.Lock()
	defer s.mu.Unlock()

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
