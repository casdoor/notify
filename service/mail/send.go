package mail

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	mail "github.com/xhit/go-simple-mail/v2"

	"github.com/nikoksr/notify/v2"
)

func (s *Service) buildEmailPayload(conf *SendConfig) (*mail.Email, error) {
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
		if _, err := buf.ReadFrom(attachment.Reader()); err != nil {
			return nil, fmt.Errorf("read attachment: %w", err)
		}

		content := base64.StdEncoding.EncodeToString(buf.Bytes())

		email.Attach(&mail.File{
			Name:     attachment.Name(),
			B64Data:  content,
			Inline:   attachment.Inline(),
			MimeType: attachment.ContentType(),
		})

		s.logger.Info().Str("attachment", attachment.Name()).Msg("Attachment added")
	}

	return email, nil
}

// The function 'send' is responsible for the process of sending a message to every recipient in the list.
//
// It checks if context was cancelled. If yes, it immediately returns the error from context. If not, it tries to send
// the message to the recipients.
//
// Compared to other implementations of the send function, this one does not have an error handling routine. This is
// because the mail package supports sending to multiple recipients at once. We're not looping over the recipients
// ourselves, but instead let the mail package do it for us. If an error occurs, it will be returned immediately.
func (s *Service) send(ctx context.Context, conf *SendConfig) error {
	s.logger.Debug().Msg("Sending message to all recipients")

	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Build email payload
	email, err := s.buildEmailPayload(conf)
	if err != nil {
		return fmt.Errorf("build email payload: %w", err)
	}

	// Quit early if dry run is enabled
	if conf.DryRun {
		s.logger.Info().Strs("recipients", s.recipients).Msg("Dry run enabled - Message not sent.")
		return nil
	}

	// Send the email to the SMTP server.
	if err := email.Send(s.client); err != nil {
		return &notify.SendError{
			FailedRecipients: []string{"unknown"}, // TODO: Not sure how to get the recipient from the error
			Errors:           []error{err},        // TODO: convert to Notify error
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

	// Send message to all recipients
	return s.send(ctx, conf)
}
