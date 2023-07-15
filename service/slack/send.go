package slack

import (
	"context"

	"github.com/slack-go/slack"

	"github.com/nikoksr/notify/v2"
)

func (s *Service) sendFile(ctx context.Context, channelID string, conf SendConfig, isFirst bool, attachment notify.Attachment) error {
	params := slack.UploadFileV2Parameters{
		Reader:   attachment,
		Filename: attachment.Name(),
		AltTxt:   attachment.Name(),
		Channel:  channelID,
	}

	if isFirst {
		params.InitialComment = conf.message
	}

	_, err := s.client.UploadFileV2Context(ctx, params)
	return err
}

func (s *Service) sendFileAttachments(ctx context.Context, channelID string, conf SendConfig) error {
	for idx, attachment := range conf.attachments {
		isFirst := idx == 0
		if err := s.sendFile(ctx, channelID, conf, isFirst, attachment); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) sendTextMessage(ctx context.Context, channelID string, conf SendConfig) error {
	_, _, err := s.client.PostMessageContext(ctx, channelID, slack.MsgOptionText(conf.message, conf.escapeMessage))
	return err
}

// sendToChannel sends a message to a specific channel utilizing the SendConfig settings. If no message or attachments
// are defined, the function will return without error.
func (s *Service) sendToChannel(ctx context.Context, channelID string, conf SendConfig) error {
	if conf.message == "" {
		return nil
	}

	// decide the way of sending the message based on attachments
	if len(conf.attachments) == 0 {
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
		subject:       subject,
		message:       message,
		escapeMessage: s.escapeMessage,
	}

	for _, opt := range opts {
		opt(&conf)
	}

	conf.message = s.renderMessage(conf)

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

	return nil
}
