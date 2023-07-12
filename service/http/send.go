package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/nikoksr/notify/v2"
)

var _ notify.SendConfig = (*SendConfig)(nil)

// SendConfig is the configuration for sending a message. It implements the
// notify.SendConfig interface.
type SendConfig struct {
	subject     string
	message     string
	attachments []notify.Attachment
	metadata    map[string]any
}

// Common fields

// Subject returns the subject of the message.
func (c *SendConfig) Subject() string {
	return c.subject
}

// Message returns the message.
func (c *SendConfig) Message() string {
	return c.message
}

// Attachments returns the attachments.
func (c *SendConfig) Attachments() []notify.Attachment {
	return c.attachments
}

// Metadata returns the metadata.
func (c *SendConfig) Metadata() map[string]any {
	return c.metadata
}

// notify.SendConfig implementation

// SetAttachments adds attachments to the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetAttachments(attachments ...notify.Attachment) {
	c.attachments = attachments
}

// SetMetadata sets the metadata of the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetMetadata(metadata map[string]any) {
	c.metadata = metadata
}

// newRequest creates a new http request with the given method, content-type, url and payload. Request created by this
// function will usually be passed to the Service.do method.
func newRequest(ctx context.Context, hook *Webhook, payload io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, hook.Method, hook.URL, payload)
	if err != nil {
		return nil, err
	}

	req.Header = hook.Header

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", defaultUserAgent)
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", hook.ContentType)
	}

	return req, nil
}

// do sends the given request and returns an error if the request failed. A failed request gets identified by either
// an unsuccessful status code or a non-nil error. The given request is expected to be valid and was usually created
// by the newRequest function.
func (s *Service) do(req *http.Request) error {
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	// Check if response code is 2xx. Should this be configurable?
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("responded with status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *Service) sendToWebhook(ctx context.Context, webhook *Webhook, conf SendConfig) error {
	payload := s.renderMessage(conf)

	req, err := newRequest(ctx, webhook, bytes.NewReader([]byte(payload)))
	if err != nil {
		return fmt.Errorf("create request for %q: %w", webhook, err)
	}
	defer func() { _ = req.Body.Close() }()

	return s.do(req)
}

// Send takes a message and sends it to all webhooks.
func (s *Service) Send(ctx context.Context, subject, message string, opts ...notify.SendOption) error {
	if len(s.webhooks) == 0 {
		return notify.ErrNoRecipients
	}

	conf := SendConfig{
		subject: subject,
		message: message,
	}

	for _, opt := range opts {
		opt(&conf)
	}

	conf.message = s.renderMessage(conf)

	for _, webhook := range s.webhooks {
		if err := s.sendToWebhook(ctx, webhook, conf); err != nil {
			return notify.NewErrSendNotification(webhook, err)
		}
	}

	return nil
}
