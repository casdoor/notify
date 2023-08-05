package twilio

import (
	"context"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/nikoksr/notify/v2"
)

func (s *Service) buildMessagePayload(phoneNumber string, conf *SendConfig) *twilioApi.CreateMessageParams {
	params := &twilioApi.CreateMessageParams{}
	params.SetFrom(conf.SenderPhoneNumber)
	params.SetTo(phoneNumber)
	params.SetBody(conf.Message)

	return params
}

// sendToPhoneNumber sends a message to a chat. It returns an error if the message could not be sent.
func (s *Service) sendToPhoneNumber(phoneNumber string, conf *SendConfig) error {
	s.logger.Debug().Str("recipient", phoneNumber).Msg("Sending text message to chat")

	// Build message payload
	message := s.buildMessagePayload(phoneNumber, conf)

	// Quit early if dry run is enabled
	if conf.DryRun {
		s.logger.Info().Str("recipient", phoneNumber).Msg("Dry run enabled - Message not sent.")
		return nil
	}

	// Send the message
	if _, err := s.client.Api.CreateMessage(message); err != nil {
		return err
	}

	s.logger.Info().Str("recipient", phoneNumber).Msg("Text message sent to chat")

	return nil
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

	handleError := func(phoneNumber string, err error) {
		// Append error info and log
		failedRecipients = append(failedRecipients, phoneNumber)
		errorList = append(errorList, asNotifyError(err))
		s.logger.Warn().Err(err).Str("recipient", phoneNumber).Msg("Error sending message to recipient")
	}

	for _, phoneNumber := range conf.Recipients {
		// If context is cancelled, return error immediately
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Send the message
		if err := s.sendToPhoneNumber(phoneNumber, conf); err != nil {
			handleError(phoneNumber, err) // Handle the error

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

// Send sends a message to all chats that are configured to receive messages. It returns an error if the message could
// not be sent.
func (s *Service) Send(ctx context.Context, subject, message string, opts ...notify.SendOption) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create new send config from service's default values and passed options
	conf := s.newSendConfig(subject, message, opts...)

	if len(conf.Recipients) == 0 {
		return notify.ErrNoRecipients
	}

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
