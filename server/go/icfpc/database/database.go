package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func SetUp(ctx context.Context, connStr string) (*bun.DB, error) {
	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr))), pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		// bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	if os.Getenv("DROP_TABLES") == "true" {
		if err := dropTables(ctx, db); err != nil {
			return nil, err
		}
	}

	err := createSchema(ctx, db)
	for err != nil {
		log.Printf("can't connect, retrying in 1 sec: %s\n", err.Error())
		time.Sleep(time.Second)

		err = createSchema(ctx, db)
	}

	return db, nil
}

func createSchema(ctx context.Context, db *bun.DB) error {
	for _, model := range allModels {
		if _, err := db.NewCreateTable().IfNotExists().Model(model).Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

func dropTables(ctx context.Context, db *bun.DB) error {
	for _, model := range allModels {
		if _, err := db.NewDropTable().Model(model).IfExists().Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}
