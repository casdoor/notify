package telegram

import (
	"context"
	"fmt"
	"strconv"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/nikoksr/notify/v2"
)

func (s *Service) buildMessagePayload(chatID int64, conf *SendConfig) telegram.Chattable {
	message := telegram.NewMessage(chatID, conf.Message)
	message.ParseMode = conf.ParseMode

	return message
}

// sendTextMessage sends a text message
func (s *Service) sendTextMessage(chatID int64, conf *SendConfig) error {
	s.logger.Debug().Int64("recipient", chatID).Msg("Sending text message to chat")

	// Build the message
	message := s.buildMessagePayload(chatID, conf)

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

func (s *Service) buildFilePayload(chatID int64, conf *SendConfig, isFirst bool, attachment notify.Attachment) (telegram.Chattable, error) {
	document := telegram.NewDocument(chatID, telegram.FileReader{
		Reader: attachment.Reader(),
		Name:   attachment.Name(),
	})

	// Set caption only for the first file
	if isFirst {
		document.Caption = conf.Message
		document.ParseMode = conf.ParseMode
	}

	return document, nil
}

// sendFile sends an individual file
func (s *Service) sendFile(chatID int64, conf *SendConfig, isFirst bool, attachment notify.Attachment) error {
	s.logger.Debug().Int64("recipient", chatID).Str("file", attachment.Name()).Msg("Sending file to chat")

	// Build the payload
	file, err := s.buildFilePayload(chatID, conf, isFirst, attachment)
	if err != nil {
		return fmt.Errorf("build file payload: %w", err)
	}

	// Quit early if dry run is enabled
	if conf.DryRun {
		s.logger.Info().Str("recipient", strconv.FormatInt(chatID, 10)).Msg("Dry run enabled - Message not sent.")
	}

	// Send the file
	if _, err := s.client.Send(file); err != nil {
		return err
	}

	s.logger.Info().Int64("recipient", chatID).Str("file", attachment.Name()).Msg("File sent to chat")

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

// sendToChat sends a message to a chat. It returns an error if the message could not be sent.
func (s *Service) sendToChat(chatID int64, conf *SendConfig) error {
	if len(conf.Attachments) == 0 {
		return s.sendTextMessage(chatID, conf)
	}

	return s.sendFileAttachments(chatID, conf)
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

	for _, chatID := range conf.Recipients {
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
