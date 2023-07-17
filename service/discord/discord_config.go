package discord

import "github.com/bwmarrin/discordgo"

// Option is a function that applies an option to the service.
type Option = func(*Service)

// WithClient sets the discord client (session) to use for sending messages. Naming it WithClient to stay consistent
// with the other services.
func WithClient(session *discordgo.Session) Option {
	return func(s *Service) {
		s.client.setSession(session)
	}
}

// WithRecipients sets the channel IDs or webhook URLs to send messages to. You can add more channel IDs or webhook URLs
// by calling Service.AddRecipients.
func WithRecipients(recipients ...string) Option {
	return func(d *Service) {
		d.recipients = recipients
	}
}

// WithName sets the name of the service. The default is "discord".
func WithName(name string) Option {
	return func(d *Service) {
		d.name = name
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
	return func(t *Service) {
		t.renderMessage = builder
	}
}
