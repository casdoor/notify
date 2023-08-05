package twilio

import "github.com/nikoksr/notify/v2"

var _ notify.SendConfig = (*SendConfig)(nil)

// SendConfig represents the configuration needed for sending a message.
//
// This struct complies with the notify.SendConfig interface and allows you to alter
// the behavior of the send function. This can be achieved by either passing send options
// to the send function or by manipulating the fields of this struct in your custom
// message renderer.
//
// All fields of this struct are exported to offer maximum flexibility to users.
// However, users must be aware that they are responsible for managing thread-safety
// and other similar concerns when manipulating these fields directly.
type SendConfig struct {
	Subject       string
	Message       string
	Attachments   []notify.Attachment
	Metadata      map[string]any
	DryRun        bool
	ContinueOnErr bool

	// Twilio specific
}

// SetAttachments adds attachments to the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetAttachments(attachments ...notify.Attachment) {
	c.Attachments = attachments
}

// SetMetadata sets the metadata of the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetMetadata(metadata map[string]any) {
	c.Metadata = metadata
}

// SetDryRun sets the dry run flag of the message. This method is needed as part of the notify.SendConfig interface.
func (c *SendConfig) SetDryRun(dryRun bool) {
	c.DryRun = dryRun
}

// SetContinueOnErr sets the continue on error flag of the message. This method is needed as part of the
// notify.SendConfig interface.
func (c *SendConfig) SetContinueOnErr(continueOnErr bool) {
	c.ContinueOnErr = continueOnErr
}

// newSendConfig creates a new send config with default values.
func (s *Service) newSendConfig(subject, message string, opts ...notify.SendOption) *SendConfig {
	conf := &SendConfig{
		Subject:       subject,
		Message:       message,
		DryRun:        s.dryRun,
		ContinueOnErr: s.continueOnErr,
	}

	for _, opt := range opts {
		opt(conf)
	}

	conf.Message = s.renderMessage(conf)

	return conf
}
