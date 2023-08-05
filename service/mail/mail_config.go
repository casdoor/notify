package mail

import (
	"crypto/tls"
	"time"

	"github.com/nikoksr/onelog"

	mail "github.com/xhit/go-simple-mail/v2"
)

// Option represents a function type used to configure the Mail service.
type Option = func(*Service)

// Set of functions to provide optional configuration for the Mail service.

// WithServer allows using a custom Mail server. This should usually not be used, as the default server created in the
// New function should suffice and offer enough flexibility. However, to avoid any future inconveniences, this option
// is provided.
func WithServer(server *mail.SMTPServer) Option {
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.server = server
		s.logger.Debug().Msg("Mail server set")
	}
}

// WithClient allows using a custom Mail client.
func WithClient(client *mail.SMTPClient) Option {
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.client = client
		s.logger.Debug().Msg("Mail client set")
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

// WithRecipients sets the email addresses that should receive messages. You can add more email addresses by calling AddRecipients.
func WithRecipients(phoneNumbers ...string) Option {
	return func(s *Service) {
		s.recipients = phoneNumbers
		s.logger.Debug().Int("count", len(phoneNumbers)).Int("total", len(s.recipients)).Msg("Recipients set")
	}
}

// WithParseMode sets the parse mode for sending messages. The default is ModeHTML.
func WithParseMode(mode Mode) Option {
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.parseMode = mode
		s.logger.Debug().Int("mode", int(mode)).Msg("Parse mode set")
	}
}

// WithPriority sets the priority for sending messages. The default is PriorityNormal.
func WithPriority(priority Priority) Option {
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.priority = priority
		s.logger.Debug().Int("priority", int(priority)).Msg("Priority set")
	}
}

// WithSenderName sets the sender name for the notifications on the service. The default is "From Notify <no-reply>".
func WithSenderName(name string) Option {
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.senderName = name
		s.logger.Debug().Str("name", name).Msg("Sender name set")
	}
}

// WithKeepAlive sets the keep alive option on the Mail server. KeepAlive is enabled by default.
func WithKeepAlive(keepAlive bool) Option {
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.server.KeepAlive = keepAlive
		s.logger.Debug().Bool("keepAlive", keepAlive).Msg("Keep alive set")
	}
}

// WithConnectTimeout sets the connect timeout option on the Mail server.
func WithConnectTimeout(timeout time.Duration) Option {
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.server.ConnectTimeout = timeout
		s.logger.Debug().Dur("timeout", timeout).Msg("Connect timeout set")
	}
}

// WithSendTimeout sets the send timeout option on the Mail server.
func WithSendTimeout(timeout time.Duration) Option {
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.server.SendTimeout = timeout
		s.logger.Debug().Dur("timeout", timeout).Msg("Send timeout set")
	}
}

// WithEncryption sets the encryption option on the Mail server. The default is EncryptionNone.
func WithEncryption(encryption mail.Encryption) Option {
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.server.Encryption = encryption
		s.logger.Debug().Int("encryption", int(encryption)).Msg("Encryption set")
	}
}

// WithTLSConfig sets the TLS config option on the Mail server.
func WithTLSConfig(config *tls.Config) Option {
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.server.TLSConfig = config
		s.logger.Debug().Msg("TLS config set")
	}
}

// WithAuthType sets the authentication type on the Mail server. The default is AuthAuto, and it is usually not
// necessary to change this.
func WithAuthType(authentication mail.AuthType) Option {
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.server.Authentication = authentication
		s.logger.Debug().Int("authentication", int(authentication)).Msg("Authentication set")
	}
}

// WithHelo sets the HELO option on the Mail server. The default is localhost. HELO is the hostname that the client
// sends to the server when the connection is established.
func WithHelo(helo string) Option {
	return func(s *Service) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.server.Helo = helo
		s.logger.Debug().Str("helo", helo).Msg("HELO set")
	}
}
