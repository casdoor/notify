package slack

import (
	"context"

	"github.com/slack-go/slack"

	"github.com/nikoksr/notify/v2"
)

func (s *Service) sendFile(ctx context.Context, channelID string, conf SendConfig, isFirst bool, attachment notify.Attachment) error {
	s.logger.Debug().Str("recipient", channelID).Str("file", attachment.Name()).Msg("Sending file to channel")

	params := slack.UploadFileV2Parameters{
		Reader:   attachment,
		Filename: attachment.Name(),
		AltTxt:   attachment.Name(),
		Channel:  channelID,
	}

	if isFirst {
		params.InitialComment = conf.Message
	}

	if _, err := s.client.UploadFileV2Context(ctx, params); err != nil {
		return err
	}

	s.logger.Info().Str("recipient", channelID).Str("file", attachment.Name()).Msg("File sent to channel")

	return nil
}

func (s *Service) sendFileAttachments(ctx context.Context, channelID string, conf SendConfig) error {
	for idx, attachment := range conf.Attachments {
		isFirst := idx == 0
		if err := s.sendFile(ctx, channelID, conf, isFirst, attachment); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) sendTextMessage(ctx context.Context, channelID string, conf SendConfig) error {
	s.logger.Debug().Str("recipient", channelID).Msg("Sending text message to channel")

	if _, _, err := s.client.PostMessageContext(ctx, channelID, slack.MsgOptionText(conf.Message, conf.EscapeMessage)); err != nil {
		return err
	}

	s.logger.Info().Str("recipient", channelID).Msg("Text message sent to channel")

	return nil
}

// sendToChannel sends a message to a specific channel utilizing the SendConfig settings. If no message or attachments
// are defined, the function will return without error.
func (s *Service) sendToChannel(ctx context.Context, channelID string, conf SendConfig) error {
	if len(conf.Attachments) == 0 {
		return s.sendTextMessage(ctx, channelID, conf)
	}

	return s.sendFileAttachments(ctx, channelID, conf)
}

// Send sends a notification to each Slack channel defined in Service. The sender is configured through SendOption and
// SendConfig. Returns an error upon failure to send the message, or if there are no recipients identified.
func (s *Service) Send(ctx context.Context, subject, message string, opts ...notify.SendOption) error {
	if len(s.channelIDs) == 0 {
		return notify.ErrNoRecipients
	}

	conf := SendConfig{
		Subject:       subject,
		Message:       message,
		EscapeMessage: s.escapeMessage,
	}

	for _, opt := range opts {
		opt(&conf)
	}

	conf.Message = s.renderMessage(conf)

	if conf.Message == "" && len(conf.Attachments) == 0 {
		s.logger.Warn().Msg("Message is empty and no attachments are present. Aborting send.")
		return nil
	}

	s.logger.Debug().Msg("Sending message to all recipients")

	for _, channelID := range s.channelIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := s.sendToChannel(ctx, channelID, conf); err != nil {
			return &notify.SendNotificationError{
				Recipient: channelID,
				Cause:     asNotifyError(err),
			}
		}
	}

	s.logger.Info().Msg("Message successfully sent to all recipients")

	return nil
}
