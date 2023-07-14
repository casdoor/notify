package slack

import (
	"strings"

	"github.com/slack-go/slack"

	"github.com/nikoksr/notify/v2"
)

// Compile time check to make sure the service implements the notify.Service interface.
var _ notify.Service = (*Service)(nil)

// defaultMessageRenderer is a helper function to format messages for Slack.
func defaultMessageRenderer(conf SendConfig) string {
	var builder strings.Builder

	builder.WriteString(conf.subject)
	builder.WriteString("\n\n")
	builder.WriteString(conf.message)

	return builder.String()
}

// Service is a structure that contains data needed for interaction with Slack's APIs.
type Service struct {
	// client is the Slack client used for API requests.
	client *slack.Client

	// channelIDs represents the Slack channels messages will be sent to.
	channelIDs []string

	// name is a descriptive identifier for the service, by default "slack".
	name string

	// renderMessage is the function used to format messages.
	renderMessage func(conf SendConfig) string

	// Slack specific fields

	// escapeMessage is a flag used to escape characters in messages that have special meanings in Slack's markup.
	escapeMessage bool
}

// New creates a new instance of the Slack service with a default configuration. It receives as input the required Slack token and optional configurations.
func New(token string, opts ...Option) *Service {
	client := slack.New(token)

	s := &Service{
		client:        client,
		name:          "slack",
		renderMessage: defaultMessageRenderer,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Name returns the name of the service, which identifies the type of service in use (in this case, Slack).
func (s *Service) Name() string {
	return s.name
}

// AddReceivers appends given channel IDs onto an internal list that Send uses to distribute the notifications.
func (s *Service) AddReceivers(channelIDs ...string) {
	s.channelIDs = append(s.channelIDs, channelIDs...)
}

// Option represents a function type used to configure the Slack service.
type Option = func(*Service)

// Set of functions to provide optional configuration for the Slack service.

// WithClient allows using a custom Slack client.
func WithClient(client *slack.Client) Option {
	return func(t *Service) {
		t.client = client
	}
}

// WithReceivers sets the default recipients for the notifications on the service.
func WithReceivers(channelIDs ...string) Option {
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
