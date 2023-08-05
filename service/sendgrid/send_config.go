package sendgrid

import (
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/nikoksr/notify/v2"
)

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

	// Sendgrid specific

	SenderAddress    string
	SenderName       string
	Recipients       []string
	CCRecipients     []string
	BCCRecipients    []string
	ParseMode        Mode
	Headers          map[string]string
	CustomArgs       map[string]string
	BatchID          string
	IPPoolID         string
	ASM              *mail.Asm
	MailSettings     *mail.MailSettings
	TrackingSettings *mail.TrackingSettings
}

// SendWithSenderAddress is a send option that sets the sender address of the message.
func SendWithSenderAddress(address string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.SenderAddress = address
		}
	}
}

// SendWithSenderName is a send option that sets the sender name of the message.
func SendWithSenderName(name string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.SenderName = name
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
func SendWithParseMode(mode Mode) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.ParseMode = mode
		}
	}
}

// SendWithHeaders is a send option that sets the headers of the message.
func SendWithHeaders(headers map[string]string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.Headers = headers
		}
	}
}

// SendWithCustomArgs is a send option that sets the custom args of the message.
func SendWithCustomArgs(customArgs map[string]string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.CustomArgs = customArgs
		}
	}
}

// SendWithBatchID is a send option that sets the batch id of the message.
func SendWithBatchID(batchID string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.BatchID = batchID
		}
	}
}

// SendWithIPoolID is a send option that sets the ipool id of the message.
func SendWithIPoolID(iPoolID string) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.IPPoolID = iPoolID
		}
	}
}

// SendWithASM is a send option that sets the asm of the message.
func SendWithASM(asm *mail.Asm) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.ASM = asm
		}
	}
}

// SendWithMailSettings is a send option that sets the mail settings of the message.
func SendWithMailSettings(mailSettings *mail.MailSettings) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.MailSettings = mailSettings
		}
	}
}

// SendWithTrackingSettings is a send option that sets the tracking settings of the message.
func SendWithTrackingSettings(trackingSettings *mail.TrackingSettings) notify.SendOption {
	return func(config notify.SendConfigurer) {
		if typedConf, ok := config.(*SendConfig); ok {
			typedConf.TrackingSettings = trackingSettings
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
		SenderAddress:    s.senderAddress,
		SenderName:       s.senderName,
		Recipients:       s.recipients,
		CCRecipients:     s.ccRecipients,
		BCCRecipients:    s.bccRecipients,
		ParseMode:        s.parseMode,
		Headers:          s.headers,
		CustomArgs:       s.customArgs,
		BatchID:          s.batchID,
		IPPoolID:         s.ipPoolID,
		ASM:              s.asm,
		MailSettings:     s.mailSettings,
		TrackingSettings: s.trackingSettings,
	}

	for _, opt := range opts {
		opt(conf)
	}

	conf.Message = s.renderMessage(conf)

	return conf
}
