package main

import (
	"context"
	"icfpc/database"
	"icfpc/evaluation"
	"icfpc/front"
	"icfpc/runner"
	"icfpc/server"
	_ "net/http/pprof"
	"os"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// TODO: what is this?
type RunEvalResult struct {
	evaluation.Result

	RunResult database.RunResult `bun:"rel:belongs-to,join:run_result_id=id"`
}

func main() {
	if app.IsClient {
		app.Route("/", &front.RunList{})
		app.RunWhenOnBrowser()

		return
	}

	ctx := context.Background()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" { // в тесте
		connStr = "postgresql://postgres:password@localhost/postgres?sslmode=disable"
	}

	db, err := database.SetUp(ctx, connStr)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	algoRunner := runner.New(db)

	go func() {
		if err := algoRunner.Run(ctx, runner.AllTasks, runner.AllAlgorithms); err != nil {
			panic(err)
		}
	}()

	srv := server.New(db, 8080)
	srv.SetupEndpoints()

	if err = srv.Run(); err != nil {
		panic(err)
	}
}
