package mail

import "github.com/nikoksr/notify/v2"

var _ notify.SendConfigurer = (*SendConfig)(nil)

// SendConfig represents the configuration needed for sending a message.
//
// This struct complies with the notify.SendConfigurer interface and allows you to alter
// the behavior of the send function. This can be achieved by either passing send options
// to the send function or by manipulating the fields of this struct in your custom
// message renderer.
//
// All fields of this struct are exported to offer maximum flexibility to users.
// However, users must be aware that they are responsible for managing thread-safety
// and other similar concerns when manipulating these fields directly.
type SendConfig struct {
	*notify.SendConfig

	// Mail specific

	SenderName    string
	Recipients    []string
	CCRecipients  []string
	BCCRecipients []string
	ParseMode     Mode
	Priority      Priority
}

// SendWithSenderName is a send option that sets the sender name of the message. This will be displayed in the
// recipient's email client. E.g. "From Example <john.doe@example>", where "From Example" is the sender name.
func SendWithSenderName(senderName string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.SenderName = senderName
		}
	}
}

// SendWithRecipients is a send option that sets the recipients of the message.
func SendWithRecipients(recipients ...string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.Recipients = recipients
		}
	}
}

// SendWithCCRecipients is a send option that sets the cc recipients of the message.
func SendWithCCRecipients(ccRecipients ...string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.CCRecipients = ccRecipients
		}
	}
}

// SendWithBCCRecipients is a send option that sets the bcc recipients of the message.
func SendWithBCCRecipients(bccRecipients ...string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.BCCRecipients = bccRecipients
		}
	}
}

// SendWithParseMode is a send option that sets the parse mode of the message.
func SendWithParseMode(parseMode Mode) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.ParseMode = parseMode
		}
	}
}

// SendWithPriority is a send option that sets the priority of the message.
func SendWithPriority(priority Priority) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.Priority = priority
		}
	}
}

// newSendConfig creates a new send config with default values.
func (s *Service) newSendConfig(subject, message string, opts ...notify.SendOption) *SendConfig {
	conf := &SendConfig{
		SendConfig: &notify.SendConfig{
			Subject:       subject,
			Message:       message,
			DryRun:        s.dryRun,
			ContinueOnErr: s.continueOnErr,
		},
		Recipients:    s.recipients,
		CCRecipients:  s.ccRecipients,
		BCCRecipients: s.bccRecipients,
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
