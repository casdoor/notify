package twilio

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

	// Twilio specific
	SenderPhoneNumber string
	Recipients        []string
}

// SendWithSenderPhoneNumber is a send option that sets the sender phone number of the message.
func SendWithSenderPhoneNumber(senderPhoneNumber string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.SenderPhoneNumber = senderPhoneNumber
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

// newSendConfig creates a new send config with default values.
func (s *Service) newSendConfig(subject, message string, opts ...notify.SendOption) *SendConfig {
	conf := &SendConfig{
		SendConfig: &notify.SendConfig{
			Subject:       subject,
			Message:       message,
			DryRun:        s.dryRun,
			ContinueOnErr: s.continueOnErr,
		},
		SenderPhoneNumber: s.senderPhoneNumber,
		Recipients:        s.recipients,
	}

	for _, opt := range opts {
		opt(conf)
	}

	conf.Message = s.renderMessage(conf)

	return conf
}
