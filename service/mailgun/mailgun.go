package mailgun

import (
	"strings"
	"sync"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/nikoksr/onelog"
	nopadapter "github.com/nikoksr/onelog/adapter/nop"

	"github.com/nikoksr/notify/v2"
)

var _ notify.Service = (*Service)(nil)

// Mode is the mode for sending messages. It is a string to allow for custom parsing modes. The default mode is
// ModeText.
type Mode string

const (
	// ModeText is treating messages as plain text. It is the default mode.
	ModeText Mode = "text/plain"
	// ModeHTML is formatting messages with HTML.
	ModeHTML Mode = "text/html"
	// ModeAMPHTML is formatting messages with AMP HTML.
	ModeAMPHTML Mode = "text/amphtml"
)

func defaultMessageRenderer(conf *SendConfig) string {
	var builder strings.Builder

	builder.WriteString(conf.Subject)
	builder.WriteString("\n\n")
	builder.WriteString(conf.Message)

	return builder.String()
}

// Service is the telegram service. It is used to send messages to Mailgun chats.
type Service struct {
	client mailgun.Mailgun

	name          string
	mu            sync.RWMutex
	logger        onelog.Logger
	renderMessage func(conf *SendConfig) string
	dryRun        bool
	continueOnErr bool

	// Mailgun specific

	senderAddress        string
	recipients           []string
	ccRecipients         []string
	bccRecipients        []string
	parseMode            Mode
	domain               string
	headers              map[string]string
	tags                 []string
	setDKIM              bool
	enableNativeSend     bool
	requireTLS           bool
	skipVerification     bool
	enableTestMode       bool
	enableTracking       bool
	enableTrackingClicks bool
	enableTrackingOpens  bool
}

func (s *Service) applyOptions(opts ...Option) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, opt := range opts {
		opt(s)
	}
}

// New creates a new telegram service. It returns an error if the telegram client could not be created.
func New(domain, apiKey, senderAddress string, opts ...Option) (*Service, error) {
	client := mailgun.NewMailgun(domain, apiKey)

	s := &Service{
		client:        client,
		name:          "mailgun",
		logger:        nopadapter.NewAdapter(),
		renderMessage: defaultMessageRenderer,
		senderAddress: senderAddress,
	}

	s.applyOptions(opts...)

	return s, nil
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return s.name
}

// AddRecipients adds chat IDs that should receive messages.
func (s *Service) AddRecipients(recipients ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recipients = append(s.recipients, recipients...)
	s.logger.Debug().Int("count", len(recipients)).Int("total", len(s.recipients)).Msg("Recipients added")
}

// AddCCRecipients adds chat IDs that should receive messages.
func (s *Service) AddCCRecipients(recipients ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ccRecipients = append(s.ccRecipients, recipients...)
	s.logger.Debug().Int("count", len(recipients)).Int("total", len(s.ccRecipients)).Msg("CC Recipients added")
}

// AddBCCRecipients adds chat IDs that should receive messages.
func (s *Service) AddBCCRecipients(recipients ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bccRecipients = append(s.bccRecipients, recipients...)
	s.logger.Debug().Int("count", len(recipients)).Int("total", len(s.bccRecipients)).Msg("BCC Recipients added")
}
