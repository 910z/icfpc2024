package workers

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

type TasksFetcher struct {
	db *bun.DB
}

func NewTasksFetcher(db *bun.DB) *TasksFetcher {
	return &TasksFetcher{db: db}
}

func (t *TasksFetcher) isRetriable(err error) bool {
	return errors.Is(err, integration.Error)
}

func (t *TasksFetcher) Run(
	ctx context.Context,
	getTasks func() ([]database.Task, error),
) error {
	handleTasks := func() error {
		tasks, err := getTasks()
		if err != nil {
			slog.ErrorContext(ctx, "failed to fetch tasks", slog.Any("error", err))
			return err // тут не валимся, потому что ошибку от апи icfpc ретраим
		}

		err = t.saveTasks(ctx, tasks)
		if err != nil {
			slog.ErrorContext(ctx, "failed to save tasks", slog.Any("error", err))
			return err
		}
		return nil
	}

	err := handleTasks()
	if err != nil && !t.isRetriable(err) {
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
			if err != nil && !t.isRetriable(err) {
				return err
			}
		}
	}
}

func (t *TasksFetcher) saveTasks(ctx context.Context, tasks []database.Task) error {
	res, err := t.db.NewInsert().
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
