package slack

import "github.com/nikoksr/notify/v2"

// Compile time check to make sure the struct implements the notify.SendConfig interface.
var _ notify.SendConfig = (*SendConfig)(nil)

// SendConfig constructs the configuration for sending notifications over Slack and the struct implements the
// notify.SendConfig interface. It facilitates features such as subject, message, attachments and metadata. Slack-
// specific fields are also accommodated.
type SendConfig struct {
	subject     string
	message     string
	attachments []notify.Attachment
	metadata    map[string]any

	// Slack specific fields
	escapeMessage bool
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

// Slack specific fields

// EscapeMessage is a function to determine whether the message content should be escaped or not.
func (c *SendConfig) EscapeMessage() bool {
	return c.escapeMessage
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

// Send options

// SendWithEscapeMessage sets the escape message option.
func SendWithEscapeMessage(escapeMessage bool) notify.SendOption {
	return func(c notify.SendConfig) {
		if conf, ok := c.(*SendConfig); ok {
			conf.escapeMessage = escapeMessage
		}
	}
}
