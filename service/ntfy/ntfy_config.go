package ntfy

import (
	"net/http"

	"github.com/nikoksr/onelog"
)

// Option is a function that can be used to configure the ntfy service.
type Option = func(*Service)

// WithClient sets the ntfy client. This is useful if you want to use a custom client.
func WithClient(client *http.Client) Option {
	return func(s *Service) {
		s.client = client
		s.logger.Info().Msg("Ntfy client set")
	}
}

// WithLogger sets the logger. The default logger is a no-op logger.
func WithLogger(logger onelog.Logger) Option {
	return func(s *Service) {
		logger = logger.With("service", s.Name()) // Add service name to logger
		s.logger = logger
		s.logger.Info().Msg("Logger set")
	}
}

// WithRecipients sets the topics that should receive messages. You can add more topics by calling AddRecipients.
func WithRecipients(topics ...string) Option {
	return func(s *Service) {
		s.AddRecipients(topics...)
		s.logger.Info().Int("count", len(topics)).Int("total", len(s.topics)).Msg("Recipients set")
	}
}

// WithName sets the name of the service. The default name is "ntfy".
func WithName(name string) Option {
	return func(s *Service) {
		s.name = name
		s.logger.Info().Str("name", name).Msg("Service name set")
	}
}

// WithMessageRenderer sets the message renderer. The default function will put the subject and message on separate lines.
//
// Example:
//
//	ntfy.WithMessageRenderer(func(conf *SendConfig) string {
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
		s.logger.Info().Msg("Message renderer set")
	}
}

// WithAPIBaseURL sets the API base URL. The default is "https://ntfy.sh/".
func WithAPIBaseURL(url string) Option {
	return func(s *Service) {
		s.apiBaseURL = url
		s.logger.Info().Str("url", url).Msg("API base URL set")
	}
}

// WithParseMode sets the parse mode for sending messages. The default is ModeText.
func WithParseMode(mode Mode) Option {
	return func(s *Service) {
		s.parseMode = mode
		s.logger.Info().Str("mode", string(mode)).Msg("Parse mode set")
	}
}

// WithPriority sets the priority for sending messages. The default is PriorityDefault.
func WithPriority(priority Priority) Option {
	return func(s *Service) {
		s.priority = priority
		s.logger.Info().Int("priority", int(priority)).Msg("Priority set")
	}
}

// WithTags sets the tags for sending messages. The default is no tags.
func WithTags(tags ...string) Option {
	return func(s *Service) {
		s.tags = tags
		s.logger.Info().Int("count", len(tags)).Int("total", len(s.tags)).Msg("Tags set")
	}
}

// WithIcon sets the icon for sending messages. The default is "" (no icon).
func WithIcon(icon string) Option {
	return func(s *Service) {
		s.icon = icon
		s.logger.Info().Str("icon", icon).Msg("Icon set")
	}
}

// WithDelay sets the delay for sending messages. The default is "" (no delay).
func WithDelay(delay string) Option {
	return func(s *Service) {
		s.delay = delay
		s.logger.Info().Str("delay", delay).Msg("Delay set")
	}
}

// WithClickAction sets the click action for sending messages. The default is "" (no click action).
func WithClickAction(action string) Option {
	return func(s *Service) {
		s.clickAction = action
		s.logger.Info().Str("clickAction", action).Msg("Click action set")
	}
}
