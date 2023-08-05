package slack

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

	// Slack specific

	Recipients    []string
	EscapeMessage bool
}

// SendWithRecipients sets the recipients option.
func SendWithRecipients(recipients ...string) notify.SendOption {
	return func(c notify.SendConfigurer) {
		if conf, ok := c.(*SendConfig); ok {
			conf.Recipients = recipients
		}
	}
}

// SendWithEscapeMessage sets the escape message option.
func SendWithEscapeMessage(escapeMessage bool) notify.SendOption {
	return func(c notify.SendConfigurer) {
		if conf, ok := c.(*SendConfig); ok {
			conf.EscapeMessage = escapeMessage
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
		EscapeMessage: s.escapeMessage,
	}

	for _, opt := range opts {
		opt(conf)
	}

	conf.Message = s.renderMessage(conf)

	return conf
}
