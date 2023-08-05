package slack

import (
	"github.com/nikoksr/onelog"
	"github.com/slack-go/slack"
)

// Option represents a function type used to configure the Slack service.
type Option = func(*Service)

// Set of functions to provide optional configuration for the Slack service.

// WithClient allows using a custom Slack client.
func WithClient(client *slack.Client) Option {
	return func(s *Service) {
		s.client = client
		s.logger.Debug().Msg("Slack client set")
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

// WithRecipients sets the recipients that should receive messages. You can add more recipients by calling AddRecipients.
func WithRecipients(channelIDs ...string) Option {
	return func(s *Service) {
		s.recipients = channelIDs
		s.logger.Debug().Int("count", len(channelIDs)).Int("total", len(s.recipients)).Msg("Recipients set")
	}
}

// WithEscapeMessage sets whether messages should be escaped or not before sending.
func WithEscapeMessage(escapeMessage bool) Option {
	return func(s *Service) {
		s.escapeMessage = escapeMessage
		s.logger.Debug().Bool("escape-message", escapeMessage).Msg("Escape message set")
	}
}
