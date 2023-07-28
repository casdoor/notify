package twilio

import (
	"github.com/nikoksr/onelog"
	"github.com/twilio/twilio-go"
)

// Option is a function that can be used to configure the twilio service.
type Option = func(*Service)

// WithClient sets the twilio client. This is useful if you want to use a custom client.
func WithClient(client *twilio.RestClient) Option {
	return func(s *Service) {
		s.client = client
		s.logger.Info().Msg("Twilio client set")
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

// WithRecipients sets the phonenumbers that should receive messages. You can add more phonenumbers by calling AddRecipients.
func WithRecipients(phoneNumbers ...string) Option {
	return func(s *Service) {
		s.phoneNumbers = phoneNumbers
		s.logger.Info().Int("count", len(phoneNumbers)).Int("total", len(s.phoneNumbers)).Msg("Recipients set")
	}
}

// WithName sets the name of the service. The default name is "twilio".
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
//	twilio.WithMessageRenderer(func(conf *SendConfig) string {
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

// WithDryRun sets the dry run flag. If set to true, no messages will be sent.
func WithDryRun(dryRun bool) Option {
	return func(s *Service) {
		s.dryRun = dryRun
		s.logger.Info().Bool("dry-run", dryRun).Msg("Dry run set")
	}
}

// WithEdge sets the default edge to use.
func WithEdge(edge string) Option {
	return func(s *Service) {
		s.client.SetEdge(edge)
		s.logger.Info().Str("edge", edge).Msg("Edge set")
	}
}

// WithRegion sets the default region to use.
func WithRegion(region string) Option {
	return func(s *Service) {
		s.client.SetRegion(region)
		s.logger.Info().Str("region", region).Msg("Region set")
	}
}
