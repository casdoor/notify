package mailgun

import (
	"github.com/mailgun/mailgun-go/v4"
	"github.com/nikoksr/onelog"
)

// Option is a function that can be used to configure the telegram service.
type Option = func(*Service)

// WithClient sets the telegram client. This is useful if you want to use a custom client.
func WithClient(client mailgun.Mailgun) Option {
	return func(s *Service) {
		s.client = client
		s.logger.Debug().Msg("Mailgun client set")
	}
}

// WithName sets the name of the service.
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
//	WithMessageRenderer(func(conf *SendConfig) string {
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

// WithSenderAddress sets the sender address. This is the address that will be shown as sender in the chat.
func WithSenderAddress(senderAddress string) Option {
	return func(s *Service) {
		s.senderAddress = senderAddress
		s.logger.Debug().Str("sender-address", senderAddress).Msg("Sender address set")
	}
}

// WithRecipients sets the email addresses that should receive messages. You can add more email addresses by calling AddRecipients.
func WithRecipients(phoneNumbers ...string) Option {
	return func(s *Service) {
		s.recipients = phoneNumbers
		s.logger.Debug().Int("count", len(phoneNumbers)).Int("total", len(s.recipients)).Msg("Recipients set")
	}
}

// WithCCRecipients sets the chat IDs that should receive messages as CC. You can add more chat IDs by calling
// AddCCRecipients.
func WithCCRecipients(recipients ...string) Option {
	return func(s *Service) {
		s.ccRecipients = recipients
		s.logger.Debug().Int("count", len(recipients)).Int("total", len(s.ccRecipients)).Msg("CC recipients set")
	}
}

// WithBCCRecipients sets the chat IDs that should receive messages as BCC. You can add more chat IDs by calling
// AddBCCRecipients.
func WithBCCRecipients(recipients ...string) Option {
	return func(s *Service) {
		s.bccRecipients = recipients
		s.logger.Debug().Int("count", len(recipients)).Int("total", len(s.bccRecipients)).Msg("BCC recipients set")
	}
}

// WithParseMode sets the parse mode. The default parse mode is ModeText.
func WithParseMode(parseMode Mode) Option {
	return func(s *Service) {
		s.parseMode = parseMode
		s.logger.Debug().Str("parse-mode", string(parseMode)).Msg("Parse mode set")
	}
}

// WithDomain specifies a separate domain for the type of messages being sent.
func WithDomain(domain string) Option {
	return func(s *Service) {
		s.domain = domain
		s.logger.Debug().Str("domain", domain).Msg("Domain set")
	}
}

// WithHeaders assigns custom MIME headers to the message being sent.
func WithHeaders(headers map[string]string) Option {
	return func(s *Service) {
		s.headers = headers
		s.logger.Debug().Int("count", len(headers)).Msg("Header set")
	}
}

// WithTags attaches tags to the message for metrics gathering and event tracking, as per Mailgun documentation.
func WithTags(tags ...string) Option {
	return func(s *Service) {
		s.tags = tags
		s.logger.Debug().Int("count", len(tags)).Msg("Tags set")
	}
}

// WithDKIM sets the 'o:dkim' header with its respective value for the sent message, as detailed in the Mailgun's
// documentation.
func WithDKIM(dkim bool) Option {
	return func(s *Service) {
		s.setDKIM = dkim
		s.logger.Debug().Bool("dkim", dkim).Msg("DKIM set")
	}
}

// WithNativeSend enables the return path to match the 'From' address in the Message.Headers when sending from Mailgun,
// instead of the default bounce+ address.
func WithNativeSend(nativeSend bool) Option {
	return func(s *Service) {
		s.enableNativeSend = nativeSend
		s.logger.Debug().Bool("nativeSend", nativeSend).Msg("NativeSend set")
	}
}

// WithRequireTLS sets the option according to TLS requirements as specified in the Mailgun's documentation.
func WithRequireTLS(requireTLS bool) Option {
	return func(s *Service) {
		s.requireTLS = requireTLS
		s.logger.Debug().Bool("requireTLS", requireTLS).Msg("RequireTLS set")
	}
}

// WithSkipVerification configures the verification setting as per instructions in the Mailgun's documentation.
func WithSkipVerification(skipVerification bool) Option {
	return func(s *Service) {
		s.skipVerification = skipVerification
		s.logger.Debug().Bool("skipVerification", skipVerification).Msg("SkipVerification set")
	}
}

// WithTestMode allows sending of a message that will be discarded by Mailgun, enabling client-side software testing
// without consuming email resources.
func WithTestMode(testMode bool) Option {
	return func(s *Service) {
		s.enableTestMode = testMode
		s.logger.Debug().Bool("testMode", testMode).Msg("TestMode set")
	}
}

// WithTracking adjusts the 'o:tracking' message parameter for each message, enabling or disabling URL rewriting to
// facilitate event tracking (such as opens, clicks, unsubscribes, etc.), according to the method's parameter. This
// header, however, isn't passed to the final recipient(s). More details can be found in the Mailgun's documentation.
func WithTracking(tracking bool) Option {
	return func(s *Service) {
		s.enableTracking = tracking
		s.logger.Debug().Bool("tracking", tracking).Msg("Tracking set")
	}
}

// WithTrackingClicks configures click tracking as per the details provided in the Mailgun's documentation.
func WithTrackingClicks(trackingClicks bool) Option {
	return func(s *Service) {
		s.enableTrackingClicks = trackingClicks
		s.logger.Debug().Bool("trackingClicks", trackingClicks).Msg("TrackingClicks set")
	}
}

// WithTrackingOpens determines the open tracking setting as specified in the Mailgun's documentation.
func WithTrackingOpens(trackingOpens bool) Option {
	return func(s *Service) {
		s.enableTrackingOpens = trackingOpens
		s.logger.Debug().Bool("trackingOpens", trackingOpens).Msg("TrackingOpens set")
	}
}
