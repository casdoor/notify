package notify

var _ SendConfigurer = (*SendConfig)(nil)

// SendConfigurer is used to configure the Send call.
type SendConfigurer interface {
	// SetAttachments sets attachments that can be sent alongside the message.
	SetAttachments(attachments ...Attachment)

	// SetMetadata sets additional metadata that can be sent with the message.
	SetMetadata(metadata map[string]any)

	// SetDryRun sets the dry run flag.
	SetDryRun(dryRun bool)

	// SetContinueOnErr sets the continue on error flag.
	SetContinueOnErr(continueOnError bool)
}

// SendConfig is the shared configuration for all senders.
type SendConfig struct {
	// Subject is the subject of the message. Subject will be used directly by services that support a distinct subject
	// field. For services that do not support a distinct subject field, this field will be ignored, unless used otherwise
	// by the service's message renderer function.
	Subject string

	// Message is the body of the message. This is the central part of the message and will be used by all services.
	Message string

	// Attachments are optional files that can be sent alongside the message.
	Attachments []Attachment

	// Metadata is optional additional metadata that can be sent with the message. Commonly used to provide additional
	// context in the service's message renderer function.
	Metadata map[string]any

	// DryRun is a flag that, if set to true, will prevent the service from trying to send the message.
	DryRun bool

	// ContinueOnErr is a flag that, if set to true, will continue sending the message to the next recipient even if an
	// error occurred.
	ContinueOnErr bool
}

// SetAttachments sets attachments that can be sent alongside the message.
func (c *SendConfig) SetAttachments(attachments ...Attachment) {
	c.Attachments = attachments
}

// SetMetadata sets additional metadata that can be sent with the message.
func (c *SendConfig) SetMetadata(metadata map[string]any) {
	c.Metadata = metadata
}

// SetDryRun sets the dry run flag.
func (c *SendConfig) SetDryRun(dryRun bool) {
	c.DryRun = dryRun
}

// SetContinueOnErr sets the continue on error flag.
func (c *SendConfig) SetContinueOnErr(continueOnError bool) {
	c.ContinueOnErr = continueOnError
}

// SendOption is a function that modifies the configuration of a Send call.
type SendOption = func(SendConfigurer)

// SendWithAttachments attaches the provided files to the message being sent.
func SendWithAttachments(attachments ...Attachment) SendOption {
	return func(c SendConfigurer) {
		c.SetAttachments(attachments...)
	}
}

// SendWithMetadata attaches the provided metadata to the message being sent.
func SendWithMetadata(metadata map[string]any) SendOption {
	return func(c SendConfigurer) {
		c.SetMetadata(metadata)
	}
}

// SendWithDryRun sets the dry run flag. If set to true, the service will not try to authenticate or send the message.
func SendWithDryRun(dryRun bool) SendOption {
	return func(c SendConfigurer) {
		c.SetDryRun(dryRun)
	}
}

// SendWithContinueOnErr sets the continue on error flag. If set to true, the service will continue sending the
// message to the next recipient even if an error occurred.
func SendWithContinueOnErr(continueOnErr bool) SendOption {
	return func(c SendConfigurer) {
		c.SetContinueOnErr(continueOnErr)
	}
}
