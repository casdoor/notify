package telegram

import (
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/nikoksr/notify/v2"
)

var _ notify.Service = (*Service)(nil)

const (
	// ModeHTML is the default mode for sending messages.
	ModeHTML = telegram.ModeHTML
	// ModeMarkdown is the markdown mode for sending messages.
	ModeMarkdown = telegram.ModeMarkdown
)

func defaultMessageRenderer(conf SendConfig) string {
	var builder strings.Builder

	builder.WriteString(conf.subject)
	builder.WriteString("\n\n")
	builder.WriteString(conf.message)

	return builder.String()
}

// Service is the telegram service. It is used to send messages to Telegram chats.
type Service struct {
	client        *telegram.BotAPI
	chatIDs       []int64
	name          string
	renderMessage func(conf SendConfig) string

	// Send option fields
	parseMode string
}

// New creates a new telegram service. It returns an error if the telegram client could not be created.
func New(token string, opts ...Option) (*Service, error) {
	client, err := telegram.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	s := &Service{
		client:        client,
		name:          "telegram",
		renderMessage: defaultMessageRenderer,
		parseMode:     ModeMarkdown,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return s.name
}

// AddReceivers adds chat IDs that should receive messages.
func (s *Service) AddReceivers(chatIDs ...int64) {
	s.chatIDs = append(s.chatIDs, chatIDs...)
}

// Option is a function that can be used to configure the telegram service.
type Option = func(*Service)

// WithClient sets the telegram client. This is useful if you want to use a custom client.
func WithClient(client *telegram.BotAPI) Option {
	return func(s *Service) {
		s.client = client
	}
}

// WithReceivers sets the chat IDs that should receive messages. You can add more chat IDs by calling AddReceivers.
func WithReceivers(chatIDs ...int64) Option {
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
