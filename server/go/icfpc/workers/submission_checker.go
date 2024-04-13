package workers

import (
	"context"
	"errors"
	"icfpc/database"
	"icfpc/integration"
	"log/slog"
	"time"

	"github.com/uptrace/bun"
)

type submissionChecker struct {
	db  *bun.DB
	bus bus
}

func NewSubmissionChecker(db *bun.DB, bus bus) *submissionChecker {
	return &submissionChecker{
		db:  db,
		bus: bus,
	}
}

func (s submissionChecker) Run(
	ctx context.Context,
	fetch func(context.Context, []database.RunResult) (integration.CheckedSubmissions, error),
) error {
	return runPeriodical(ctx, time.Second, make(chan struct{}), func() error {
		var pending []database.RunResult
		q := s.db.NewSelect().
			Model(&pending).
			Where("submission_status = ?", database.SubmissionStatusPending).
			Relation("Task")
		if err := q.Scan(ctx); err != nil {
			return err
		}
		fetched, err := fetch(ctx, pending)
		if errors.Is(err, integration.Error) {
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
