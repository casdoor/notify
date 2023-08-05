package mail

import (
	"sync"

	"github.com/nikoksr/onelog"
	nopadapter "github.com/nikoksr/onelog/adapter/nop"
	mail "github.com/xhit/go-simple-mail/v2"

	"github.com/nikoksr/notify/v2"
)

// Compile time check to make sure the service implements the notify.Service interface.
var _ notify.Service = (*Service)(nil)

type (
	// Mode is the mode for sending messages. The default mode is ModeHTML.
	Mode int

	// Priority is the priority for sending messages. The default priority is PriorityNormal.
	Priority int

	// Service is a structure that contains data needed for interaction with Mail's APIs.
	Service struct {
		server *mail.SMTPServer
		client *mail.SMTPClient

		name          string
		mu            sync.RWMutex
		logger        onelog.Logger
		renderMessage func(conf *SendConfig) string
		dryRun        bool
		continueOnErr bool // no-op for Mail

		// Mail specific

		senderName    string
		recipients    []string
		ccRecipients  []string
		bccRecipients []string
		parseMode     Mode
		priority      Priority
	}
)

const (
	// ModeHTML sends the message as HTML. It is the default mode.
	ModeHTML = Mode(mail.TextHTML)
	// ModePlain sends the message as plain text.
	ModePlain = Mode(mail.TextPlain)
	// ModeCalendar sends the message as a calendar invite.
	ModeCalendar = Mode(mail.TextCalendar)

	// PriorityLow is the low priority for sending messages.
	PriorityLow = Priority(mail.PriorityLow)
	// PriorityNormal is the default priority for sending messages.
	PriorityNormal = -1
	// PriorityHigh is the high priority for sending messages.
	PriorityHigh = Priority(mail.PriorityHigh)
)

// defaultMessageRenderer is a helper function to format messages for Mail.
func defaultMessageRenderer(conf *SendConfig) string {
	return conf.Message
}

// newServer creates a new Mail server from the given parameters. It also uses the following default values:
//
// - Authentication: AuthAuto
// - Encryption:     None
// - ConnectTimeout: 10 seconds
// - SendTimeout:    10 seconds
// - Helo:           localhost
// - KeepAlive:      true
func newServer(host string, port int, username, password string) *mail.SMTPServer {
	server := mail.NewSMTPClient()

	server.Host = host
	server.Port = port
	server.Username = username
	server.Password = password

	server.KeepAlive = true

	return server
}

func (s *Service) applyOptions(opts ...Option) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, opt := range opts {
		opt(s)
	}
}

func New(host string, port int, username, password string, opts ...Option) (*Service, error) {
	// Create a new Mail server
	server := newServer(host, port, username, password)

	s := &Service{
		server:        server,
		name:          "mail",
		logger:        nopadapter.NewAdapter(),
		renderMessage: defaultMessageRenderer,
		senderName:    "From Notify <no-reply>",
		parseMode:     ModeHTML,
		priority:      PriorityNormal,
	}

	s.applyOptions(opts...)

	// Connect to the SMTP server and return the client
	client, err := server.Connect()
	if err != nil {
		return nil, err
	}

	// Set the authenticated client that will be used to send the notifications
	s.client = client

	return s, nil
}

// Name returns the name of the service, which identifies the type of service in use (in this case, Mail).
func (s *Service) Name() string {
	return s.name
}

// AddRecipients appends given channel IDs onto an internal list that Send uses to distribute the notifications.
func (s *Service) AddRecipients(recipients ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recipients = append(s.recipients, recipients...)
	s.logger.Debug().Int("count", len(recipients)).Int("total", len(s.recipients)).Msg("Recipients added")
}

// AddCCRecipients appends given channel IDs onto an internal list that Send uses to distribute the notifications.
func (s *Service) AddCCRecipients(recipients ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ccRecipients = append(s.ccRecipients, recipients...)
	s.logger.Debug().Int("count", len(recipients)).Int("total", len(s.ccRecipients)).Msg("CC Recipients added")
}

// AddBCCRecipients appends given channel IDs onto an internal list that Send uses to distribute the notifications.
func (s *Service) AddBCCRecipients(recipients ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bccRecipients = append(s.bccRecipients, recipients...)
	s.logger.Debug().Int("count", len(recipients)).Int("total", len(s.bccRecipients)).Msg("BCC Recipients added")
}
