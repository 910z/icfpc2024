package database

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func waitForConnection(ctx context.Context, db *bun.DB) {
	for {
		err := db.PingContext(ctx)
		if err == nil {
			slog.InfoContext(ctx, "connected to db")
			break
		} else {
			slog.WarnContext(ctx, "can't connect, retrying in 1 sec", slog.Any("error", err))
			time.Sleep(time.Second)
		}
	}
}

func SetUp(ctx context.Context, connStr string) (*bun.DB, error) {
	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr))), pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.FromEnv("BUNDEBUG"),
	))

	if os.Getenv("DROP_TABLES") == "true" {
		if err := dropTables(ctx, db); err != nil {
			return nil, err
		}
	}

	waitForConnection(ctx, db)

	if err := createSchema(ctx, db); err != nil {
		return nil, err
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
