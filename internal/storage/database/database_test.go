package database

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/daremove/go-metrics-service/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockDBExpectedResult struct {
	execCalls     int
	queryCalls    int
	queryRowCalls int
	beginTxCalls  int
	pingCalls     int
}

type MockDB struct {
	mu           sync.Mutex
	ExecFunc     func(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	QueryFunc    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRowFunc func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	BeginTxFunc  func(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	PingFunc     func(ctx context.Context) error

	ExecCalls     int
	QueryCalls    int
	QueryRowCalls int
	BeginTxCalls  int
	PingCalls     int

	expected MockDBExpectedResult
}

func (m *MockDB) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ExecCalls++
	return m.ExecFunc(ctx, sql, arguments...)
}

func (m *MockDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.QueryCalls++
	return m.QueryFunc(ctx, sql, args...)
}

func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.QueryRowCalls++
	return m.QueryRowFunc(ctx, sql, args...)
}

func (m *MockDB) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.BeginTxCalls++
	return m.BeginTxFunc(ctx, opts)
}

func (m *MockDB) Ping(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.PingCalls++
	return m.PingFunc(ctx)
}

func (m *MockDB) SetExpectedCalls(expected MockDBExpectedResult) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.expected = expected
}

func (m *MockDB) ExpectationsWereMet() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ExecCalls != m.expected.execCalls {
		return fmt.Errorf("expected Exec to be called %d times, but got %d", m.expected.execCalls, m.ExecCalls)
	}
	if m.QueryCalls != m.expected.queryCalls {
		return fmt.Errorf("expected Query to be called %d times, but got %d", m.expected.queryCalls, m.QueryCalls)
	}
	if m.QueryRowCalls != m.expected.queryRowCalls {
		return fmt.Errorf("expected QueryRow to be called %d times, but got %d", m.expected.queryRowCalls, m.QueryRowCalls)
	}
	if m.BeginTxCalls != m.expected.beginTxCalls {
		return fmt.Errorf("expected BeginTx to be called %d times, but got %d", m.expected.beginTxCalls, m.BeginTxCalls)
	}
	if m.PingCalls != m.expected.pingCalls {
		return fmt.Errorf("expected Ping to be called %d times, but got %d", m.expected.pingCalls, m.PingCalls)
	}
	return nil
}

func setupMockDB() (*Database, *MockDB) {
	mock := &MockDB{}
	db := &Database{db: mock}
	return db, mock
}

func TestDatabase(t *testing.T) {
	ctx := context.Background()

	t.Run("Should ping the database", func(t *testing.T) {
		db, mock := setupMockDB()
		mock.PingFunc = func(ctx context.Context) error {
			return nil
		}
		mock.SetExpectedCalls(MockDBExpectedResult{pingCalls: 1})

		err := db.Ping(ctx)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Should add and retrieve a gauge metric", func(t *testing.T) {
		db, mock := setupMockDB()
		key := "test_gauge"
		value := 1.23

		mock.ExecFunc = func(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
			return pgconn.CommandTag{}, nil
		}
		mock.QueryRowFunc = func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
			return &MockRow{
				ScanFunc: func(dest ...interface{}) error {
					if len(dest) == 1 {
						if v, ok := dest[0].(*float64); ok {
							*v = value
						}
					}
					return nil
				},
			}
		}
		mock.SetExpectedCalls(MockDBExpectedResult{execCalls: 1, queryRowCalls: 1})

		err := db.AddGaugeMetric(ctx, key, value)
		require.NoError(t, err)

		metric, err := db.GetGaugeMetric(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, key, metric.Name)
		assert.Equal(t, value, metric.Value)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return not found for nonexistent gauge metric", func(t *testing.T) {
		db, mock := setupMockDB()
		key := "nonexistent_gauge"

		mock.QueryRowFunc = func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
			return &MockRow{
				ScanFunc: func(dest ...interface{}) error {
					return pgx.ErrNoRows
				},
			}
		}
		mock.SetExpectedCalls(MockDBExpectedResult{queryRowCalls: 1})

		_, err := db.GetGaugeMetric(ctx, key)
		assert.ErrorIs(t, err, storage.ErrDataNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Should add and retrieve a counter metric", func(t *testing.T) {
		db, mock := setupMockDB()
		key := "test_counter"
		value := int64(123)

		mock.ExecFunc = func(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
			return pgconn.CommandTag{}, nil
		}
		mock.QueryRowFunc = func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
			return &MockRow{
				ScanFunc: func(dest ...interface{}) error {
					if len(dest) == 1 {
						if v, ok := dest[0].(*int64); ok {
							*v = value
						}
					}
					return nil
				},
			}
		}
		mock.SetExpectedCalls(MockDBExpectedResult{execCalls: 1, queryRowCalls: 1})

		err := db.AddCounterMetric(ctx, key, value)
		require.NoError(t, err)

		metric, err := db.GetCounterMetric(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, key, metric.Name)
		assert.Equal(t, value, metric.Value)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return not found for nonexistent counter metric", func(t *testing.T) {
		db, mock := setupMockDB()
		key := "nonexistent_counter"

		mock.QueryRowFunc = func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
			return &MockRow{
				ScanFunc: func(dest ...interface{}) error {
					return pgx.ErrNoRows
				},
			}
		}
		mock.SetExpectedCalls(MockDBExpectedResult{queryRowCalls: 1})

		_, err := db.GetCounterMetric(ctx, key)
		assert.ErrorIs(t, err, storage.ErrDataNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Should add multiple metrics in a transaction", func(t *testing.T) {
		db, mock := setupMockDB()
		gaugeMetrics := []storage.GaugeMetric{
			{Name: "gauge1", Value: 1.1},
			{Name: "gauge2", Value: 2.2},
		}
		counterMetrics := []storage.CounterMetric{
			{Name: "counter1", Value: 111},
			{Name: "counter2", Value: 222},
		}

		mock.ExecFunc = func(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
			return pgconn.CommandTag{}, nil
		}
		mock.BeginTxFunc = func(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
			return &MockTx{DB: mock}, nil
		}
		mock.SetExpectedCalls(MockDBExpectedResult{execCalls: 4, beginTxCalls: 1})

		err := db.AddMetrics(ctx, gaugeMetrics, counterMetrics)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Should retrieve all gauge metrics", func(t *testing.T) {
		db, mock := setupMockDB()
		expected := []storage.GaugeMetric{
			{Name: "gauge1", Value: 1.1},
			{Name: "gauge2", Value: 2.2},
		}

		mock.QueryFunc = func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
			return &MockRows{
				Rows: [][]interface{}{
					{"gauge1", 1.1},
					{"gauge2", 2.2},
				},
				Index: -1,
			}, nil
		}
		mock.SetExpectedCalls(MockDBExpectedResult{queryCalls: 1})

		retrieved, err := db.GetGaugeMetrics(ctx)
		require.NoError(t, err)
		assert.ElementsMatch(t, expected, retrieved)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Should retrieve all counter metrics", func(t *testing.T) {
		db, mock := setupMockDB()
		expected := []storage.CounterMetric{
			{Name: "counter1", Value: 111},
			{Name: "counter2", Value: 222},
		}

		mock.QueryFunc = func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
			return &MockRows{
				Rows: [][]interface{}{
					{"counter1", int64(111)},
					{"counter2", int64(222)},
				},
				Index: -1,
			}, nil
		}
		mock.SetExpectedCalls(MockDBExpectedResult{queryCalls: 1})

		retrieved, err := db.GetCounterMetrics(ctx)
		require.NoError(t, err)
		assert.ElementsMatch(t, expected, retrieved)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

type MockRow struct {
	ScanFunc func(dest ...interface{}) error
}

func (m *MockRow) Scan(dest ...interface{}) error {
	return m.ScanFunc(dest...)
}

type MockRows struct {
	Rows  [][]interface{}
	Index int
}

func (m *MockRows) Next() bool {
	m.Index++
	return m.Index < len(m.Rows)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	for i, d := range dest {
		switch v := d.(type) {
		case *string:
			*v = m.Rows[m.Index][i].(string)
		case *int:
			*v = m.Rows[m.Index][i].(int)
		case *int64:
			*v = m.Rows[m.Index][i].(int64)
		case *float64:
			*v = m.Rows[m.Index][i].(float64)
		default:
			return fmt.Errorf("unsupported scan destination type: %T", d)
		}
	}
	return nil
}

func (m *MockRows) Close()                        {}
func (m *MockRows) Err() error                    { return nil }
func (m *MockRows) CommandTag() pgconn.CommandTag { return pgconn.CommandTag{} }
func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	return []pgconn.FieldDescription{}
}
func (m *MockRows) Values() ([]any, error) {
	return []any{}, nil
}
func (m *MockRows) RawValues() [][]byte {
	return [][]byte{}
}
func (m *MockRows) Conn() *pgx.Conn {
	return nil
}

type MockTx struct {
	DB *MockDB
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) { return m, nil }
func (m *MockTx) Commit(ctx context.Context) error          { return nil }
func (m *MockTx) Rollback(ctx context.Context) error        { return nil }
func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	m.DB.ExecCalls += 1
	return pgconn.CommandTag{}, nil
}
func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return nil, nil
}
func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return &MockRow{}
}
func (m *MockTx) Prepare(ctx context.Context, name string, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}
func (m *MockTx) Conn() *pgx.Conn { return nil }
func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}
func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return nil
}
func (m *MockTx) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}
