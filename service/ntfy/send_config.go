package ntfy

import "github.com/nikoksr/notify/v2"

var _ notify.SendConfig = (*SendConfig)(nil)

// SendConfig is the configuration for sending a message. It implements the
// notify.SendConfig interface.
type SendConfig struct {
	subject     string
	message     string
	attachments []notify.Attachment
	metadata    map[string]any

	// Custom fields
	parseMode   Mode
	priority    Priority
	tags        []string
	delay       string
	clickAction string
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

// Ntfy specific fields

// ParseMode returns the parse mode of the message. This is a Ntfy specific option.
func (c *SendConfig) ParseMode() Mode {
	return c.parseMode
}

// Priority returns the priority of the message. This is a Ntfy specific option.
func (c *SendConfig) Priority() Priority {
	return c.priority
}

// Tags returns the tags of the message. This is a Ntfy specific option.
func (c *SendConfig) Tags() []string {
	return c.tags
}

// Delay returns the delay of the message. This is a Ntfy specific option.
func (c *SendConfig) Delay() string {
	return c.delay
}

// ClickAction returns the click action of the message. This is a Ntfy specific option.
func (c *SendConfig) ClickAction() string {
	return c.clickAction
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

// SendWithPriority is a send option that sets the priority of the message.
func SendWithPriority(priority Priority) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.priority = priority
		}
	}
}

// SendWithParseMode is a send option that sets the parse mode of the message.
func SendWithParseMode(parseMode Mode) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.parseMode = parseMode
		}
	}
}

// SendWithTags is a send option that sets the tags of the message.
func SendWithTags(tags ...string) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.tags = tags
		}
	}
}

// SendWithDelay is a send option that sets the delay of the message.
func SendWithDelay(delay string) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.delay = delay
		}
	}
}

// SendWithClickAction is a send option that sets the click action of the message.
func SendWithClickAction(clickAction string) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.clickAction = clickAction
		}
	}
}
