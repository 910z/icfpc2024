package workers

import (
	"context"
	"icfpc/database"
	"icfpc/integration"
	"log/slog"
	"time"

	"github.com/uptrace/bun"
)

type submissionChecker struct {
	db *bun.DB
}

func NewSubmissionChecker(db *bun.DB) *submissionChecker {
	return &submissionChecker{
		db: db,
	}
}

func (s submissionChecker) Run(
	ctx context.Context,
	fetch func(context.Context, []database.RunResult) (integration.CheckedSubmissions, error),
) error {
	return runPeriodical(ctx, time.Second, func(ctx context.Context) error {
		var pending []database.RunResult
		q := s.db.NewSelect().
			Model(&pending).
			Where("submission_status = ?", database.SubmissionStatusPending).
			Relation("Task")
		if err := q.Scan(ctx); err != nil {
			return err
		}
		fetched, err := fetch(ctx, pending)
		if err == integration.Error {
			slog.WarnContext(ctx, "can't fetch submissions, retrying", slog.Any("error", err))
			return nil
		}
		if err != nil {
			return err
		}
		for _, pr := range pending {
			check, ok := fetched[pr.Task.ExternalID]
			if !ok {
				// еще пендится
				continue
			}
			pr.SubmissionCheckedAt = time.Now().UTC()
			pr.ExternalScore = check.Score
			pr.SubmissionStatus = check.Status

			if err := database.UpdateEnsured(ctx, s.db.NewUpdate().Model(&pr).WherePK()); err != nil {
				return err
			}
		}

		return nil
	})
}
