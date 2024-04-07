package main

import (
	"context"
	"encoding/json"
	"icfpc/types"
	"log"
	"net/http"
	"strconv"

	"github.com/uptrace/bun"
)

func GetPage(db *bun.DB, skip int, take int) ([]types.RunResult, error) {
	log.Print("querying list")
	ctx := context.Background()
	var results []types.RunResult
	err := db.NewSelect().
		Model(&results).
		OrderExpr("task_id ASC").
		Limit(take).
		Offset(skip).
		Scan(ctx)
	log.Print("queried list", err)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func writeJson(w http.ResponseWriter, v interface{}, err error) {
	if err != nil {
		serveError(w, 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func serveError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func serveDb(db *bun.DB) {

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		skip64, err := strconv.ParseInt(r.URL.Query().Get("skip"), 10, 64)
		skip := int(skip64)
		if err != nil {
			skip = 0
		}

		take64, err := strconv.ParseInt(r.URL.Query().Get("take"), 10, 64)
		take := int(take64)
		if err != nil {
			take = 10
		}

		page, err := GetPage(db, skip, take)
		writeJson(w, page, err)
	})

}
