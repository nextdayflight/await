package main

import (
	"context"
	"errors"
	"time"
)

const retryDelay = 500 * time.Millisecond

type timeoutError struct {
	Reason error
}

// Error implements the error interface.
func (e *timeoutError) Error() string {
	return e.Reason.Error()
}

type awaiter struct {
	logger  *LevelLogger
	timeout time.Duration
}

func (a *awaiter) run(resources []resource) error {
	if a.logger == nil {
		a.logger = NewLogger(errorLevel)
	}

	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	var latestErr error

	go func() {
		for _, res := range resources {
			a.logger.Infof("Awaiting resource: %s", res)

			for {
				select {
				case <-ctx.Done():
					// Exceeded timeout
					return
				default:
					// Still time left, let's continue
				}

				if latestErr = res.Await(ctx); latestErr != nil {
					if e, ok := latestErr.(*unavailabilityError); ok {
						// transient error
						a.logger.Debugf("Resource unavailable: %v", e)
					} else {
						// Maybe transient error
						a.logger.Errorf("Error: failed to await resource: %v", latestErr)
					}
					time.Sleep(retryDelay)
				} else {
					a.logger.Infof("Resource found: %s", res)
					// Resource found, move on to next one
					break
				}
			}
		}

		cancel() // All resources are available
	}()

	<-ctx.Done()
	switch ctx.Err() {
	case context.Canceled:
		if latestErr != nil {
			return latestErr
		}
		// All resources are available
		return nil
	case context.DeadlineExceeded:
		if latestErr == nil {
			// Time out even before the first try
			latestErr = errors.New("initial await did not finish")
		}
		return &unavailabilityError{latestErr}
	default:
		return errors.New("unknown error")
	}
}
