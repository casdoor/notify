package mail

import (
	"crypto/tls"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

// Option represents a function type used to configure the Mail service.
type Option = func(*Service)

// Set of functions to provide optional configuration for the Mail service.

// WithServer allows using a custom Mail server. This should usually not be used, as the default server created in the
// New function should suffice and offer enough flexibility. However, to avoid any future inconveniences, this option
// is provided.
func WithServer(server *mail.SMTPServer) Option {
	return func(t *Service) {
		t.server = server
	}
}

// WithClient allows using a custom Mail client.
func WithClient(client *mail.SMTPClient) Option {
	return func(t *Service) {
		t.client = client
	}
}

// WithRecipients sets the default recipients for the notifications on the service.
func WithRecipients(recipients ...string) Option {
	return func(t *Service) {
		t.recipients = recipients
	}
}

// WithCCRecipients sets the default CC recipients for the notifications on the service.
func WithCCRecipients(recipients ...string) Option {
	return func(t *Service) {
		t.ccRecipients = recipients
	}
}

// WithBCCRecipients sets the default BCC recipients for the notifications on the service.
func WithBCCRecipients(recipients ...string) Option {
	return func(t *Service) {
		t.bccRecipients = recipients
	}
}

// WithName sets an alternative name for the service.
func WithName(name string) Option {
	return func(t *Service) {
		t.name = name
	}
}

// WithMessageRenderer sets the function to render the message.
//
// Example:
//
//	email.WithMessageRenderer(func(conf SendConfig) string {
//		var builder strings.Builder
//
//		builder.WriteString(conf.subject)
//		builder.WriteString("\n")
//		builder.WriteString(conf.message)
//
//		return builder.String()
//	})
func WithMessageRenderer(builder func(conf SendConfig) string) Option {
	return func(t *Service) {
		t.renderMessage = builder
	}
}

// WithParseMode sets the parse mode for sending messages. The default is ModeHTML.
func WithParseMode(mode Mode) Option {
	return func(t *Service) {
		t.parseMode = mode
	}
}

// WithPriority sets the priority for sending messages. The default is PriorityNormal.
func WithPriority(priority Priority) Option {
	return func(t *Service) {
		t.priority = priority
	}
}

// WithSenderName sets the sender name for the notifications on the service. The default is "From Notify <no-reply>".
func WithSenderName(name string) Option {
	return func(t *Service) {
		t.senderName = name
	}
}

// WithInlineAttachments sets the inline attachments option on the Mail server. The default is false.
func WithInlineAttachments(inline bool) Option {
	return func(t *Service) {
		t.inlineAttachments = inline
	}
}

// WithKeepAlive sets the keep alive option on the Mail server. KeepAlive is enabled by default.
func WithKeepAlive(keepAlive bool) Option {
	return func(t *Service) {
		t.server.KeepAlive = keepAlive
	}
}

// WithConnectTimeout sets the connect timeout option on the Mail server.
func WithConnectTimeout(timeout time.Duration) Option {
	return func(t *Service) {
		t.server.ConnectTimeout = timeout
	}
}

// WithSendTimeout sets the send timeout option on the Mail server.
func WithSendTimeout(timeout time.Duration) Option {
	return func(t *Service) {
		t.server.SendTimeout = timeout
	}
}

// WithEncryption sets the encryption option on the Mail server. The default is EncryptionNone.
func WithEncryption(encryption mail.Encryption) Option {
	return func(t *Service) {
		t.server.Encryption = encryption
	}
}

// WithTLSConfig sets the TLS config option on the Mail server.
func WithTLSConfig(config *tls.Config) Option {
	return func(t *Service) {
		t.server.TLSConfig = config
	}
}

// WithAuthType sets the authentication type on the Mail server. The default is AuthAuto, and it is usually not
// necessary to change this.
func WithAuthType(authentication mail.AuthType) Option {
	return func(t *Service) {
		t.server.Authentication = authentication
	}
}

// WithHelo sets the HELO option on the Mail server. The default is localhost. HELO is the hostname that the client
// sends to the server when the connection is established.
func WithHelo(helo string) Option {
	return func(t *Service) {
		t.server.Helo = helo
	}
}
