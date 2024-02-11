package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Database struct {
	dsn string
}

func (s *Database) Ping() error {
	dbpool, err := pgxpool.New(context.Background(), s.dsn)

	if err != nil {
		return err
	}

	defer dbpool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	if err := dbpool.Ping(ctx); err != nil {
		return err
	}

	return nil
}

func New(dsn string) *Database {
	return &Database{dsn}
}
