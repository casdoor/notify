package ntfy

import "net/http"

// Option is a function that can be used to configure the ntfy service.
type Option = func(*Service)

// WithClient sets the ntfy client. This is useful if you want to use a custom client.
func WithClient(client *http.Client) Option {
	return func(s *Service) {
		s.client = client
	}
}

// WithRecipients sets the topics that should receive messages. You can add more topics by calling AddRecipients.
func WithRecipients(topics ...string) Option {
	return func(s *Service) {
		s.AddRecipients(topics...)
	}
}

// WithName sets the name of the service. The default name is "ntfy".
func WithName(name string) Option {
	return func(s *Service) {
		s.name = name
	}
}

// WithMessageRenderer sets the message renderer. The default function will put the subject and message on separate lines.
//
// Example:
//
//	ntfy.WithMessageRenderer(func(conf SendConfig) string {
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
	}
}

// WithAPIBaseURL sets the API base URL. The default is "https://ntfy.sh/".
func WithAPIBaseURL(url string) Option {
	return func(s *Service) {
		s.apiBaseURL = url
	}
}

// WithParseMode sets the parse mode for sending messages. The default is ModeText.
func WithParseMode(mode Mode) Option {
	return func(s *Service) {
		s.parseMode = mode
	}
}

// WithPriority sets the priority for sending messages. The default is PriorityDefault.
func WithPriority(priority Priority) Option {
	return func(s *Service) {
		s.priority = priority
	}
}

// WithTags sets the tags for sending messages. The default is no tags.
func WithTags(tags ...string) Option {
	return func(s *Service) {
		s.tags = tags
	}
}

// WithIcon sets the icon for sending messages. The default is "" (no icon).
func WithIcon(icon string) Option {
	return func(s *Service) {
		s.icon = icon
	}
}

// WithDelay sets the delay for sending messages. The default is "" (no delay).
func WithDelay(delay string) Option {
	return func(s *Service) {
		s.delay = delay
	}
}

// WithClickAction sets the click action for sending messages. The default is "" (no click action).
func WithClickAction(action string) Option {
	return func(s *Service) {
		s.clickAction = action
	}
}
