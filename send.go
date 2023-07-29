package notify

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Send sends a notification with the given subject and message through all the services of n. It performs these
// operations concurrently and returns the first encountered error, if any.
func (d *Dispatcher) Send(ctx context.Context, subject, message string, opts ...SendOption) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	eg, ctx := errgroup.WithContext(ctx)
	for _, service := range d.services {
		service := service

		eg.Go(func() error {
			if err := service.Send(ctx, subject, message, opts...); err != nil {
				return &ServiceFailureError{Service: service.Name(), Cause: err}
			}

			return nil
		})
	}

	return eg.Wait()
}
