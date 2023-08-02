package ntfy

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/nikoksr/notify/v2"
)

func sendRequest(client *http.Client, req *http.Request) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := asNotifyError(resp); err != nil {
		return err
	}

	return nil
}

func (s *Service) sendTextMessage(ctx context.Context, topic string, conf *SendConfig) error {
	s.logger.Debug().Str("recipient", topic).Msg("Sending text message to topic")

	payload := &sendMessageRequest{
		Topic:       topic,
		Title:       conf.Subject,
		Message:     conf.Message,
		Tags:        conf.Tags,
		Priority:    conf.Priority,
		ClickAction: conf.ClickAction,
		Markdown:    conf.ParseMode == ModeMarkdown,
		Delay:       conf.Delay,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.apiBaseURL, bytes.NewReader(payloadJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	// Quit early if dry run is enabled
	if conf.DryRun {
		s.logger.Info().Str("recipient", topic).Msg("Dry run enabled - Message not sent.")
		return nil
	}

	// Send the message
	if err := sendRequest(s.client, req); err != nil {
		return err
	}

	s.logger.Info().Str("recipient", topic).Msg("Text message sent to topic")

	return nil
}

func (s *Service) sendFile(ctx context.Context, topic string, conf *SendConfig, attachment notify.Attachment) error {
	s.logger.Debug().Str("recipient", topic).Str("file", attachment.Name()).Msg("Sending file to topic")

	// Append topic to base URL, e.g. https://ntfy.sh/my_topic
	endpoint := s.apiBaseURL + topic

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, attachment.Reader())
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)

	// Quit early if dry run is enabled
	if conf.DryRun {
		s.logger.Info().Str("recipient", topic).Str("file", attachment.Name()).Msg("Dry run enabled - File not sent.")
		return nil
	}

	// Send the file
	if err := sendRequest(s.client, req); err != nil {
		return err
	}

	s.logger.Info().Str("recipient", topic).Str("file", attachment.Name()).Msg("File sent to topic")

	return nil
}

func (s *Service) sendFileAttachments(ctx context.Context, topic string, conf *SendConfig) error {
	for _, attachment := range conf.Attachments {
		if err := s.sendFile(ctx, topic, conf, attachment); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) sendToTopic(ctx context.Context, topic string, conf *SendConfig) error {
	if err := s.sendTextMessage(ctx, topic, conf); err != nil {
		return err
	}

	return s.sendFileAttachments(ctx, topic, conf)
}

// newSendConfig creates a new send config with default values.
func (s *Service) newSendConfig(subject, message string, opts ...notify.SendOption) *SendConfig {
	conf := &SendConfig{
		Subject:       subject,
		Message:       message,
		DryRun:        s.dryRun,
		ContinueOnErr: s.continueOnErr,
		ParseMode:     s.parseMode,
		Priority:      s.priority,
		Tags:          s.tags,
		Delay:         s.delay,
		ClickAction:   s.clickAction,
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

	handleError := func(topic string, err error) {
		// Append error info and log
		failedRecipients = append(failedRecipients, topic)
		errorList = append(errorList, err) // asNotifyError has been called in sendToTopic, as it requires the http response
		s.logger.Warn().Err(err).Str("recipient", topic).Msg("Error sending message to recipient")
	}

	for _, topic := range s.topics {
		// If context is cancelled, return error immediately
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := s.sendToTopic(ctx, topic, conf); err != nil {
			handleError(topic, err) // Handle the error

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

// Send sends a message to all topics that are configured to receive messages. It returns an error if the message could
// not be sent.
func (s *Service) Send(ctx context.Context, subject, message string, opts ...notify.SendOption) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.topics) == 0 {
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
