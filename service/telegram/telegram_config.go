package telegram

import telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Option is a function that can be used to configure the telegram service.
type Option = func(*Service)

// WithClient sets the telegram client. This is useful if you want to use a custom client.
func WithClient(client *telegram.BotAPI) Option {
	return func(s *Service) {
		s.client = client
	}
}

// WithRecipients sets the chat IDs that should receive messages. You can add more chat IDs by calling AddRecipients.
func WithRecipients(chatIDs ...int64) Option {
	return func(s *Service) {
		s.chatIDs = chatIDs
	}
}

// WithName sets the name of the service. The default name is "telegram".
func WithName(name string) Option {
	return func(s *Service) {
		s.name = name
	}
}

// WithMessageRenderer sets the message renderer. The default function will put the subject and message on separate lines.
//
// Example:
//
//	telegram.WithMessageRenderer(func(conf SendConfig) string {
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

// WithParseMode sets the parse mode for sending messages. The default is ModeHTML.
func WithParseMode(mode string) Option {
	return func(s *Service) {
		s.parseMode = mode
	}
}
