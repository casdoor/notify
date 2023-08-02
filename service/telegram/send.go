package telegram

import (
	"context"
	"strconv"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/nikoksr/notify/v2"
)

// sendToChat sends a message to a chat. It returns an error if the message could not be sent.
func (s *Service) sendToChat(chatID int64, conf *SendConfig) error {
	if len(conf.Attachments) == 0 {
		return s.sendTextMessage(chatID, conf)
	}

	return s.sendFileAttachments(chatID, conf)
}

// sendTextMessage sends a text message
func (s *Service) sendTextMessage(chatID int64, conf *SendConfig) error {
	s.logger.Debug().Int64("recipient", chatID).Msg("Sending text message to chat")

	message := telegram.NewMessage(chatID, conf.Message)
	message.ParseMode = conf.ParseMode

	// Quit early if dry run is enabled
	if conf.DryRun {
		s.logger.Info().Str("recipient", strconv.FormatInt(chatID, 10)).Msg("Dry run enabled - Message not sent.")
		return nil
	}

	// Send the message
	if _, err := s.client.Send(message); err != nil {
		return err
	}

	s.logger.Info().Int64("recipient", chatID).Msg("Text message sent to chat")

	return nil
}

// sendFileAttachments sends file attachments
func (s *Service) sendFileAttachments(chatID int64, conf *SendConfig) error {
	for idx, attachment := range conf.Attachments {
		isFirst := idx == 0
		if err := s.sendFile(chatID, conf, isFirst, attachment); err != nil {
			return err
		}
	}

	return nil
}

// sendFile sends an individual file
func (s *Service) sendFile(chatID int64, conf *SendConfig, isFirst bool, attachment notify.Attachment) error {
	s.logger.Debug().Int64("recipient", chatID).Str("file", attachment.Name()).Msg("Sending file to chat")

	document := telegram.NewDocument(chatID, telegram.FileReader{
		Reader: attachment.Reader(),
		Name:   attachment.Name(),
	})

	// Set caption only for the first file
	if isFirst {
		document.Caption = conf.Message
		document.ParseMode = conf.ParseMode
	}

	// Quit early if dry run is enabled
	if conf.DryRun {
		s.logger.Info().Str("recipient", strconv.FormatInt(chatID, 10)).Msg("Dry run enabled - Message not sent.")
	}

	// Send the file
	if _, err := s.client.Send(document); err != nil {
		return err
	}

	s.logger.Info().Int64("recipient", chatID).Str("file", attachment.Name()).Msg("File sent to chat")

	return nil
}

// newSendConfig creates a new send config with default values.
func (s *Service) newSendConfig(subject, message string, opts ...notify.SendOption) *SendConfig {
	conf := &SendConfig{
		Subject:       subject,
		Message:       message,
		DryRun:        s.dryRun,
		ContinueOnErr: s.continueOnErr,
		ParseMode:     s.parseMode,
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

	handleError := func(chatID int64, err error) {
		// Append error info and log
		chatIDStr := strconv.FormatInt(chatID, 10)
		failedRecipients = append(failedRecipients, chatIDStr)
		errorList = append(errorList, asNotifyError(err))
		s.logger.Warn().Err(err).Str("recipient", chatIDStr).Msg("Error sending message to recipient")
	}

	for _, chatID := range s.chatIDs {
		// If context is cancelled, return error immediately
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := s.sendToChat(chatID, conf); err != nil {
			handleError(chatID, err) // Handle the error

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

	if len(s.chatIDs) == 0 {
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
