package ntfy

import (
	"net/http"
	"sync"

	"github.com/nikoksr/onelog"
	nopadapter "github.com/nikoksr/onelog/adapter/nop"

	"github.com/nikoksr/notify/v2"
)

var _ notify.Service = (*Service)(nil)

type (
	// Priority is the priority for sending messages. The default priority is PriorityDefault.
	Priority int
	// Mode is the mode for sending messages. It is a string to allow for custom parsing modes. The default mode is
	// ModeText.
	Mode string

	clickAction struct {
		Action string `json:"action"`
		Label  string `json:"label"`
		URL    string `json:"url"`
	}

	sendMessageRequest struct {
		Topic       string        `json:"topic"`
		Title       string        `json:"title"`
		Message     string        `json:"message"`
		Tags        []string      `json:"tags"`
		Priority    Priority      `json:"priority"`
		Actions     []clickAction `json:"actions"`
		ClickAction string        `json:"click"`
		Markdown    bool          `json:"markdown"`
		Delay       string        `json:"delay"`
	}
)

const (
	// PriorityMin is the minimum priority for sending messages.
	PriorityMin Priority = iota + 1
	// PriorityLow is the low priority for sending messages.
	PriorityLow
	// PriorityDefault is the default priority for sending messages.
	PriorityDefault
	// PriorityHigh is the high priority for sending messages.
	PriorityHigh
	// PriorityMax is the maximum priority for sending messages.
	PriorityMax

	// ModeText is treating messages as plain text. It is the default mode.
	ModeText Mode = "text/plain"
	// ModeMarkdown is formatting messages with Markdown.
	ModeMarkdown Mode = "text/markdown"

	defaultAPIBaseURL = "https://ntfy.sh/"
)

func defaultMessageRenderer(conf *SendConfig) string {
	return conf.Message
}

// Service is the ntfy service. It is used to send messages to Ntfy chats.
type Service struct {
	client *http.Client

	name          string
	mu            sync.RWMutex
	logger        onelog.Logger
	renderMessage func(conf *SendConfig) string
	dryRun        bool

	// Ntfy specific
	token       string
	topics      []string
	apiBaseURL  string
	parseMode   Mode
	priority    Priority
	tags        []string
	icon        string
	delay       string
	clickAction string
}

func (s *Service) applyOptions(opts ...Option) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, opt := range opts {
		opt(s)
	}
}

// New creates a new ntfy service. It returns an error if the ntfy client could not be created.
func New(token string, opts ...Option) (*Service, error) {
	s := &Service{
		client:        http.DefaultClient,
		logger:        nopadapter.NewAdapter(),
		name:          "ntfy",
		token:         token,
		renderMessage: defaultMessageRenderer,
		dryRun:        false,
		apiBaseURL:    defaultAPIBaseURL,
		parseMode:     ModeText,
		priority:      PriorityDefault,
	}

	s.applyOptions(opts...)

	return s, nil
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return s.name
}

// AddRecipients adds topics that should receive messages.
func (s *Service) AddRecipients(topics ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.topics = append(s.topics, topics...)
	s.logger.Info().Int("count", len(topics)).Int("total", len(s.topics)).Msg("Recipients added")
}
