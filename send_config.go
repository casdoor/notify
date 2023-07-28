package notify

import "io"

// SendConfig is used to configure the Send call.
type SendConfig interface {
	// SetAttachments sets attachments that can be sent alongside the message.
	SetAttachments(attachments ...Attachment)

	// SetMetadata sets additional metadata that can be sent with the message.
	SetMetadata(metadata map[string]any)

	// SetDryRun sets the dry run flag.
	SetDryRun(dryRun bool)

	// SetContinueOnErr sets the continue on error flag.
	SetContinueOnErr(continueOnError bool)
}

// SendOption is a function that modifies the configuration of a Send call.
type SendOption = func(SendConfig)

// Attachment represents a file that can be attached to a notification message.
type Attachment interface {
	// Reader is used to read the contents of the attachment.
	io.Reader

	// Name is used as the filename when sending the attachment.
	Name() string
}

// SendWithAttachments attaches the provided files to the message being sent.
func SendWithAttachments(attachments ...Attachment) SendOption {
	return func(c SendConfig) {
		c.SetAttachments(attachments...)
	}
}

// SendWithMetadata attaches the provided metadata to the message being sent.
func SendWithMetadata(metadata map[string]any) SendOption {
	return func(c SendConfig) {
		c.SetMetadata(metadata)
	}
}

// SendWithDryRun sets the dry run flag. If set to true, the service will not try to authenticate or send the message.
func SendWithDryRun(dryRun bool) SendOption {
	return func(c SendConfig) {
		c.SetDryRun(dryRun)
	}
}

// SendWithContinueOnErr sets the continue on error flag. If set to true, the service will continue sending the
// message to the next recipient even if an error occurred.
func SendWithContinueOnErr(continueOnErr bool) SendOption {
	return func(c SendConfig) {
		c.SetContinueOnErr(continueOnErr)
	}
}
