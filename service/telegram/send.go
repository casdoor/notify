package telegram

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/nikoksr/notify/v2"
)

// sendToChat sends a message to a chat. It returns an error if the message could not be sent.
func (s *Service) sendToChat(chatID int64, conf SendConfig) error {
	if conf.message == "" {
		return nil
	}

	if len(conf.attachments) == 0 {
		return s.sendTextMessage(chatID, conf)
	}

	return s.sendFileAttachments(chatID, conf)
}

// sendTextMessage sends a text message
func (s *Service) sendTextMessage(chatID int64, conf SendConfig) error {
	message := telegram.NewMessage(chatID, conf.message)
	message.ParseMode = conf.parseMode

	_, err := s.client.Send(message)
	return err
}

// sendFileAttachments sends file attachments
func (s *Service) sendFileAttachments(chatID int64, conf SendConfig) error {
	for idx, attachment := range conf.attachments {
		isFirst := idx == 0
		if err := s.sendFile(chatID, conf, isFirst, attachment); err != nil {
			return err
		}
	}

	return nil
}

// sendFile sends an individual file
func (s *Service) sendFile(chatID int64, conf SendConfig, isFirst bool, attachment notify.Attachment) error {
	document := telegram.NewDocument(chatID, telegram.FileReader{
		Reader: attachment,
		Name:   attachment.Name(),
	})

	// Set caption only for the first file
	if isFirst {
		document.Caption = conf.message
		document.ParseMode = conf.parseMode
	}

	_, err := s.client.Send(document)
	return err
}

// Send sends a message to all chats that are configured to receive messages. It returns an error if the message could
// not be sent.
func (s *Service) Send(ctx context.Context, subject, message string, opts ...notify.SendOption) error {
	if len(s.chatIDs) == 0 {
		return notify.ErrNoRecipients
	}

	conf := SendConfig{
		parseMode: s.parseMode,
		subject:   subject,
		message:   message,
	}

	for _, opt := range opts {
		opt(&conf)
	}

	conf.message = s.renderMessage(conf)

	for _, chatID := range s.chatIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := s.sendToChat(chatID, conf); err != nil {
				return &notify.ErrSendNotification{Recipient: chatID, Cause: err}
			}
		}
	}

	return nil
}
