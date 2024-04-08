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
	db *bun.DB
}

type plannedRun struct {
	task      *database.Task
	algorithm algorithms.IAlgorithm
}

/*
SELECT run_results.*
FROM run_results
RIGHT JOIN (
    SELECT *
    FROM tasks
    CROSS JOIN (
        VALUES
        ('*algorithms.Tripler', '1.0.0'),
        ('*algorithms.Tripler', '1.0.0')
    ) AS algorithms(algorithm_name, algorithm_version)
) AS ttt_algorithms ON
    run_results.algorithm_name = ttt_algorithms.algorithm_name
    AND run_results.algorithm_version = ttt_algorithms.algorithm_version
    AND run_results.task_id = ttt_algorithms.id
WHERE ttt_algorithms.algorithm_name IS NULL
    AND ttt_algorithms.algorithm_version IS NULL
    AND ttt_algorithms.id IS NULL;
*/

type algorithmData struct {
	Name    string
	Version string
}

func toAlgorithmDatas(algs []algorithms.IAlgorithm) []algorithmData {
	var result []algorithmData
	for _, algorithm := range algs {
		result = append(result, algorithmData{
			Name:    algorithms.GetName(algorithm),
			Version: algorithm.Version(),
		})
	}
	return result
}

func (r Runner) planRuns(
	ctx context.Context,
	currentAlgorithms []algorithms.IAlgorithm,
) ([]database.RunResult, error) {
	algDatas := toAlgorithmDatas(currentAlgorithms)
	algValues := r.db.NewValues(&algDatas)
	var runResults []database.RunResult
	err := r.db.NewSelect().
		With("algorithms", algValues).
		Model(&runResults).
		Join(`FULL OUTER JOIN (select * from tasks CROSS JOIN algorithms) as task_algo
			ON 
				task_algo.name = run_result.algorithm_name 
				AND task_algo.version = run_result.algorithm_version
				AND task_algo.id = run_result.task_id`).
		Where("run_result.id IS NULL").
		Column("task_id", "algorithm_name").
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return runResults, nil
}

func (r Runner) RunAlgorithms(ctx context.Context, algorithms []algorithms.IAlgorithm) error {
	for {
		runs, err := r.planRuns(ctx, algorithms)
		if err != nil {
			return err
		}
		slog.InfoContext(ctx, "runs planned", slog.Any("runs", runs))

		for _, task := range AllTasks {
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
	ctx context.Context,
	algorithm algorithms.IAlgorithm,
	task database.Task,
) (
	_ database.Solution,
	_ database.SolutionExplanation,
	err error,
) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered in algorithm.Solve: %v", r)
		}
	}()
	solution, explanation, err := algorithm.Solve(ctx, task)
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
		runResult.Error = err.Error()

		if _, err = r.db.NewUpdate().Model(&runResult).WherePK().Exec(ctx); err != nil {
			panic(err)
		}
	}

	solution, explanation, err := safeSolve(ctx, algorithm, task)
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
