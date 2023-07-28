package telegram

import (
	"context"

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
		Reader: attachment,
		Name:   attachment.Name(),
	})

	// Set caption only for the first file
	if isFirst {
		document.Caption = conf.Message
		document.ParseMode = conf.ParseMode
	}

	if _, err := s.client.Send(document); err != nil {
		return err
	}

	s.logger.Info().Int64("recipient", chatID).Str("file", attachment.Name()).Msg("File sent to chat")

	return nil
}

// newSendConfig creates a new send config with default values.
func (s *Service) newSendConfig(subject, message string, opts ...notify.SendOption) *SendConfig {
	conf := &SendConfig{
		Subject:   subject,
		Message:   message,
		DryRun:    s.dryRun,
		ParseMode: s.parseMode,
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

	for _, chatID := range s.chatIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := s.sendToChat(chatID, conf); err != nil {
			return &notify.SendNotificationError{
				Recipient: chatID,
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

	if conf.DryRun {
		s.logger.Info().Str("message", conf.Message).Msg("Dry run enabled - Message not sent.")
		return nil
	}

	// Send message to all recipients
	return s.send(ctx, conf)
}
