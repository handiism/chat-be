package psql

import (
	"context"
	"fmt"
	"log"
	"os"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool() *pgxpool.Pool {
	username := os.Getenv("PSQL_USERNAME")
	password := os.Getenv("PSQL_PASSWORD")
	host := os.Getenv("PSQL_HOST")
	port := os.Getenv("PSQL_PORT")
	database := os.Getenv("PSQL_DBNAME")

	url := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s", host, username, password, port, database)
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err.Error())
	}

	config.AfterConnect = func(_ context.Context, c *pgx.Conn) error {
		pgxuuid.Register(c.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf(err.Error())
	}

	return pool
}
