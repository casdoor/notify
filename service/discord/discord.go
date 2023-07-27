package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nikoksr/onelog"
	nopadapter "github.com/nikoksr/onelog/adapter/nop"

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
	logger  onelog.Logger
}

type webhookClient struct {
	session *discordgo.Session
	logger  onelog.Logger
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
	client client

	logger        onelog.Logger
	recipients    []string
	name          string
	renderMessage func(conf SendConfig) string
}

func newService(client client, name string, opts ...Option) (*Service, error) {
	svc := &Service{
		client:        client,
		logger:        nopadapter.NewAdapter(),
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
	s, err := newService(nil, "discord", opts...)
	if err != nil {
		return nil, err
	}

	session, err := authenticateWithOAuth2Token(token)
	if err != nil {
		return nil, err
	}

	s.client = &authClient{
		session: session,
		logger:  s.logger,
	}

	return s, nil
}

// NewBot creates a new Discord bot service.
func NewBot(token string, opts ...Option) (*Service, error) {
	s, err := newService(nil, "discord-bot", opts...)
	if err != nil {
		return nil, err
	}

	session, err := authenticateWithBotToken(token)
	if err != nil {
		return nil, err
	}

	s.client = &authClient{
		session: session,
		logger:  s.logger,
	}

	return s, nil
}

// NewWebhook creates a new Discord webhook service. The recipient string must be a webhook URL.
func NewWebhook(opts ...Option) (*Service, error) {
	s, err := newService(nil, "discord-webhook", opts...)
	if err != nil {
		return nil, err
	}

	session, err := authenticate("") // Create an unauthenticated session.
	if err != nil {
		return nil, err
	}

	s.client = &webhookClient{
		session: session,
		logger:  s.logger,
	}

	return s, nil
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return s.name
}

// AddRecipients takes Service channel IDs or webhook URLs and adds them to the list of recipients. You can add more
// channel IDs or webhook URLs by calling AddRecipients again.
func (s *Service) AddRecipients(recipients ...string) {
	s.recipients = append(s.recipients, recipients...)
	s.logger.Info().Int("count", len(recipients)).Int("total", len(s.recipients)).Msg("Recipients added")
}
