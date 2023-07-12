package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/nikoksr/notify/v2"
)

var _ notify.Service = (*Service)(nil)

const (
	defaultUserAgent     = "notify/" + notify.Version
	defaultContentType   = "application/json; charset=utf-8"
	defaultRequestMethod = http.MethodPost

	// Defining these as constants for testing purposes.
	defaultSubjectKey = "subject"
	defaultMessageKey = "message"
)

type (
	// PreSendHookFn defines a function signature for a pre-send hook.
	PreSendHookFn func(req *http.Request) error

	// PostSendHookFn defines a function signature for a post-send hook.
	PostSendHookFn func(req *http.Request, resp *http.Response) error

	// Webhook represents a single webhook recipient. It contains all the information needed to sendToWebhook a valid request to
	// the recipient. The buildPayload function is used to build the payload that will be sent to the recipient from the
	// given subject and message.
	Webhook struct {
		ContentType string
		Header      http.Header
		Method      string
		URL         string
	}

	// Service is the main struct of this package. It contains all the information needed to sendToWebhook notifications to a
	// list of recipients. The recipients are represented by Webhooks and are expected to be valid HTTP endpoints. The
	// Service also allows
	Service struct {
		client        *http.Client
		webhooks      []*Webhook
		name          string
		renderMessage func(conf SendConfig) string
		preSendHooks  []PreSendHookFn
		postSendHooks []PostSendHookFn
	}
)

func defaultMessageRenderer(conf SendConfig) string {
	payload := map[string]string{
		defaultSubjectKey: conf.subject,
		defaultMessageKey: conf.message,
	}

	// Marshal the payload to a JSON byte slice.
	payloadRaw, _ := json.Marshal(payload)

	return string(payloadRaw)
}

func newWebhook(url string) *Webhook {
	return &Webhook{
		ContentType: defaultContentType,
		Header:      http.Header{},
		Method:      defaultRequestMethod,
		URL:         url,
	}
}

// String returns a string representation of the webhook. It implements the fmt.Stringer interface.
func (w *Webhook) String() string {
	if w == nil {
		return ""
	}

	return strings.TrimSpace(fmt.Sprintf("%s %s %s", strings.ToUpper(w.Method), w.URL, w.ContentType))
}

// New returns a new instance of a Service notification service. Parameter 'tag' is used as a log prefix and may be left
// empty, it has a fallback value.
func New(opts ...Option) *Service {
	s := &Service{
		client:        http.DefaultClient,
		webhooks:      []*Webhook{},
		name:          "http",
		renderMessage: defaultMessageRenderer,
		preSendHooks:  []PreSendHookFn{},
		postSendHooks: []PostSendHookFn{},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Name returns the name of the service. It implements the notify.Service interface.
func (s *Service) Name() string {
	return s.name
}

// AddRecipients accepts a list of Webhooks and adds them as recipients. The Webhooks are expected to be valid HTTP
// endpoints.
func (s *Service) AddRecipients(webhooks ...*Webhook) {
	s.webhooks = append(s.webhooks, webhooks...)
}

// AddRecipientsURLs accepts a list of URLs and adds them as recipients. Internally it converts the URLs to Webhooks by
// using the default content-type ("application/json") and request method ("POST").
func (s *Service) AddRecipientsURLs(urls ...string) {
	for _, url := range urls {
		s.AddRecipients(newWebhook(url))
	}
}

// Option is a function that can be used to configure the http service.
type Option = func(*Service)

// WithClient sets the HTTP client. This is useful if you want to use a custom client.
func WithClient(client *http.Client) Option {
	return func(t *Service) {
		t.client = client
	}
}

func WithRecipients(webhooks ...*Webhook) Option {
	return func(s *Service) {
		s.AddRecipients(webhooks...)
	}
}

func WithRecipientsURLS(webhookURLs ...string) Option {
	return func(s *Service) {
		s.AddRecipientsURLs(webhookURLs...)
	}
}

// WithName sets the name of the service. The default name is "http".
func WithName(name string) Option {
	return func(s *Service) {
		s.name = name
	}
}

// WithMessageRenderer sets the message renderer. The default function will put the subject and message on separate lines.
//
// Example:
//
//	http.WithMessageRenderer(func(conf SendConfig) string {
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
