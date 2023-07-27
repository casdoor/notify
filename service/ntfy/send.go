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

func (s *Service) sendTextMessage(ctx context.Context, topic string, conf SendConfig) error {
	s.logger.Debug().Str("topic", topic).Msg("Sending text message to topic")

	payload := &sendMessageRequest{
		Topic:       topic,
		Title:       conf.subject,
		Message:     conf.message,
		Tags:        conf.tags,
		Priority:    conf.priority,
		ClickAction: conf.clickAction,
		Markdown:    conf.parseMode == ModeMarkdown,
		Delay:       conf.delay,
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

	if err := sendRequest(s.client, req); err != nil {
		return err
	}

	s.logger.Info().Str("topic", topic).Msg("Text message sent to topic")

	return nil
}

func (s *Service) sendFile(ctx context.Context, topic string, attachment notify.Attachment) error {
	s.logger.Debug().Str("topic", topic).Str("file", attachment.Name()).Msg("Sending file to topic")

	// Append topic to base URL, e.g. https://ntfy.sh/my_topic
	endpoint := s.apiBaseURL + topic

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, attachment)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	// req.Header.Set("X-Title", attachment.Name())

	if err := sendRequest(s.client, req); err != nil {
		return err
	}

	s.logger.Info().Str("topic", topic).Str("file", attachment.Name()).Msg("File sent to topic")

	return nil
}

func (s *Service) sendFileAttachments(ctx context.Context, topic string, conf SendConfig) error {
	for _, attachment := range conf.attachments {
		if err := s.sendFile(ctx, topic, attachment); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) sendToTopic(ctx context.Context, topic string, conf SendConfig) error {
	if err := s.sendTextMessage(ctx, topic, conf); err != nil {
		return err
	}

	return s.sendFileAttachments(ctx, topic, conf)
}

// Send sends a message to all topics that are configured to receive messages. It returns an error if the message could
// not be sent.
func (s *Service) Send(ctx context.Context, subject, message string, opts ...notify.SendOption) error {
	if len(s.topics) == 0 {
		return notify.ErrNoRecipients
	}

	conf := SendConfig{
		subject:     subject,
		message:     message,
		parseMode:   s.parseMode,
		priority:    s.priority,
		tags:        s.tags,
		delay:       s.delay,
		clickAction: s.clickAction,
	}

	for _, opt := range opts {
		opt(&conf)
	}

	conf.message = s.renderMessage(conf)

	if conf.message == "" && len(conf.attachments) == 0 {
		s.logger.Warn().Msg("Message is empty and no attachments are present. Aborting send.")
		return nil
	}

	s.logger.Debug().Msg("Sending message to topics")

	for _, topic := range s.topics {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := s.sendToTopic(ctx, topic, conf); err != nil {
			return &notify.SendNotificationError{
				Recipient: topic,
				Cause:     err, // asNotifyError has been called in sendToTopic, as it requires the http response
			}
		}
	}

	s.logger.Info().Msg("Message successfully sent to all topics")

	return nil
}
