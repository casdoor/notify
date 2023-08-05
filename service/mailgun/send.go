package mailgun

import (
	"context"
	"fmt"
	"io"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/nikoksr/notify/v2"
)

func (s *Service) buildEmailPayload(conf *SendConfig) (*mailgun.Message, error) {
	email := s.client.NewMessage(
		conf.SenderAddress,
		conf.Subject,
		conf.Message,
		s.recipients...,
	)

	for _, cc := range s.ccRecipients {
		email.AddCC(cc)
	}

	for _, bcc := range s.bccRecipients {
		email.AddBCC(bcc)
	}

	// Domain: string
	email.AddDomain(s.domain)

	// Set headers
	for k, v := range s.headers {
		email.AddHeader(k, v)
	}

	// Set tags
	for _, tag := range s.tags {
		if err := email.AddTag(tag); err != nil {
			return nil, fmt.Errorf("add tags: %w", err)
		}
	}

	// Set content type
	switch conf.ParseMode {
	case ModeText:
		// Plain text is already set by default
	case ModeHTML:
		email.SetHtml(conf.Message)
	case ModeAMPHTML:
		email.SetAMPHtml(conf.Message)
	}

	// Security
	if conf.EnableNativeSend {
		email.EnableNativeSend()
	}

	email.SetDKIM(conf.SetDKIM)
	email.SetRequireTLS(conf.RequireTLS)
	email.SetSkipVerification(conf.SkipVerification)

	// Tracking
	email.SetTracking(conf.EnableTracking)
	email.SetTrackingClicks(conf.EnableTrackingClicks)
	email.SetTrackingOpens(conf.EnableTrackingOpens)

	// Add attachments
	for _, attachment := range conf.Attachments {
		s.logger.Debug().Str("attachment", attachment.Name()).Msg("Adding attachment")

		if attachment.Inline() {
			email.AddReaderInline(attachment.Name(), io.NopCloser(attachment.Reader()))
		} else {
			email.AddReaderAttachment(attachment.Name(), io.NopCloser(attachment.Reader()))
		}

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

	// Build the email payload
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
	if _, _, err := s.client.Send(ctx, email); err != nil {
		err = asNotifyError(err)

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
