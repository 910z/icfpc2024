package workers

import (
	"context"
	"errors"
	"fmt"
	"icfpc/database"
	"icfpc/integration"
	"icfpc/logs"
	"log/slog"
	"reflect"
	"time"

	"github.com/uptrace/bun"
)

type tasksFetcher struct {
	db  *bun.DB
	bus bus
}

func NewTasksFetcher(db *bun.DB, bus bus) *tasksFetcher {
	return &tasksFetcher{db: db, bus: bus}
}

func (t tasksFetcher) Run(
	ctx context.Context,
	getTasks func() ([]database.Task, error),
) error {
	fetchTasks := func() error {
		tasks, err := getTasks()

		if errors.Is(err, integration.Error) {
			slog.WarnContext(ctx, "can't save tasks, retrying", slog.Any("error", err))
			return nil
		}

		if err != nil {
			slog.ErrorContext(ctx, "failed to fetch tasks", slog.Any("error", err))
			return err
		}

		err = t.saveTasks(ctx, tasks)
		if err != nil {
			slog.ErrorContext(ctx, "failed to save tasks", slog.Any("error", err))
			return err
		}

		return nil
	}
	return runPeriodical(logs.WithType(ctx, reflect.TypeOf(t)), time.Second, make(chan struct{}), fetchTasks)
}

func (t tasksFetcher) saveTasks(ctx context.Context, tasks []database.Task) error {
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
	if affected != 0 {
		slog.InfoContext(ctx, "saved tasks", slog.Int64("affected", affected))
		t.bus.onTasksAdded(tasks)
	}
	return nil
}
