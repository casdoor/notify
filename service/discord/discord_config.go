package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nikoksr/onelog"
)

// Option is a function that applies an option to the service.
type Option = func(*Service)

// WithClient sets the discord client (session) to use for sending messages. Naming it WithClient to stay consistent
// with the other services.
func WithClient(session *discordgo.Session) Option {
	return func(s *Service) {
		s.client.setSession(session)
		s.logger.Info().Msg("Discord client set")
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

// WithRecipients sets the channel IDs or webhook URLs to send messages to. You can add more channel IDs or webhook URLs
// by calling Service.AddRecipients.
func WithRecipients(recipients ...string) Option {
	return func(d *Service) {
		d.recipients = recipients
		d.logger.Info().Int("count", len(recipients)).Int("total", len(d.recipients)).Msg("Recipients set")
	}
}

// WithName sets the name of the service. The default is "discord".
func WithName(name string) Option {
	return func(d *Service) {
		d.name = name
		d.logger.Info().Str("name", name).Msg("Service name set")
	}
}

// WithMessageRenderer sets the message renderer. The default function will put the subject and message on separate lines.
//
// Example:
//
//	telegram.WithMessageRenderer(func(conf *SendConfig) string {
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
