package logs

import (
	"context"
	"icfpc/database"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type runResultKey struct{}

type contextHandler struct {
	slog.Handler
}

func WithRunResultLogging(ctx context.Context, runResult database.RunResult) context.Context {
	return context.WithValue(ctx, runResultKey{}, runResult)
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	runResult := ctx.Value(runResultKey{})

	if runResult != nil {
		runResult := runResult.(database.RunResult)
		r.AddAttrs(slog.Int64("task_id", runResult.TaskID))
		r.AddAttrs(slog.String("algorithm_name", runResult.AlgorithmName))
		r.AddAttrs(slog.String("algorithm_version", runResult.AlgorithmVersion))
	}

	return h.Handler.Handle(ctx, r)
}

func New() *slog.Logger {
	new := tint.NewHandler(os.Stdout, nil)
	withContext := (&contextHandler{new})
	return slog.New(withContext)
}
