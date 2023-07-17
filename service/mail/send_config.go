package mail

import "github.com/nikoksr/notify/v2"

// Compile time check to make sure the struct implements the notify.SendConfig interface.
var _ notify.SendConfig = (*SendConfig)(nil)

// SendConfig constructs the configuration for sending notifications over Mail and the struct implements the
// notify.SendConfig interface. It facilitates features such as subject, message, attachments and metadata. Mail-
// specific fields are also accommodated.
type SendConfig struct {
	subject     string
	message     string
	attachments []notify.Attachment
	metadata    map[string]any

	// Mail specific fields
	parseMode         Mode
	priority          Priority
	senderName        string
	inlineAttachments bool
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

// Mail specific fields

// ParseMode returns the parse mode of the message. This is a Mail specific option.
func (c *SendConfig) ParseMode() Mode {
	return c.parseMode
}

// Priority returns the priority of the message. This is a Mail specific option.
func (c *SendConfig) Priority() Priority {
	return c.priority
}

// SenderName returns the sender name of the message. This is a Mail specific option.
func (c *SendConfig) SenderName() string {
	return c.senderName
}

// InlineAttachments returns whether attachments should be sent inline or not. This is a Mail specific option.
func (c *SendConfig) InlineAttachments() bool {
	return c.inlineAttachments
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

// SendWithParseMode is a send option that sets the parse mode of the message.
func SendWithParseMode(parseMode Mode) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.parseMode = parseMode
		}
	}
}

// SendWithPriority is a send option that sets the priority of the message.
func SendWithPriority(priority Priority) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.priority = priority
		}
	}
}

// SendWithSenderName is a send option that sets the sender name of the message. This will be displayed in the
// recipient's email client. E.g. "From Example <john.doe@example>", where "From Example" is the sender name. The default
// is "From Notify <no-reply>".
func SendWithSenderName(senderName string) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.senderName = senderName
		}
	}
}

// SendWithInlineAttachments is a send option that sets whether attachments should be sent inline or not. If set to
// true, attachments will be sent inline. This is a Mail specific option. The default is false.
func SendWithInlineAttachments(inlineAttachments bool) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.inlineAttachments = inlineAttachments
		}
	}
}
