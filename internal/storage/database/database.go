package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/daremove/go-metrics-service/internal/storage"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"time"
)

type Database struct {
	db *sql.DB
}

func checkConnection(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

func (d *Database) Ping(ctx context.Context) error {
	return checkConnection(ctx, d.db)
}

func (d *Database) GetGaugeMetric(ctx context.Context, key string) (storage.GaugeMetric, error) {
	var result float64

	row := d.db.QueryRowContext(ctx, "SELECT value FROM gauge_metrics WHERE id = $1", key)

	if err := row.Scan(&result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.GaugeMetric{}, storage.ErrDataNotFound
		}

		return storage.GaugeMetric{}, err
	}

	return storage.GaugeMetric{Name: key, Value: result}, nil
}

func (d *Database) GetGaugeMetrics(ctx context.Context) ([]storage.GaugeMetric, error) {
	var result []storage.GaugeMetric

	rows, err := d.db.QueryContext(ctx, "SELECT * FROM gauge_metrics")

	if err != nil {
		return []storage.GaugeMetric{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var item storage.GaugeMetric

		if err := rows.Scan(&item.Name, &item.Value); err != nil {
			return []storage.GaugeMetric{}, err
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return []storage.GaugeMetric{}, err
	}

	return result, nil
}

func (d *Database) GetCounterMetric(ctx context.Context, key string) (storage.CounterMetric, error) {
	var result int64

	row := d.db.QueryRowContext(ctx, "SELECT value FROM counter_metrics WHERE id = $1", key)

	if err := row.Scan(&result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.CounterMetric{}, storage.ErrDataNotFound
		}

		return storage.CounterMetric{}, err
	}

	return storage.CounterMetric{Name: key, Value: result}, nil
}

func (d *Database) GetCounterMetrics(ctx context.Context) ([]storage.CounterMetric, error) {
	var result []storage.CounterMetric

	rows, err := d.db.QueryContext(ctx, "SELECT * FROM counter_metrics")

	if err != nil {
		return []storage.CounterMetric{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var item storage.CounterMetric

		if err := rows.Scan(&item.Name, &item.Value); err != nil {
			return []storage.CounterMetric{}, err
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return []storage.CounterMetric{}, err
	}

	return result, nil
}

func (d *Database) AddGauge(ctx context.Context, key string, value float64) error {
	if _, err := d.db.ExecContext(ctx, `
		INSERT INTO
			gauge_metrics (id, value)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET value = EXCLUDED.value
	`, key, value); err != nil {
		return err
	}

	return nil
}

func (d *Database) AddCounter(ctx context.Context, key string, value int64) error {
	if _, err := d.db.ExecContext(ctx, `
		INSERT INTO
			counter_metrics (id, value)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET value = EXCLUDED.value
	`, key, value); err != nil {
		return err
	}

	return nil
}

func New(ctx context.Context, dsn string) (*Database, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if err := checkConnection(ctx, db); err != nil {
		return nil, err
	}

	//migrationsPath, err := utils.GetRootPath("internal/storage/database/migrations")

	//if err != nil {
	//	return nil, err
	//}

	//if _, err := os.Stat(migrationsPath); errors.Is(err, os.ErrNotExist) {
	//	return &Database{db}, nil
	//}
	//
	//driver, err := postgres.WithInstance(db, &postgres.Config{})

	//if err != nil {
	//	return nil, err
	//}

	//m, err := migrate.NewWithDatabaseInstance(
	//	"file:"+migrationsPath,
	//	"postgres",
	//	driver,
	//)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//if err := m.Run(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
	//	return nil, err
	//}

	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS gauge_metrics (
		    id      VARCHAR PRIMARY KEY,
    		value   DOUBLE PRECISION NOT NULL
		);

		CREATE TABLE IF NOT EXISTS counter_metrics (
			id      VARCHAR PRIMARY KEY,
			value   INTEGER NOT NULL
		);
	`); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}
