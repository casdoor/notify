package slack

import "github.com/slack-go/slack"

// Option represents a function type used to configure the Slack service.
type Option = func(*Service)

// Set of functions to provide optional configuration for the Slack service.

// WithClient allows using a custom Slack client.
func WithClient(client *slack.Client) Option {
	return func(t *Service) {
		t.client = client
	}
}

// WithRecipients sets the default recipients for the notifications on the service.
func WithRecipients(channelIDs ...string) Option {
	return func(t *Service) {
		t.channelIDs = channelIDs
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
//	slack.WithMessageRenderer(func(conf SendConfig) string {
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

// WithEscapeMessage sets whether messages should be escaped or not before sending.
func WithEscapeMessage(escapeMessage bool) Option {
	return func(t *Service) {
		t.escapeMessage = escapeMessage
	}
}
