package workers

import (
	"context"
	"database/sql"
	"fmt"
	"icfpc/algorithms"
	"icfpc/database"
	"icfpc/logs"
	"log/slog"
	"time"

	"github.com/uptrace/bun"
)

func NewAlgorithmRunner(db *bun.DB) *Runner {
	return &Runner{
		db: db,
	}
}

type Runner struct {
	db *bun.DB
}

type algorithmData struct {
	AlgorithmName    string
	AlgorithmVersion string
}

func toAlgorithmDatas(algs []algorithms.IAlgorithm) []algorithmData {
	var result []algorithmData
	for _, algorithm := range algs {
		result = append(result, algorithmData{algorithms.GetName(algorithm), algorithm.Version()})
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
		Column("task_algo.task_id", "task_algo.algorithm_name", "task_algo.algorithm_version").
		With("algorithms", algValues).
		Model(&runResults).
		Join(`FULL OUTER JOIN (SELECT id AS task_id, algorithms.* FROM tasks CROSS JOIN algorithms) as task_algo
			ON 
				task_algo.algorithm_name = run_result.algorithm_name 
				AND task_algo.algorithm_version = run_result.algorithm_version
				AND task_algo.task_id = run_result.task_id`).
		Where("run_result.id IS NULL").
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return runResults, nil
}

func getTasksToAdd(taskCache map[int64]database.Task, plannedRuns []database.RunResult) []int64 {
	uniqueTaskIds := make(map[int64]struct{})
	for _, run := range plannedRuns {
		uniqueTaskIds[run.TaskID] = struct{}{}
	}
	keys := make([]int64, 0, len(uniqueTaskIds))
	for k := range uniqueTaskIds {
		if _, ok := taskCache[k]; !ok {
			keys = append(keys, k)
		}
	}
	return keys
}

// может, брать его сразу из памяти, когда из апишки icfpc читаем
func addToTaskCache(ctx context.Context, db *bun.DB, taskIds []int64, taskCache map[int64]database.Task) error {
	tasks := make([]database.Task, 0, len(taskIds))
	err := db.NewSelect().Model(&tasks).Where("id IN (?)", bun.In(taskIds)).Scan(ctx)
	if err != nil {
		return err
	}
	for _, task := range tasks {
		taskCache[task.ID] = task
	}
	return nil
}

func getAlgorithm(algs []algorithms.IAlgorithm, name string, version string) algorithms.IAlgorithm {
	for _, algorithm := range algs {
		if algorithms.GetName(algorithm) == name && algorithm.Version() == version {
			return algorithm
		}
	}
	return nil
}

func (r Runner) Run(ctx context.Context, algs []algorithms.IAlgorithm) error {
	taskCache := make(map[int64]database.Task)

	for {
		runs, err := r.planRuns(ctx, algs)
		if err != nil {
			return err
		}

		slog.InfoContext(ctx, "runs planned", slog.Any("runs", runs))
		addToTaskCache(ctx, r.db, getTasksToAdd(taskCache, runs), taskCache)
		for _, plannedRun := range runs {
			task := taskCache[plannedRun.TaskID]
			algorithm := getAlgorithm(algs, plannedRun.AlgorithmName, plannedRun.AlgorithmVersion)
			if algorithm == nil {
				return fmt.Errorf("algorithm not found: %s %s", plannedRun.AlgorithmName, plannedRun.AlgorithmVersion)
			}

			runResult := database.RunResult{
				TaskID:           task.ID,
				AlgorithmName:    algorithms.GetName(algorithm),
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
