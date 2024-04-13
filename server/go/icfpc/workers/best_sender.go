package workers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"icfpc/database"
	"icfpc/evaluation"
	"icfpc/integration"
	"icfpc/logs"

	"github.com/uptrace/bun"
)

type bestSender struct {
	db  *bun.DB
	bus bus
}

func NewBestSender(db *bun.DB, bus bus) *bestSender {
	return &bestSender{
		db:  db,
		bus: bus,
	}
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

type fullResult struct {
	Solution         database.Solution
	RunResultId      int64
	TaskID           int64
	TaskExternalId   string
	AlgorithmVersion string
	AlgorithmName    string
	SubmissionStatus database.SubmissionStatus
	Score            database.Score
}

func (b bestSender) Run(
	ctx context.Context,
	ord SortOrder,
	send func(context.Context, string, database.Solution) (string, error),
) error {
	return runPeriodical(ctx, time.Second, b.bus.solutionEvaluated, func() error {
		var best []fullResult
		err := b.db.NewRaw(fmt.Sprintf(`with cte as (
			select
			  run_results.solution,
			  run_results.id AS run_result_id,
			  tasks.external_id AS task_external_id,
			  ROW_NUMBER() over w as row_num,
			  run_results.submission_status,
			  run_results.algorithm_version,
			  run_results.algorithm_name,
			  run_results.task_id,
			  case run_results.submission_status
				when 'checked' then run_results.external_score
				else run_eval_results.score
			  end as score
			from
			  run_eval_results
			  join run_results on run_results.id = run_eval_results.run_result_id
			  join tasks on tasks.id = run_results.task_id
			where
			  run_eval_results.version = ?
			  AND run_eval_results.status = ?
			window w as (
				partition by run_results.task_id
				order by score %s
			  )
		  )
		  select
		    solution,
		    run_result_id,
			algorithm_version,
			algorithm_name,
			task_id,
			task_external_id,
			submission_status,
			score
		  from cte where row_num = 1
		  `, ord), evaluation.Version, database.ProgressStatusFinished).Scan(ctx, &best)
		if err != nil {
			return err
		}
		for i := range best {
			if best[i].SubmissionStatus != database.SubmissionStatusNotSubmitted {
				// уже отправлялся
				continue
			}
			runCtx := logs.WithRunResultLogging(ctx, database.RunResult{
				TaskID:           best[i].TaskID,
				AlgorithmName:    best[i].AlgorithmName,
				AlgorithmVersion: best[i].AlgorithmVersion,
			})
			slog.InfoContext(runCtx, "sending best", slog.Any("score", best[i].Score))
			token, err := send(runCtx, best[i].TaskExternalId, best[i].Solution)
			if errors.Is(err, integration.Error) {
				slog.WarnContext(ctx, "can't send, will retry next time", slog.Any("error", err))
				continue
			}
			if err != nil {
				return err
			}
			runRes := database.RunResult{
				ID: best[i].RunResultId,
				Submission: database.Submission{
					SubmissionStatus: database.SubmissionStatusPending,
					SubmittedAt:      time.Now().UTC(),
					SubmissionToken:  token,
				},
			}
			err = database.UpdateEnsured(ctx, b.db.NewUpdate().
				Model(&runRes).
				WherePK().
				ExcludeColumn("solution").
				OmitZero())
			if err != nil {
				return err
			}
		}

		return nil
	})
}
