package sendgrid

import (
	"github.com/nikoksr/onelog"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Option is a function that can be used to configure the sendgrid service.
type Option = func(*Service)

// WithClient sets the sendgrid client. This is useful if you want to use a custom client.
func WithClient(client *sendgrid.Client) Option {
	return func(s *Service) {
		s.client = client
		s.logger.Debug().Msg("Sendgrid client set")
	}
}

// WithName sets the name of the service. The default name is "sendgrid".
func WithName(name string) Option {
	return func(s *Service) {
		s.name = name
		s.logger.Debug().Str("name", name).Msg("Service name set")
	}
}

// WithLogger sets the logger. The default logger is a no-op logger.
func WithLogger(logger onelog.Logger) Option {
	return func(s *Service) {
		logger = logger.With("service", s.Name()) // Add service name to logger
		s.logger = logger
		s.logger.Debug().Msg("Logger set")
	}
}

// WithMessageRenderer sets the message renderer. The default function will put the subject and message on separate lines.
//
// Example:
//
//	sendgrid.WithMessageRenderer(func(conf *SendConfig) string {
//		var builder strings.Builder
//
//		builder.WriteString(conf.subject)
//		builder.WriteString("\n")
//		builder.WriteString(conf.message)
//
//		return builder.String()
//	})
func WithMessageRenderer(builder func(conf *SendConfig) string) Option {
	return func(s *Service) {
		s.renderMessage = builder
		s.logger.Debug().Msg("Message renderer set")
	}
}

// WithDryRun sets the dry run flag. If set to true, no messages will be sent.
func WithDryRun(dryRun bool) Option {
	return func(s *Service) {
		s.dryRun = dryRun
		s.logger.Debug().Bool("dry-run", dryRun).Msg("Dry run set")
	}
}

// WithContinueOnErr sets the continue on error flag. If set to true, the service will continue sending the message to
// the next recipient even if an error occurred.
func WithContinueOnErr(continueOnErr bool) Option {
	return func(s *Service) {
		s.continueOnErr = continueOnErr
		s.logger.Debug().Bool("continue-on-error", continueOnErr).Msg("Continue on error set")
	}
}

// WithRecipients sets the email addresses that should receive messages. You can add more email addresses by calling AddRecipients.
func WithRecipients(phoneNumbers ...string) Option {
	return func(s *Service) {
		s.recipients = phoneNumbers
		s.logger.Debug().Int("count", len(phoneNumbers)).Int("total", len(s.recipients)).Msg("Recipients set")
	}
}

// WithCCRecipients sets the email addresses that should receive messages as CC. You can add more email addresses by
// calling AddCCRecipients.
func WithCCRecipients(phoneNumbers ...string) Option {
	return func(s *Service) {
		s.ccRecipients = phoneNumbers
		s.logger.Debug().Int("count", len(phoneNumbers)).Int("total", len(s.ccRecipients)).Msg("CC recipients set")
	}
}

// WithBCCRecipients sets the email addresses that should receive messages as BCC. You can add more email addresses by
// calling AddBCCRecipients.
func WithBCCRecipients(phoneNumbers ...string) Option {
	return func(s *Service) {
		s.bccRecipients = phoneNumbers
		s.logger.Debug().Int("count", len(phoneNumbers)).Int("total", len(s.bccRecipients)).Msg("BCC recipients set")
	}
}

// WithSenderAddress sets the sender address.
func WithSenderAddress(address string) Option {
	return func(s *Service) {
		s.senderAddress = address
		s.logger.Debug().Str("sender-address", address).Msg("Sender address set")
	}
}

// WithSenderName sets the sender name.
func WithSenderName(name string) Option {
	return func(s *Service) {
		s.senderName = name
		s.logger.Debug().Str("sender-name", name).Msg("Sender name set")
	}
}

// WithParseMode sets the parse mode. The default parse mode is ModeHTML.
func WithParseMode(mode Mode) Option {
	return func(s *Service) {
		s.parseMode = mode
		s.logger.Debug().Str("parse-mode", string(mode)).Msg("Parse mode set")
	}
}

// WithHeaders sets the Email headers.
func WithHeaders(headers map[string]string) Option {
	return func(s *Service) {
		s.headers = headers
		s.logger.Debug().Int("count", len(headers)).Msg("Headers set")
	}
}

// WithCustomArgs sets custom arguments for SendGrid.
func WithCustomArgs(customArgs map[string]string) Option {
	return func(s *Service) {
		s.customArgs = customArgs
		s.logger.Debug().Int("count", len(customArgs)).Msg("Custom args set")
	}
}

// WithBatchID sets the batch ID for the Email.
func WithBatchID(batchID string) Option {
	return func(s *Service) {
		s.batchID = batchID
		s.logger.Debug().Str("batch-id", batchID).Msg("Batch ID set")
	}
}

// WithIPPoolID sets the IP Pool ID for the Email.
func WithIPPoolID(ipPoolID string) Option {
	return func(s *Service) {
		s.ipPoolID = ipPoolID
		s.logger.Debug().Str("ip-pool-id", ipPoolID).Msg("IP Pool ID set")
	}
}

// WithASM sets the ASM for the Email.
func WithASM(asm *mail.Asm) Option {
	return func(s *Service) {
		s.asm = asm
		s.logger.Debug().Msg("ASM set")
	}
}

// WithMailSettings sets the MailSettings for the Email.
func WithMailSettings(mailSettings *mail.MailSettings) Option {
	return func(s *Service) {
		s.mailSettings = mailSettings
		s.logger.Debug().Msg("MailSettings set")
	}
}

// WithTrackingSettings sets the TrackingSettings for the Email.
func WithTrackingSettings(trackingSettings *mail.TrackingSettings) Option {
	return func(s *Service) {
		s.trackingSettings = trackingSettings
		s.logger.Debug().Msg("TrackingSettings set")
	}
}
