package slack

import "github.com/slack-go/slack"

// Option represents a function type used to configure the Slack service.
type Option = func(*Service)

// Set of functions to provide optional configuration for the Slack service.

// WithClient allows using a custom Slack client.
func WithClient(client *slack.Client) Option {
	return func(s *Service) {
		s.client = client
		s.logger.Info().Msg("Slack client set")
	}
}

// WithRecipients sets the default recipients for the notifications on the service.
func WithRecipients(channelIDs ...string) Option {
	return func(s *Service) {
		s.channelIDs = channelIDs
		s.logger.Info().Int("count", len(channelIDs)).Int("total", len(s.channelIDs)).Msg("Recipients set")
	}
}

// WithName sets an alternative name for the service.
func WithName(name string) Option {
	return func(s *Service) {
		s.name = name
		s.logger.Info().Str("name", name).Msg("Service name set")
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
	return func(s *Service) {
		s.renderMessage = builder
		s.logger.Info().Msg("Message renderer set")
	}
}

// WithEscapeMessage sets whether messages should be escaped or not before sending.
func WithEscapeMessage(escapeMessage bool) Option {
	return func(s *Service) {
		s.escapeMessage = escapeMessage
		s.logger.Info().Bool("escapeMessage", escapeMessage).Msg("Escape message set")
	}
}
