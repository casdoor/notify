package notify

import (
	"context"
	"io"

	"golang.org/x/sync/errgroup"
)

// send calls the underlying notification services to send the given subject and message to their respective endpoints.
func (n *Notify) send(ctx context.Context, subject, message string, opts ...SendOption) error {
	var eg errgroup.Group
	for _, service := range n.services {
		service := service
		eg.Go(func() error {
			return newErrSendNotification(service.Name(), service.Send(ctx, subject, message, opts...))
		})
	}

	return eg.Wait()
}

// Send calls the underlying notification services to send the given subject and message to their respective endpoints.
func (n *Notify) Send(ctx context.Context, subject, message string, opts ...SendOption) error {
	return n.send(ctx, subject, message, opts...)
}

// Send calls the underlying notification services to send the given subject and message to their respective endpoints.
func Send(ctx context.Context, subject, message string, opts ...SendOption) error {
	return defaultNotify.Send(ctx, subject, message, opts...)
}

// SendConfig is an interface that can be used to configure a Send call. It is intended to be implemented by the various
// notification services and their respective SendConfig implementations.
type SendConfig interface {
	SetAttachments(attachments ...Attachment)
	SetMetadata(metadata map[string]any)
}

// SendOption is a function that configures a Send call. More specifically, it configures the SendConfig instance that
// is passed to the Send call.
type SendOption = func(SendConfig)

// Attachment is a type that can be used to attach a file to a notification message. It implements the io.Reader interface,
// so it can be used to attach a file to a notification message. Name is used as the filename when the attachment is sent
// to the notification service. Size is the size of the attachment in bytes, or -1 if the size is unknown.
type Attachment interface {
	io.Reader
	Name() string
}

// SendWithAttachments is a SendOption that attaches the given attachments to the SendConfig instance. This is a so-
// called shared SendOption, meaning that it is intended to be used by multiple notification services and is not specific
// to a single service.
// The Attachment type implements the io.Reader interface, so it can be used to attach a file to a notification message.
// The name of the attachment is used as the filename when the attachment is sent to the notification service.
func SendWithAttachments(attachments ...Attachment) SendOption {
	return func(c SendConfig) {
		c.SetAttachments(attachments...)
	}
}

// SendWithMetadata is a SendOption that attaches the given metadata to the SendConfig instance. This is a so-called
// shared SendOption, meaning that it is intended to be used by multiple notification services and is not specific to a
// single service.
func SendWithMetadata(metadata map[string]any) SendOption {
	return func(c SendConfig) {
		c.SetMetadata(metadata)
	}
}
