package main

import (
	"context"
	"database/sql"
	"fmt"
	"icfpc/front"
	"icfpc/types"
	"log"
	"net/http"
	"os"
	"time"

	_ "net/http/pprof"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var models = []interface{}{
	(*types.RunResult)(nil),
	(*types.Task)(nil),
}

func createSchema(db *bun.DB) error {
	ctx := context.Background()

	for _, model := range models {
		_, err := db.NewCreateTable().IfNotExists().Model(model).Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func drop(db *bun.DB) error {
	ctx := context.Background()

	for _, model := range models {
		_, err := db.NewDropTable().Model(model).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {

	app.Route("/", &front.RunList{})
	app.RunWhenOnBrowser()
	if app.IsClient {
		return
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" { // в тесте
		dsn = "postgresql://postgres:password@localhost/postgres?sslmode=disable"
	}
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	db.AddQueryHook(bundebug.NewQueryHook(
		// bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	if os.Getenv("DROP_TABLES") == "true" {
		err := drop(db)
		if err != nil {
			panic(err)
		}
	}
	err := createSchema(db)
	for err != nil {
		log.Println("can't connect, retrying in 1 sec", err)
		time.Sleep(time.Second)
		err = createSchema(db)
	}

	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})
	serveDb(db)
	go runAlgorithms(db)

	port := 8080
	log.Printf("Server listening on port %d...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
