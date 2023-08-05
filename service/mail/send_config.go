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
	Subject       string
	Message       string
	Attachments   []notify.Attachment
	Metadata      map[string]any
	DryRun        bool
	ContinueOnErr bool

	// Mail specific

	SenderName string
	ParseMode  Mode
	Priority   Priority
}

// SetAttachments adds attachments to the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetAttachments(attachments ...notify.Attachment) {
	c.Attachments = attachments
}

// SetMetadata sets the metadata of the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetMetadata(metadata map[string]any) {
	c.Metadata = metadata
}

// SetDryRun sets the dry run flag of the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetDryRun(dryRun bool) {
	c.DryRun = dryRun
}

// SetContinueOnErr sets the continue on error flag of the message. This method is needed as part of the
// notify.SendConfig interface. Compared to other services, this is a no-op, as the Mail service will always send its
// messages to all recipients at once.
func (c *SendConfig) SetContinueOnErr(continueOnErr bool) {
	c.ContinueOnErr = continueOnErr
}

// Send options

// SendWithSenderName is a send option that sets the sender name of the message. This will be displayed in the
// recipient's email client. E.g. "From Example <john.doe@example>", where "From Example" is the sender name.
func SendWithSenderName(senderName string) notify.SendOption {
	return func(config notify.SendConfig) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.SenderName = senderName
		}
	}
}

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

// newSendConfig creates a new send config with default values.
func (s *Service) newSendConfig(subject, message string, opts ...notify.SendOption) *SendConfig {
	conf := &SendConfig{
		Subject:       subject,
		Message:       message,
		DryRun:        s.dryRun,
		ContinueOnErr: s.continueOnErr,
		SenderName:    s.senderName,
		ParseMode:     s.parseMode,
		Priority:      s.priority,
	}

	for _, opt := range opts {
		opt(conf)
	}

	conf.Message = s.renderMessage(conf)

	return conf
}
