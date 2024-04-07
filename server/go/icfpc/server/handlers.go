package server

import (
	"context"
	"icfpc/database"
	"log"
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

func (s *Server) getPage(ctx context.Context, skip int, take int) ([]database.RunResult, error) {
	log.Println("querying list...")

	var results []database.RunResult

	err := s.db.NewSelect().
		Model(&results).
		OrderExpr("task_id ASC").
		Limit(take).
		Offset(skip).
		Scan(ctx)
	if err != nil {
		log.Printf("error querying list: %s", err.Error())

		return nil, err
	}

	log.Println("queried list")

	return results, nil
}
