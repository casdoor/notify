package twilio

import (
	"strings"
	"sync"

	"github.com/nikoksr/onelog"
	nopadapter "github.com/nikoksr/onelog/adapter/nop"
	"github.com/twilio/twilio-go"

	"github.com/nikoksr/notify/v2"
)

var _ notify.Service = (*Service)(nil)

func defaultMessageRenderer(conf *SendConfig) string {
	var builder strings.Builder

	builder.WriteString(conf.Subject)
	builder.WriteString("\n\n")
	builder.WriteString(conf.Message)

	return builder.String()
}

// Service is the twilio service. It is used to send messages to Twilio chats.
type Service struct {
	client *twilio.RestClient

	name          string
	mu            sync.RWMutex
	logger        onelog.Logger
	renderMessage func(conf *SendConfig) string
	dryRun        bool
	continueOnErr bool

	// Twilio specific
	senderPhoneNumber string
	phoneNumbers      []string
}

func (s *Service) applyOptions(opts ...Option) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, opt := range opts {
		opt(s)
	}
}

// Common function to create a new service
func newService(username, password, accountSid, phoneNumber string, opts ...Option) (*Service, error) {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   username,
		Password:   password,
		AccountSid: accountSid,
	})

	s := &Service{
		client:            client,
		name:              "twilio",
		logger:            nopadapter.NewAdapter(),
		renderMessage:     defaultMessageRenderer,
		dryRun:            false,
		continueOnErr:     false,
		senderPhoneNumber: phoneNumber,
	}

	s.applyOptions(opts...)

	return s, nil
}

// New creates a new twilio service. It returns an error if the twilio client could not be created.
func New(accountSID, apiKey, apiSecret, phoneNumber string, opts ...Option) (*Service, error) {
	return newService(apiKey, apiSecret, accountSID, phoneNumber, opts...)
}

func NewWithCredentials(username, token, phoneNumber string, opts ...Option) (*Service, error) {
	return newService(username, token, "", phoneNumber, opts...)
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return s.name
}

// AddRecipients adds phonenumbers that should receive messages.
func (s *Service) AddRecipients(phoneNumbers ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.phoneNumbers = append(s.phoneNumbers, phoneNumbers...)
	s.logger.Info().Int("count", len(phoneNumbers)).Int("total", len(s.phoneNumbers)).Msg("Recipients added")
}
