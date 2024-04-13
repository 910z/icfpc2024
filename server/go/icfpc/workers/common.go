package workers

import (
	"context"
	"time"
)

func runPeriodical(ctx context.Context, interval time.Duration, f func() error) error {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := f(); err != nil {
				return err
			}
		}
	}
}
