package runner

import (
	"context"
	"errors"
	"fmt"
	"icfpc/database"
	"icfpc/integration"
	"log/slog"
	"time"

	"github.com/uptrace/bun"
)

func isRetriable(err error) bool {
	return errors.Is(err, integration.Error{})
}

func RunTasksFetcher(
	ctx context.Context,
	getTasks func() ([]database.Task, error),
	db *bun.DB,
) error {
	handleTasks := func() error {
		tasks, err := getTasks()
		if err != nil {
			slog.ErrorContext(ctx, "failed to fetch tasks", slog.Any("error", err))
			return err // тут не валимся, потому что ошибку от апи icfpc ретраим
		}

		err = saveTasks(ctx, db, tasks)
		if err != nil {
			slog.ErrorContext(ctx, "failed to save tasks", slog.Any("error", err))
			return err
		}
		return nil
	}

	err := handleTasks()
	if err != nil && !isRetriable(err) {
		return err
	}

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := handleTasks()
			if err != nil && !isRetriable(err) {
				return err
			}
		}
	}
}

func saveTasks(ctx context.Context, db *bun.DB, tasks []database.Task) error {
	res, err := db.NewInsert().
		Model(&tasks).
		On("CONFLICT DO NOTHING"). // когда совпали контент и айдишник. если контент поменялся, добавит дубль
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to save tasks: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to save tasks: %w", err)
	}
	slog.InfoContext(ctx, "saved tasks", slog.Int64("affected", affected))
	return nil
}
