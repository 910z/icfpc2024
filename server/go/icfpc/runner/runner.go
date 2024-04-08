package runner

import (
	"context"
	"database/sql"
	"fmt"
	"icfpc/algorithms"
	"icfpc/database"
	"icfpc/logs"
	"log/slog"
	"reflect"
	"time"

	"github.com/uptrace/bun"
)

func New(db *bun.DB) *Runner {
	return &Runner{
		db: db,
	}
}

type Runner struct {
	db         *bun.DB
	tasks      []database.Task
	algorithms []algorithms.IAlgorithm
}

// TODO: this is bad
type explanationWithError struct {
	algorithms.Explanation

	Error string
}

func (r Runner) Run(ctx context.Context, tasks []database.Task, algorithms []algorithms.IAlgorithm) error {
	for {
		for _, task := range tasks {
			for _, algorithm := range algorithms {
				runResult := database.RunResult{
					TaskID:           task.ID,
					AlgorithmName:    reflect.TypeOf(algorithm).String(),
					AlgorithmVersion: algorithm.Version(),
					Status:           database.RunStatusStarted,
					StartedAt:        time.Now().UTC(),
				}

				err := r.db.NewInsert().Model(&runResult).
					Returning("id").
					Ignore().
					Scan(ctx, &runResult.ID)
				if err == sql.ErrNoRows { // такой прогон уже был запущен
					continue
				}
				if err != nil {
					return err
				}

				workerContext := logs.WithRunResultLogging(ctx, runResult)
				go r.runWorker(workerContext, task, algorithm, runResult)
			}
		}

		time.Sleep(time.Second)
	}
}

func safeSolve(
	algorithm algorithms.IAlgorithm,
	task database.Task,
) (
	_ algorithms.Solution,
	_ algorithms.Explanation,
	err error,
) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered in algorithm.Solve: %v", r)
		}
	}()
	solution, explanation, err := algorithm.Solve(task)
	return solution, explanation, err
}

func (r Runner) runWorker(
	ctx context.Context,
	task database.Task,
	algorithm algorithms.IAlgorithm,
	runResult database.RunResult,
) {
	defer func() {
		if err := recover(); err != nil {
			slog.ErrorContext(ctx, "recovered panic in runWorker", slog.Any("error", err))
			panic(err)
		} else {
			slog.InfoContext(ctx, "finished")
		}
	}()
	slog.InfoContext(ctx, "started")

	handleError := func(err error) {
		runResult.Status = database.RunStatusError
		runResult.FinishedAt = time.Now().UTC()
		runResult.Explanation = &explanationWithError{
			Error: err.Error(),
		}

		if _, err = r.db.NewUpdate().Model(&runResult).WherePK().Exec(ctx); err != nil {
			panic(err)
		}
	}

	solution, explanation, err := safeSolve(algorithm, task)
	if err != nil {
		handleError(err)

		return
	}

	runResult.FinishedAt = time.Now().UTC()
	runResult.Solution = solution
	runResult.Explanation = explanation
	runResult.Status = database.RunStatusFinished

	res, err := r.db.NewUpdate().Model(&runResult).WherePK().Exec(ctx)
	if err != nil {
		handleError(err)

		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		handleError(err)

		return
	}

	if rows == 0 {
		handleError(fmt.Errorf("no rows updated: %v", runResult))

		return
	}
}
