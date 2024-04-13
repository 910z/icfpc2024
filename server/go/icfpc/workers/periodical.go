package workers

import (
	"context"
	"log/slog"
	"time"
)

func runPeriodical(ctx context.Context, interval time.Duration, wakeChannel chan struct{}, f func() error) error {
	if err := f(); err != nil {
		return err
	}
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			slog.DebugContext(ctx, "woken by timer")
			if err := f(); err != nil {
				return err
			}
		case <-wakeChannel:
			slog.DebugContext(ctx, "woken by channel")
			if err := f(); err != nil {
				return err
			}
		}
	}
}
