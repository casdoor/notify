package mailgun

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

	// Mailgun specific

	SenderAddress        string
	Recipients           []string
	CCRecipients         []string
	BCCRecipients        []string
	ParseMode            Mode
	Domain               string
	Headers              map[string]string
	Tags                 []string
	SetDKIM              bool
	EnableNativeSend     bool
	RequireTLS           bool
	SkipVerification     bool
	EnableTestMode       bool
	EnableTracking       bool
	EnableTrackingClicks bool
	EnableTrackingOpens  bool
}

// SendWithSenderAddress sets the sender address of the message.
func SendWithSenderAddress(senderAddress string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.SenderAddress = senderAddress
		}
	}
}

// SendWithRecipients sets the recipients of the message.
func SendWithRecipients(recipients ...string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.Recipients = recipients
		}
	}
}

// SendWithCCRecipients sets the CC recipients of the message.
func SendWithCCRecipients(ccRecipients ...string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.CCRecipients = ccRecipients
		}
	}
}

// SendWithBCCRecipients sets the BCC recipients of the message.
func SendWithBCCRecipients(bccRecipients ...string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.BCCRecipients = bccRecipients
		}
	}
}

// SendWithParseMode sets the parse mode of the message.
func SendWithParseMode(mode Mode) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.ParseMode = mode
		}
	}
}

// SendWithDomain sets the domain of the message.
func SendWithDomain(domain string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.Domain = domain
		}
	}
}

// SendWithHeaders sets the header of the message.
func SendWithHeaders(headers map[string]string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.Headers = headers
		}
	}
}

// SendWithTags sets the tags of the message.
func SendWithTags(tags []string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.Tags = tags
		}
	}
}

// SendWithSetDKIM sets the DKIM flag of the message.
func SendWithSetDKIM(setDKIM bool) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.SetDKIM = setDKIM
		}
	}
}

// SendWithEnableNativeSend sets the native send flag of the message.
func SendWithEnableNativeSend(enableNativeSend bool) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.EnableNativeSend = enableNativeSend
		}
	}
}

// SendWithRequireTLS sets the TLS flag of the message.
func SendWithRequireTLS(requireTLS bool) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.RequireTLS = requireTLS
		}
	}
}

// SendWithSkipVerification sets the verification flag of the message.
func SendWithSkipVerification(skipVerification bool) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.SkipVerification = skipVerification
		}
	}
}

// SendWithEnableTestMode sets the test mode flag of the message.
func SendWithEnableTestMode(enableTestMode bool) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.EnableTestMode = enableTestMode
		}
	}
}

// SendWithEnableTracking sets the tracking flag of the message.
func SendWithEnableTracking(enableTracking bool) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.EnableTracking = enableTracking
		}
	}
}

// SendWithEnableTrackingClicks sets the tracking clicks flag of the message.
func SendWithEnableTrackingClicks(enableTrackingClicks bool) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.EnableTrackingClicks = enableTrackingClicks
		}
	}
}

// SendWithEnableTrackingOpens sets the tracking opens flag of the message.
func SendWithEnableTrackingOpens(enableTrackingOpens bool) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.EnableTrackingOpens = enableTrackingOpens
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
		SenderAddress:        s.senderAddress,
		Recipients:           s.recipients,
		CCRecipients:         s.ccRecipients,
		BCCRecipients:        s.bccRecipients,
		ParseMode:            s.parseMode,
		Domain:               s.domain,
		Headers:              s.headers,
		Tags:                 s.tags,
		SetDKIM:              s.setDKIM,
		EnableNativeSend:     s.enableNativeSend,
		RequireTLS:           s.requireTLS,
		SkipVerification:     s.skipVerification,
		EnableTestMode:       s.enableTestMode,
		EnableTracking:       s.enableTracking,
		EnableTrackingClicks: s.enableTrackingClicks,
		EnableTrackingOpens:  s.enableTrackingOpens,
	}

	for _, opt := range opts {
		opt(conf)
	}

	conf.Message = s.renderMessage(conf)

	return conf
}
