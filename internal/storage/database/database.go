// Package database предоставляет функциональность для работы с базой данных,
// включая операции добавления и извлечения метрик.
package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/daremove/go-metrics-service/internal/storage"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Запросы для вставки метрик в базу данных с обработкой конфликтов.
const (
	InsertGaugeMetricQuery = `
		INSERT INTO
			gauge_metrics (id, value)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET value = EXCLUDED.value
	`
	InsertCounterMetricQuery = `
		INSERT INTO
			counter_metrics (id, value)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET value = counter_metrics.value + EXCLUDED.value
	`
)

type DB interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	Ping(context.Context) error
}

// Database структура для взаимодействия с базой данных.
type Database struct {
	db DB
}

// checkConnection проверяет соединение с базой данных.
func checkConnection(ctx context.Context, db DB) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		return err
	}

	return nil
}

// Ping проверяет доступность базы данных.
func (d *Database) Ping(ctx context.Context) error {
	return checkConnection(ctx, d.db)
}

// GetGaugeMetric извлекает метрику типа gauge из базы данных.
func (d *Database) GetGaugeMetric(ctx context.Context, key string) (storage.GaugeMetric, error) {
	var result float64

	row := d.db.QueryRow(ctx, "SELECT value FROM gauge_metrics WHERE id = $1", key)

	if err := row.Scan(&result); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.GaugeMetric{}, storage.ErrDataNotFound
		}

		return storage.GaugeMetric{}, err
	}

	return storage.GaugeMetric{Name: key, Value: result}, nil
}

// GetGaugeMetrics извлекает все метрики типа gauge из базы данных.
func (d *Database) GetGaugeMetrics(ctx context.Context) ([]storage.GaugeMetric, error) {
	var result []storage.GaugeMetric

	rows, err := d.db.Query(ctx, "SELECT id, value FROM gauge_metrics")

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

// GetCounterMetric извлекает метрику типа counter из базы данных.
func (d *Database) GetCounterMetric(ctx context.Context, key string) (storage.CounterMetric, error) {
	var result int64

	row := d.db.QueryRow(ctx, "SELECT value FROM counter_metrics WHERE id = $1", key)

	if err := row.Scan(&result); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.CounterMetric{}, storage.ErrDataNotFound
		}

		return storage.CounterMetric{}, err
	}

	return storage.CounterMetric{Name: key, Value: result}, nil
}

// GetCounterMetrics извлекает все метрики типа counter из базы данных.
func (d *Database) GetCounterMetrics(ctx context.Context) ([]storage.CounterMetric, error) {
	var result []storage.CounterMetric

	rows, err := d.db.Query(ctx, "SELECT id, value FROM counter_metrics")

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

// AddGaugeMetric добавляет или обновляет метрику типа gauge в базе данных.
func (d *Database) AddGaugeMetric(ctx context.Context, key string, value float64) error {
	if _, err := d.db.Exec(ctx, InsertGaugeMetricQuery, key, value); err != nil {
		return err
	}

	return nil
}

// AddCounterMetric добавляет или обновляет метрику типа counter в базе данных.
func (d *Database) AddCounterMetric(ctx context.Context, key string, value int64) error {
	if _, err := d.db.Exec(ctx, InsertCounterMetricQuery, key, value); err != nil {
		return err
	}

	return nil
}

// AddMetrics добавляет или обновляет несколько метрик в базе данных в рамках одной транзакции.
func (d *Database) AddMetrics(ctx context.Context, gaugeMetrics []storage.GaugeMetric, counterMetrics []storage.CounterMetric) error {
	tx, err := d.db.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	for _, gaugeMetric := range gaugeMetrics {
		if _, err := tx.Exec(ctx, InsertGaugeMetricQuery, gaugeMetric.Name, gaugeMetric.Value); err != nil {
			return err
		}
	}

	for _, counterMetric := range counterMetrics {
		if _, err := tx.Exec(ctx, InsertCounterMetricQuery, counterMetric.Name, counterMetric.Value); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// New инициализирует и возвращает новый экземпляр Database.
func New(ctx context.Context, dsn string) (*Database, error) {
	db, err := pgxpool.New(ctx, dsn)

	if err != nil {
		return nil, err
	}

	if err := checkConnection(ctx, db); err != nil {
		return nil, err
	}

	if _, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS gauge_metrics (
		    id      TEXT PRIMARY KEY,
    		value   DOUBLE PRECISION NOT NULL
		);

		CREATE TABLE IF NOT EXISTS counter_metrics (
			id      TEXT PRIMARY KEY,
			value   BIGSERIAL NOT NULL
		);
	`); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}
