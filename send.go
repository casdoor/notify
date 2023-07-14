package notify

import (
	"context"
	"io"

	"golang.org/x/sync/errgroup"
)

// Send sends a notification with the given subject and message through all the services of n. It performs these
// operations concurrently and returns the first encountered error, if any.
func (n *Notify) Send(ctx context.Context, subject, message string, opts ...SendOption) error {
	var eg errgroup.Group
	for _, service := range n.services {
		service := service

		eg.Go(func() error {
			if err := service.Send(ctx, subject, message, opts...); err != nil {
				return &ErrServiceFailure{Service: service.Name(), Cause: err}
			}

			return nil
		})
	}

	return eg.Wait()
}

// Send sends a notification with the given subject and message through all the services of the defaultNotify instance.
// It performs these operations concurrently and returns the first encountered error, if any.
func Send(ctx context.Context, subject, message string, opts ...SendOption) error {
	return defaultNotify.Send(ctx, subject, message, opts...)
}

// SendConfig is used to configure the Send call.
type SendConfig interface {
	// SetAttachments sets attachments that can be sent alongside the message.
	SetAttachments(attachments ...Attachment)
	// SetMetadata sets additional metadata that can be sent with the message.
	SetMetadata(metadata map[string]any)
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
