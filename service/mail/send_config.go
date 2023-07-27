package mail

import "github.com/nikoksr/notify/v2"

var _ notify.SendConfig = (*SendConfig)(nil)

// SendConfig represents the configuration needed for sending a message.
//
// This struct complies with the notify.SendConfig interface and allows you to alter
// the behavior of the send function. This can be achieved by either passing send options
// to the send function or by manipulating the fields of this struct in your custom
// message renderer.
//
// All fields of this struct are exported to offer maximum flexibility to users.
// However, users must be aware that they are responsible for managing thread-safety
// and other similar concerns when manipulating these fields directly.
type SendConfig struct {
	Subject     string
	Message     string
	Attachments []notify.Attachment
	Metadata    map[string]any

	// Mail specific fields

	ParseMode         Mode
	Priority          Priority
	SenderName        string
	InlineAttachments bool
}

// SetAttachments adds attachments to the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetAttachments(attachments ...notify.Attachment) {
	c.Attachments = attachments
}

// SetMetadata sets the metadata of the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetMetadata(metadata map[string]any) {
	c.Metadata = metadata
}

// Send options

// SendWithParseMode is a send option that sets the parse mode of the message.
func SendWithParseMode(parseMode Mode) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.ParseMode = parseMode
		}
	}
}

// SendWithPriority is a send option that sets the priority of the message.
func SendWithPriority(priority Priority) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.Priority = priority
		}
	}
}

// SendWithSenderName is a send option that sets the sender name of the message. This will be displayed in the
// recipient's email client. E.g. "From Example <john.doe@example>", where "From Example" is the sender name. The default
// is "From Notify <no-reply>".
func SendWithSenderName(senderName string) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.SenderName = senderName
		}
	}
}

// SendWithInlineAttachments is a send option that sets whether attachments should be sent inline or not. If set to
// true, attachments will be sent inline. This is a Mail specific option. The default is false.
func SendWithInlineAttachments(inlineAttachments bool) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.InlineAttachments = inlineAttachments
		}
	}
}
