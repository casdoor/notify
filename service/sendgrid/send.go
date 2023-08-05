package sendgrid

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/nikoksr/notify/v2"
)

func (s *Service) buildEmailPayload(conf *SendConfig) (*mail.SGMailV3, error) {
	from := mail.NewEmail(conf.SenderName, conf.SenderAddress)
	content := mail.NewContent(string(conf.ParseMode), conf.Message)

	// Create a new personalization instance to be able to add multiple receiver addresses.
	personalization := mail.NewPersonalization()
	personalization.Subject = conf.Subject

	for _, recipient := range conf.Recipients {
		personalization.AddTos(mail.NewEmail(recipient, recipient))
	}

	for _, cc := range conf.CCRecipients {
		personalization.AddCCs(mail.NewEmail(cc, cc))
	}

	for _, bcc := range conf.BCCRecipients {
		personalization.AddBCCs(mail.NewEmail(bcc, bcc))
	}

	// Configure email
	email := mail.NewV3Mail()
	email.AddPersonalizations(personalization)
	email.AddContent(content)
	email.AddAttachment()
	email.SetFrom(from)

	for k, v := range conf.Headers {
		email.SetHeader(k, v)
	}

	for k, v := range conf.CustomArgs {
		email.SetCustomArg(k, v)
	}

	// Identification
	email.SetBatchID(conf.BatchID)
	email.SetIPPoolID(conf.IPPoolID)

	// Advanced settings
	email.SetASM(conf.ASM)
	email.SetTrackingSettings(conf.TrackingSettings)
	email.SetMailSettings(conf.MailSettings)

	// Add attachments
	for _, attachment := range conf.Attachments {
		s.logger.Debug().Str("attachment", attachment.Name()).Msg("Adding attachment")

		// Read attachment into string
		var buf bytes.Buffer
		_, err := buf.ReadFrom(attachment.Reader())
		if err != nil {
			return nil, fmt.Errorf("read attachment: %w", err)
		}

		// Add attachment
		disposition := "attachment"
		if attachment.Inline() {
			disposition = "inline"
		}

		email.AddAttachment(&mail.Attachment{
			Name:        attachment.Name(),
			Filename:    attachment.Name(),
			Content:     base64.StdEncoding.EncodeToString(buf.Bytes()),
			Type:        attachment.ContentType(),
			Disposition: disposition,
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

	// Create a new email message.
	email, err := s.buildEmailPayload(conf)
	if err != nil {
		return fmt.Errorf("new message: %w", err)
	}

	// Quit early if dry run is enabled
	if conf.DryRun {
		s.logger.Info().Strs("recipients", conf.Recipients).Msg("Dry run enabled - Message not sent.")
		return nil
	}

	// Send the email
	// TODO: This section needs a check-up. At the time of writing this the sendgrid-go package does not return the
	//       http response in case of an error. We're kinda in the dark here, since they don't provide any defined
	//       errors, nor do they return the http response.
	resp, err := s.client.SendWithContext(ctx, email)
	if err != nil {
		err = asNotifyError(resp, err)
		return &notify.SendError{
			FailedRecipients: []string{"unknown"}, // TODO: Not sure how to get the recipient from the error
			Errors:           []error{err},        // TODO: convert to Notify error
		}
	}

	if err := asNotifyError(resp, nil); err != nil {
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

	// Create new send config from service's default values and passed options
	conf := s.newSendConfig(subject, message, opts...)

	if len(conf.Recipients) == 0 {
		return notify.ErrNoRecipients
	}

	if conf.Message == "" && len(conf.Attachments) == 0 {
		s.logger.Warn().Msg("Message is empty and no attachments are present. Aborting send.")
		return nil
	}

	// Send message to all recipients
	return s.send(ctx, conf)
}
