package notify

import "context"

// DefaultDispatcher is the default Dispatcher instance.
var DefaultDispatcher = New()

// UseServices appends the given service(s) to the defaultNotify instance's services list. Nil services are ignored.
func UseServices(services ...Service) {
	DefaultDispatcher.UseServices(services...)
}

// Send sends a notification with the given subject and message through all the services of the defaultNotify instance.
// It performs these operations concurrently and returns the first encountered error, if any.
func Send(ctx context.Context, subject, message string, opts ...SendOption) error {
	return DefaultDispatcher.Send(ctx, subject, message, opts...)
}
