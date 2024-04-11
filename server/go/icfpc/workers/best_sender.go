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

	"github.com/uptrace/bun"
)

type bestSender struct {
	db *bun.DB
}

func NewBestSender(db *bun.DB) *bestSender {
	return &bestSender{
		db: db,
	}
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

type extEvalResult struct {
	database.RunEvalResult
	Solution   database.Solution
	ExternalID string

	// на самом деле RowNum не нужен.
	// это просто чтобы не перечислять вручную все поля runevalresult-а, а сделать select *.
	// хотя, может, можно селектить только score, а апдейтить только ненулевые поля
	RowNum int
}

func (b bestSender) Run(
	ctx context.Context,
	ord SortOrder,
	send func(context.Context, string, database.Solution) error,
) error {
	return runPeriodical(ctx, time.Second, func(ctx context.Context) error {
		var best []extEvalResult
		err := b.db.NewRaw(fmt.Sprintf(`with allRes as (
			select
			  run_eval_results.*,
			  run_results.solution,
			  tasks.external_id,
			  row_number() over w as row_num
			from
			  run_eval_results
			  join run_results on run_results.id = run_eval_results.run_result_id
			  join tasks on tasks.id = run_results.task_id
			where
			  run_eval_results.version = ?
			  AND run_eval_results.status = ?
			  AND run_eval_results.submission_status = ''
			window w as (
				partition by run_results.task_id
				order by
				  score %s
			  )
		  )
		  select
			*
		  from
			allRes
		  where
		  	row_num = 1`, ord), evaluation.Version, database.ProgressStatusFinished).Scan(ctx, &best)
		slog.InfoContext(ctx, "sending best", slog.Int("count", len(best)))
		if err != nil {
			return err
		}
		for i := range best {
			err = send(ctx, best[i].ExternalID, best[i].Solution)
			if errors.Is(err, integration.Error) {
				continue
			}
			if err != nil {
				return err
			}
			best[i].SubmissionStatus = database.SubmissionStatusPending
			best[i].SubmittedAt = time.Now().UTC()
			err := database.UpdateEnsured(ctx, b.db.NewUpdate().Model(&best[i].RunEvalResult).WherePK())
			if err != nil {
				return err
			}
		}

		return nil
	})
}
