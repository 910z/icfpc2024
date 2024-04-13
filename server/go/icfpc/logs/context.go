package logs

import (
	"context"
	"icfpc/database"
	"log/slog"
	"os"
	"reflect"

	"github.com/lmittmann/tint"
)

type runResultKey struct{}
type typeKey struct{}

type contextHandler struct {
	slog.Handler
}

func WithType(ctx context.Context, t reflect.Type) context.Context {
	return context.WithValue(ctx, typeKey{}, t)
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

	callerType := ctx.Value(typeKey{})
	if callerType != nil {
		r.AddAttrs(slog.String("type", callerType.(reflect.Type).Name()))
	}

	return h.Handler.Handle(ctx, r)
}

func New(level slog.Level) *slog.Logger {
	new := tint.NewHandler(os.Stdout, &tint.Options{
		Level: level,
	})
	withContext := (&contextHandler{new})
	return slog.New(withContext)
}
