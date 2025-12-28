package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/cabrerajulian401/ecom/internal/env"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=ecom sslmode=disable"),
			/* A DSN in databases is a stored configuration that simplifies connecting
			applications to databses by holding all necessary connection details, like the driver,
			database location, and login info under 1 reference name */

		},
	}

	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	slog.Info("connected to database", "dsn", cfg.db.dsn)

	// Database

	api := application{
		config: cfg,
		db:     conn,
	}

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	err = api.run(api.mount())

	if err != nil {
		slog.Error("server has failed to start", "error", err)
		os.Exit(1)
	}
}
