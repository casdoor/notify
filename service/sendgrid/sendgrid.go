package sendgrid

import (
	"sync"

	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/nikoksr/onelog"
	nopadapter "github.com/nikoksr/onelog/adapter/nop"
	"github.com/sendgrid/sendgrid-go"

	"github.com/nikoksr/notify/v2"
)

var _ notify.Service = (*Service)(nil)

// Mode is the mode for sending messages. It is a string to allow for custom parsing modes. The default mode is
// ModeHTML.
type Mode string

const (
	// ModeText is treating messages as plain text.
	ModeText Mode = "text/plain"
	// ModeHTML is formatting messages with HTML. It is the default mode.
	ModeHTML Mode = "text/html"
)

func defaultMessageRenderer(conf *SendConfig) string {
	return conf.Message
}

// Service is the sendgrid service. It is used to send messages to Sendgrid chats.
type Service struct {
	client *sendgrid.Client

	name          string
	mu            sync.RWMutex
	logger        onelog.Logger
	renderMessage func(conf *SendConfig) string
	dryRun        bool
	continueOnErr bool // no-op for sendgrid

	// Sendgrid specific

	senderAddress    string
	senderName       string
	recipients       []string
	ccRecipients     []string
	bccRecipients    []string
	parseMode        Mode
	headers          map[string]string
	customArgs       map[string]string
	batchID          string
	ipPoolID         string
	asm              *mail.Asm
	mailSettings     *mail.MailSettings
	trackingSettings *mail.TrackingSettings
}

func (s *Service) applyOptions(opts ...Option) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, opt := range opts {
		opt(s)
	}
}

// New creates a new sendgrid service. It returns an error if the sendgrid client could not be created.
func New(apiKey, senderAddress, senderName string, opts ...Option) (*Service, error) {
	client := sendgrid.NewSendClient(apiKey)

	s := &Service{
		client:        client,
		name:          "sendgrid",
		logger:        nopadapter.NewAdapter(),
		renderMessage: defaultMessageRenderer,
		senderAddress: senderAddress,
		senderName:    senderName,
		parseMode:     ModeHTML,
	}

	s.applyOptions(opts...)

	return s, nil
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return s.name
}

// AddRecipients adds email addresses that should receive messages.
func (s *Service) AddRecipients(addresses ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recipients = append(s.recipients, addresses...)
	s.logger.Debug().Int("count", len(addresses)).Int("total", len(s.recipients)).Msg("Recipients added")
}
