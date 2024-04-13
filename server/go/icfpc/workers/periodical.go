package workers

import (
	"context"
	"log/slog"
	"time"
)

func runPeriodical(ctx context.Context, interval time.Duration, wakeChannel chan struct{}, f func() error) error {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := f(); err != nil {
				return err
			}
		case <-wakeChannel:
			slog.DebugContext(ctx, "woken up")
			if err := f(); err != nil {
				return err
			}
		}
	}
}
