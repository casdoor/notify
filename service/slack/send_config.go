package slack

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

	// Slack specific fields

	EscapeMessage bool
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

// SendWithEscapeMessage sets the escape message option.
func SendWithEscapeMessage(escapeMessage bool) notify.SendOption {
	return func(c notify.SendConfig) {
		if conf, ok := c.(*SendConfig); ok {
			conf.EscapeMessage = escapeMessage
		}
	}
}