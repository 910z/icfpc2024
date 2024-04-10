package workers

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"icfpc/database"
	"icfpc/evaluation"
	"icfpc/logs"

	"github.com/uptrace/bun"
)

type evaluator struct {
	db *bun.DB
}

func NewSolutionEvaluator(db *bun.DB) *evaluator {
	return &evaluator{
		db: db,
	}
}

func (e evaluator) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := e.evalEverythingPresent(ctx); err != nil {
				return err
			}
		}
	}
}

func (e evaluator) evalEverythingPresent(ctx context.Context) error {
	var runResults []database.RunResult
	if err := e.db.NewSelect().
		Model(&runResults).
		Join("FULL OUTER JOIN run_eval_results ON run_result.id = run_eval_results.run_result_id").
		Where("run_eval_results.run_result_id is null").
		Scan(ctx); err != nil {

		return fmt.Errorf("failed to get task results: %w", err)
	}

	for _, runResult := range runResults {
		runCtx := logs.WithRunResultLogging(ctx, runResult)
		runEvalRes := database.RunEvalResult{
			RunResultID: runResult.ID,
			EvalResult:  database.EvalResult{},
			StartedAt:   time.Now().UTC(),
			Status:      database.ProgressStatusStarted,
			Version:     evaluation.Version,
		}
		err := e.db.NewInsert().Model(&runEvalRes).
			Returning("id").
			Ignore().
			Scan(ctx, &runEvalRes.ID)
		if err == sql.ErrNoRows { // такой прогон уже был запущен
			continue
		}
		if err != nil {
			return err
		}

		go e.runEval(runCtx, &runResult, runEvalRes)
	}
	return nil
}

func (e evaluator) runEval(ctx context.Context, runRes *database.RunResult, runEvalRes database.RunEvalResult) {
	slog.InfoContext(ctx, "evaluation started")
	defer func() {
		if err := recover(); err != nil {
			slog.ErrorContext(ctx, "recovered panic in runEval", slog.Any("error", err))
			panic(err)
		} else {
			slog.InfoContext(ctx, "evaluation finished")
		}
	}()
	updateQ := e.db.NewUpdate().Model(&runEvalRes).WherePK()

	handleError := func(err error) {
		slog.ErrorContext(ctx, "evaluation failed", slog.Any("error", err))
		runEvalRes.Status = database.ProgressStatusError
		runEvalRes.FinishedAt = time.Now().UTC()
		runEvalRes.Error = err.Error()
		if err = database.UpdateEnsured(ctx, updateQ); err != nil {
			panic(err)
		}
	}

	task := database.Task{ID: runRes.TaskID}

	if err := e.db.NewSelect().Model(&task).WherePK().Scan(ctx); err != nil {
		handleError(err)
		return
	}

	evalRes, err := evaluation.EvaluateSolution(ctx, task, runRes.Solution)
	if err != nil {
		handleError(err)
		return
	}

	runEvalRes.EvalResult = evalRes
	runEvalRes.FinishedAt = time.Now().UTC()
	runEvalRes.Status = database.ProgressStatusFinished
	if err = database.UpdateEnsured(ctx, updateQ); err != nil {
		handleError(err)
		return
	}
}
