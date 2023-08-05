package telegram

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nikoksr/onelog"
)

// Option is a function that can be used to configure the telegram service.
type Option = func(*Service)

// WithClient sets the telegram client. This is useful if you want to use a custom client.
func WithClient(client *telegram.BotAPI) Option {
	return func(s *Service) {
		s.client = client
		s.logger.Debug().Msg("Telegram client set")
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
func WithRecipients(chatIDs ...int64) Option {
	return func(s *Service) {
		s.chatIDs = chatIDs
		s.logger.Debug().Int("count", len(chatIDs)).Int("total", len(s.chatIDs)).Msg("Recipients set")
	}
}

// WithParseMode sets the parse mode for sending messages. The default is ModeHTML.
func WithParseMode(mode string) Option {
	return func(s *Service) {
		s.parseMode = mode
		s.logger.Debug().Str("mode", mode).Msg("Parse mode set")
	}
}
