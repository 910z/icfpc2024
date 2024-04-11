package workers

import (
	"context"
	"errors"
	"time"

	"icfpc/database"
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

func (b bestSender) Run(
	ctx context.Context,
	ord SortOrder,
	send func(context.Context, database.RunEvalResult) error,
) error {
	return runPeriodical(ctx, time.Second, func(ctx context.Context) error {
		var best database.RunEvalResult
		q := b.db.NewSelect().
			Model(&best).
			Where("status = ?", database.ProgressStatusFinished).
			Where("error IS NULL OR error = ''").
			Order("score " + string(ord)).Limit(1)
		if err := q.Scan(ctx); err != nil {
			return err
		}
		err := send(ctx, best)
		if errors.Is(err, integration.Error) {
			return nil
		}
		if err != nil {
			return err
		}

		best.SubmissionStatus = database.SubmissionStatusPending
		best.SubmittedAt = time.Now().UTC()
		return database.UpdateEnsured(ctx, b.db.NewUpdate().Model(&best).WherePK())
	})
}
