package main

import (
	"context"
	"icfpc/database"
	"icfpc/front"
	"icfpc/integration"
	"icfpc/logs"
	"icfpc/runner"
	"icfpc/server"
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

	logs.SetDefaultSlog()

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
		if err := algoRunner.RunAlgorithms(ctx, runner.AllAlgorithms); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := runner.RunTasksFetcher(ctx, integration.GetTasks, db); err != nil {
			panic(err)
		}
	}()

	srv := server.New(db, 8080)
	srv.SetupEndpoints()

	if err = srv.Run(); err != nil {
		panic(err)
	}
}
