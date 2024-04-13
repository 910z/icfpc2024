package server

import (
	"context"
	"icfpc/database"
	"icfpc/evaluation"
	"log/slog"
	"net/http"
	"strconv"
)

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	skip, err := strconv.ParseInt(r.URL.Query().Get("skip"), 10, 64)
	if err != nil {
		skip = 0
	}

	take, err := strconv.ParseInt(r.URL.Query().Get("take"), 10, 64)
	if err != nil {
		take = 10
	}

	page, err := s.getPage(r.Context(), int(skip), int(take))
	if err != nil {
		writeError(w, http.StatusInternalServerError)

		return
	}

	writeJson(w, page)
}

func (s *Server) getPage(ctx context.Context, skip int, take int) ([]database.RunEvalResult, error) {
	slog.DebugContext(ctx, "querying list...")
	var results []database.RunEvalResult

	err := s.db.NewSelect().
		Model(&results).
		Relation("RunResult.Task").
		Relation("RunResult").
		Where("version = ?", evaluation.Version).
		OrderExpr("task_id ASC").
		Limit(take).
		Offset(skip).
		Scan(ctx)

	for i := range results { // чтоб не лагал фронт. не селектить их через ExcludeColumn легко не получится
		var zeroData database.TaskData
		var zeroSolution database.Solution
		results[i].RunResult.Task.Data = zeroData
		results[i].RunResult.Solution = zeroSolution
	}

	if err != nil {
		slog.ErrorContext(ctx, "failed to query list", slog.Any("error", err))
		return nil, err
	}

	slog.DebugContext(ctx, "queried list")
	return results, nil
}
