package main

import (
	"context"
	"icfpc/database"
	"icfpc/front"
	"icfpc/integration"
	"icfpc/logs"
	"icfpc/server"
	"icfpc/workers"
	"log/slog"
	_ "net/http/pprof"
	"os"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	app.Route("/", &front.RunList{})

	if app.IsClient {
		app.RunWhenOnBrowser()

		return
	}

	ctx := context.Background()

	connStr := os.Getenv("DATABASE_URL")
	isDevEnv := connStr == ""
	if isDevEnv {
		slog.SetDefault(logs.New(slog.LevelDebug))
		connStr = "postgresql://postgres:password@localhost/postgres?sslmode=disable"
	} else {
		slog.SetDefault(logs.New(slog.LevelInfo))
	}

	db, err := database.SetUp(ctx, connStr)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = db.Close()
	}()

	bus := workers.NewBus()
	algoRunner := workers.NewAlgorithmRunner(db, bus)
	taskFetcher := workers.NewTasksFetcher(db, bus)
	solutionEvaluator := workers.NewSolutionEvaluator(db, bus)
	bestSender := workers.NewBestSender(db, bus)
	checker := workers.NewSubmissionChecker(db, bus)

	go func() {
		if err := checker.Run(ctx, integration.GetSubmissionsStatus); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := bestSender.Run(ctx, workers.SortOrderDesc, integration.SendSolution); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := algoRunner.Run(ctx, workers.AllAlgorithms); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := taskFetcher.Run(ctx, integration.GetTasks); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := solutionEvaluator.Run(ctx); err != nil {
			panic(err)
		}
	}()

	srv := server.New(db, 8080)
	srv.SetupEndpoints()

	if err = srv.Run(); err != nil {
		panic(err)
	}
}
