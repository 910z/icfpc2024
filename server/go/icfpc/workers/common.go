package workers

import (
	"context"
	"time"
)

func runPeriodical(ctx context.Context, interval time.Duration, f func(ctx context.Context) error) error {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := f(ctx); err != nil {
				return err
			}
		}
	}
}
