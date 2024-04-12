package workers

import (
	"context"
	"time"

	"icfpc/database"

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
	fetch func(context.Context, []database.RunEvalResult) error,
) error {
	return runPeriodical(ctx, time.Second, func(ctx context.Context) error {
		return nil
	})
}
