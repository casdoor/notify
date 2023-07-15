package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/nikoksr/notify/v2"
)

var (
	_ notify.Service = (*Service)(nil)
	_ client         = (*authClient)(nil)
	_ client         = (*webhookClient)(nil)
)

func defaultMessageRenderer(conf SendConfig) string {
	var builder strings.Builder

	builder.WriteString(conf.subject)
	builder.WriteString("\n\n")
	builder.WriteString(conf.message)

	return builder.String()
}

type authClient struct {
	session *discordgo.Session
}

type webhookClient struct {
	session *discordgo.Session
}

type client interface {
	setSession(session *discordgo.Session)
	sendTo(recipient string, conf SendConfig) error
}

func (c *authClient) setSession(session *discordgo.Session) {
	c.session = session
}

func attachmentsToFiles(attachments []notify.Attachment) []*discordgo.File {
	var files []*discordgo.File
	for _, attachment := range attachments {
		files = append(files, &discordgo.File{
			Reader: attachment,
			Name:   attachment.Name(),
		})
	}

	return files
}

func (c *webhookClient) setSession(session *discordgo.Session) {
	c.session = session
}

// Service struct holds necessary data to communicate with the Discord API.
type Service struct {
	client        client
	recipients    []string
	name          string
	renderMessage func(conf SendConfig) string
}

func newService(client client, name string, opts ...Option) (*Service, error) {
	svc := &Service{
		client:        client,
		name:          name,
		renderMessage: defaultMessageRenderer,
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc, nil
}

// New creates a new Discord service using an OAuth2 token for authentication.
func New(token string, opts ...Option) (*Service, error) {
	session, err := authenticateWithOAuth2Token(token)
	if err != nil {
		return nil, err
	}

	client := &authClient{session: session}

	return newService(client, "discord", opts...)
}

// NewBot creates a new Discord bot service.
func NewBot(token string, opts ...Option) (*Service, error) {
	session, err := authenticateWithBotToken(token)
	if err != nil {
		return nil, err
	}

	client := &authClient{session: session}

	return newService(client, "discord-bot", opts...)
}

// NewWebhook creates a new Discord webhook service. The recipient string must be a webhook URL.
func NewWebhook(opts ...Option) (*Service, error) {
	session, err := authenticate("") // Create an unauthenticated session.
	if err != nil {
		return nil, err
	}

	client := &webhookClient{session: session}

	return newService(client, "discord-webhook", opts...)
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return s.name
}

// AddRecipients takes Service channel IDs or webhook URLs and adds them to the list of recipients. You can add more
// channel IDs or webhook URLs by calling AddRecipients again.
func (s *Service) AddRecipients(recipients ...string) {
	s.recipients = append(s.recipients, recipients...)
}

// Option is a function that applies an option to the service.
type Option = func(*Service)

// WithClient sets the discord client (session) to use for sending messages. Naming it WithClient to stay consistent
// with the other services.
func WithClient(session *discordgo.Session) Option {
	return func(s *Service) {
		s.client.setSession(session)
	}
}

// WithRecipients sets the channel IDs or webhook URLs to send messages to. You can add more channel IDs or webhook URLs
// by calling Service.AddRecipients.
func WithRecipients(recipients ...string) Option {
	return func(d *Service) {
		d.recipients = recipients
	}
}

// WithName sets the name of the service. The default is "discord".
func WithName(name string) Option {
	return func(d *Service) {
		d.name = name
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
	return func(t *Service) {
		t.renderMessage = builder
	}
}
