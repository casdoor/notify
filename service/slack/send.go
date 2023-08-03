package slack

import (
	"context"

	"github.com/slack-go/slack"

	"github.com/nikoksr/notify/v2"
)

func (s *Service) sendFile(ctx context.Context, channelID string, conf *SendConfig, isFirst bool, attachment notify.Attachment) error {
	s.logger.Debug().Str("recipient", channelID).Str("file", attachment.Name()).Msg("Sending file to channel")

	params := slack.UploadFileV2Parameters{
		Reader:   attachment.Reader(),
		Filename: attachment.Name(),
		FileSize: int(attachment.Size()),
		AltTxt:   attachment.Name(),
		Channel:  channelID,
	}

	if isFirst {
		params.InitialComment = conf.Message
	}

	// Quit early if dry run is enabled
	if conf.DryRun {
		s.logger.Info().Str("recipient", channelID).Str("file", attachment.Name()).Msg("Dry run enabled - File not sent.")
		return nil
	}

	// Send the file
	if _, err := s.client.UploadFileV2Context(ctx, params); err != nil {
		return err
	}

	s.logger.Info().Str("recipient", channelID).Str("file", attachment.Name()).Msg("File sent to channel")

	return nil
}

func (s *Service) sendFileAttachments(ctx context.Context, channelID string, conf *SendConfig) error {
	for idx, attachment := range conf.Attachments {
		isFirst := idx == 0
		if err := s.sendFile(ctx, channelID, conf, isFirst, attachment); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) sendTextMessage(ctx context.Context, channelID string, conf *SendConfig) error {
	s.logger.Debug().Str("recipient", channelID).Msg("Sending text message to channel")

	// Quit early if dry run is enabled
	if conf.DryRun {
		s.logger.Info().Str("recipient", channelID).Msg("Dry run enabled - Message not sent.")
		return nil
	}

	// Send the message
	if _, _, err := s.client.PostMessageContext(ctx, channelID, slack.MsgOptionText(conf.Message, conf.EscapeMessage)); err != nil {
		return err
	}

	s.logger.Info().Str("recipient", channelID).Msg("Text message sent to channel")

	return nil
}

// sendToChannel sends a message to a specific channel utilizing the SendConfig settings. If no message or attachments
// are defined, the function will return without error.
func (s *Service) sendToChannel(ctx context.Context, channelID string, conf *SendConfig) error {
	if len(conf.Attachments) == 0 {
		return s.sendTextMessage(ctx, channelID, conf)
	}

	return s.sendFileAttachments(ctx, channelID, conf)
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

	handleError := func(channelID string, err error) {
		// Append error info and log
		failedRecipients = append(failedRecipients, channelID)
		errorList = append(errorList, asNotifyError(err))
		s.logger.Warn().Err(err).Str("recipient", channelID).Msg("Error sending message to recipient")
	}

	for _, channelID := range s.channelIDs {
		// If context is cancelled, return error immediately
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := s.sendToChannel(ctx, channelID, conf); err != nil {
			handleError(channelID, err) // Handle the error

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

// Send sends a notification to each Slack channel defined in Service. The sender is configured through SendOption and
// SendConfig. Returns an error upon failure to send the message, or if there are no recipients identified.
func (s *Service) Send(ctx context.Context, subject, message string, opts ...notify.SendOption) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.channelIDs) == 0 {
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
