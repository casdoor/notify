// Package notify provides an abstraction for sending notifications
// through multiple services. It allows easy addition of new services
// and uniform handling of notification sending.
package notify

import (
	"context"
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
