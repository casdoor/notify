package twilio

import (
	"context"
	"errors"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/nikoksr/notify/v2"
)

// sendToPhoneNumber sends a message to a chat. It returns an error if the message could not be sent.
func (s *Service) sendToPhoneNumber(phoneNumber string, conf *SendConfig) error {
	if conf == nil {
		return errors.New("send config is nil")
	}

	s.logger.Debug().Str("recipient", phoneNumber).Msg("Sending text message to chat")

	params := &twilioApi.CreateMessageParams{}
	params.SetFrom(s.senderPhoneNumber)
	params.SetTo(phoneNumber)
	params.SetBody(conf.Message)

	if _, err := s.client.Api.CreateMessage(params); err != nil {
		return err
	}

	s.logger.Info().Str("recipient", phoneNumber).Msg("Text message sent to chat")

	return nil
}

// newSendConfig creates a new send config with default values.
func (s *Service) newSendConfig(subject, message string, opts ...notify.SendOption) *SendConfig {
	conf := &SendConfig{
		Subject: subject,
		Message: message,
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

	for _, phoneNumber := range s.phoneNumbers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := s.sendToPhoneNumber(phoneNumber, conf); err != nil {
			return &notify.SendNotificationError{
				Recipient: phoneNumber,
				Cause:     asNotifyError(err),
			}
		}
	}

	s.logger.Info().Msg("Message successfully sent to all recipients")

	return nil
}

// Send sends a message to all chats that are configured to receive messages. It returns an error if the message could
// not be sent.
func (s *Service) Send(ctx context.Context, subject, message string, opts ...notify.SendOption) error {
	if len(s.phoneNumbers) == 0 {
		return notify.ErrNoRecipients
	}

	// Create new send config from service's default values and passed options
	conf := s.newSendConfig(subject, message, opts...)

	if conf.Message == "" && len(conf.Attachments) == 0 {
		s.logger.Warn().Msg("Message is empty and no attachments are present. Aborting send.")
		return nil
	}

	if len(conf.Attachments) > 0 {
		s.logger.Debug().Msg("Attachments are not supported by Twilio")
	}

	// Send message to all recipients
	return s.send(ctx, conf)
}
