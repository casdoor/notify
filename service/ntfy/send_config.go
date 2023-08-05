package ntfy

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

	// Ntfy specific

	APIBaseURL  string
	APIKey      string
	Recipients  []string
	ParseMode   Mode
	Priority    Priority
	Tags        []string
	Delay       string
	ClickAction string
}

// SendWithAPIBaseURL is a send option that sets the API base URL of the message.
func SendWithAPIBaseURL(apiBaseURL string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.APIBaseURL = apiBaseURL
		}
	}
}

// SendWithAPIKey is a send option that sets the API key of the message.
func SendWithAPIKey(apiKey string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.APIKey = apiKey
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

// SendWithTags is a send option that sets the tags of the message.
func SendWithTags(tags ...string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.Tags = tags
		}
	}
}

// SendWithDelay is a send option that sets the delay of the message.
func SendWithDelay(delay string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.Delay = delay
		}
	}
}

// SendWithClickAction is a send option that sets the click action of the message.
func SendWithClickAction(clickAction string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.ClickAction = clickAction
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
		APIBaseURL:  s.apiBaseURL,
		APIKey:      s.apiKey,
		Recipients:  s.recipients,
		ParseMode:   s.parseMode,
		Priority:    s.priority,
		Tags:        s.tags,
		Delay:       s.delay,
		ClickAction: s.clickAction,
	}

	for _, opt := range opts {
		opt(conf)
	}

	conf.Message = s.renderMessage(conf)

	return conf
}
