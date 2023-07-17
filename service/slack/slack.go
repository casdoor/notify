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

// New creates a new instance of the Slack service with a default configuration. It receives as input the required Slack
// token and optional configurations. If no configuration is provided, the default values are used.
//
// Note: This function never returns an error. It has a return value for consistency with other services.
func New(token string, opts ...Option) (*Service, error) {
	client := slack.New(token)

	s := &Service{
		client:        client,
		name:          "slack",
		renderMessage: defaultMessageRenderer,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

// Name returns the name of the service, which identifies the type of service in use (in this case, Slack).
func (s *Service) Name() string {
	return s.name
}

// AddRecipients appends given channel IDs onto an internal list that Send uses to distribute the notifications.
func (s *Service) AddRecipients(channelIDs ...string) {
	s.channelIDs = append(s.channelIDs, channelIDs...)
}
