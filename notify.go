// Package notify provides an abstraction for sending notifications
// through multiple services. It allows easy addition of new services
// and uniform handling of notification sending.
package notify

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

// Service describes a notification service that can send messages.
// Each service implementation should provide its own way of sending messages.
type Service interface {
	// Name should return a unique identifier for the service.
	Name() string

	// Send sends a message with a subject through this service.
	// Additional options can be provided to customize the sending process.
	// Returns an error if the sending process failed.
	Send(ctx context.Context, subject string, message string, opts ...SendOption) error
}

// Dispatcher handles notifications sending through multiple services.
type Dispatcher struct {
	mu       sync.RWMutex
	services []Service // services contains all services through which notifications can be send.
}

func (d *Dispatcher) applyOptions(opts ...Option) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, opt := range opts {
		opt(d)
	}
}

// New creates a Dispatcher instance with the specified options.
// If no options are provided, a Dispatcher instance is created with no services and default options.
func New(opts ...Option) *Dispatcher {
	n := &Dispatcher{}

	n.applyOptions(opts...)

	return n
}

// UseServices appends the given service(s) to the Dispatcher instance's services list. Nil services are ignored.
func (d *Dispatcher) UseServices(services ...Service) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, svc := range services {
		if svc != nil {
			d.services = append(d.services, svc)
		}
	}
}

// Option is a function that configures a Dispatcher instance.
type Option = func(*Dispatcher)

// WithServices adds the given services to a Dispatcher instance's services list.
// The services will be used in the order they are provided. Nil services are ignored.
func WithServices(services ...Service) Option {
	return func(d *Dispatcher) {
		d.services = services
	}
}

// Attachment represents a file that can be attached to a notification message.
type Attachment interface {
	// Reader returns a reader that can be used to read the attachment's data.
	Reader() io.Reader

	// Name is used as the filename when sending the attachment.
	Name() string

	// ContentType is used as the content type when sending the attachment. It's optional in most cases.
	ContentType() string

	// Size is used as the content length when sending the attachment. It's optional in most cases.
	Size() int64

	// Inline is used to determine if the attachment should be inlined or not. This is optional and not supported by
	// all services.
	Inline() bool
}

// AttachmentFromReader creates an attachment from the provided reader.
func AttachmentFromReader(reader io.Reader, name, contentType string, inline bool) (Attachment, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	size := int64(len(data))

	return &attachment{
		data:        data,
		name:        name,
		contentType: contentType,
		size:        size,
		inline:      inline,
	}, nil
}

// AttachmentFromFile creates an attachment from the provided file.
func AttachmentFromFile(file *os.File, contentType string, inlined bool) (Attachment, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	name := stat.Name()

	return AttachmentFromReader(file, name, contentType, inlined)
}

// AttachmentFromPath creates an attachment from the file at the provided path.
func AttachmentFromPath(path string, contentType string, inlined bool) (Attachment, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return AttachmentFromFile(file, contentType, inlined)
}

var _ Attachment = (*attachment)(nil)

type attachment struct {
	data        []byte
	name        string
	contentType string
	size        int64
	inline      bool
}

func (a *attachment) Read(p []byte) (n int, err error) {
	return bytes.NewReader(a.data).Read(p)
}

func (a *attachment) Reader() io.Reader {
	return bytes.NewReader(a.data)
}

func (a *attachment) Name() string {
	return a.name
}

func (a *attachment) ContentType() string {
	return a.contentType
}

func (a *attachment) Size() int64 {
	return a.size
}

func (a *attachment) Inline() bool {
	return a.inline
}
