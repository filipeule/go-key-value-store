package transact

import (
	"database/sql"
	"fmt"
	"key-value-store/core"

	_ "github.com/lib/pq"
)

type PostgresDBParams struct {
	DBName   string
	Host     string
	User     string
	Password string
}

type PostgresTransactionLogger struct {
	events chan<- core.Event
	errors <-chan error
	db     *sql.DB
}

func NewPostgresTransactionLogger(config PostgresDBParams) (core.TransactionLogger, error) {
	connStr := fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s sslmode=disable",
		config.Host,
		config.DBName,
		config.User,
		config.Password,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	logger := &PostgresTransactionLogger{
		db: db,
	}

	exists, err := logger.verifyTableExists(config.DBName)
	if err != nil {
		return nil, fmt.Errorf("failed to verify if table exists: %w", err)
	}

	if !exists {
		if err = logger.createTable(config.DBName); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	return logger, nil
}

func (pl *PostgresTransactionLogger) Run() {
	events := make(chan core.Event, 16)
	pl.events = events

	errors := make(chan error, 1)
	pl.errors = errors

	go func() {
		query := `
		INSERT INTO
			transactions(event_type, key, value)
		VALUES
			($1, $2, $3)
		`

		for e := range events {
			_, err := pl.db.Exec(query, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}
		}
	}()
}

func (pl *PostgresTransactionLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	outEvent := make(chan core.Event)
	outError := make(chan error)

	go func() {
		defer close(outEvent)
		defer close(outError)

		query := `
		SELECT
			sequence, event_type, key, value
		FROM
			transactions
		ORDER BY
			sequence
		`

		rows, err := pl.db.Query(query)
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}
		defer rows.Close()

		e := core.Event{}
		for rows.Next() {
			err = rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outError <- fmt.Errorf("error reading row: %w", err)
				return
			}

			outEvent <- e
		}

		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()

	return outEvent, outError
}

func (pl *PostgresTransactionLogger) WritePut(key, value string) {
	pl.events <- core.Event{EventType: core.EventPut, Key: key, Value: value}
}

func (pl *PostgresTransactionLogger) WriteDelete(key string) {
	pl.events <- core.Event{EventType: core.EventDelete, Key: key}
}

func (pl *PostgresTransactionLogger) Err() <-chan error {
	return pl.errors
}

func (pl *PostgresTransactionLogger) verifyTableExists(table string) (bool, error) {
	query := fmt.Sprintf(`
	SELECT EXISTS (
   		SELECT FROM pg_catalog.pg_class c
   		JOIN   pg_catalog.pg_namespace n ON n.oid = c.relnamespace
   		WHERE  n.nspname = 'public'
   		AND    c.relname = '%s'
   	)
	`, table)

	row := pl.db.QueryRow(query)

	var result bool

	err := row.Scan(&result)
	if err != nil {
		return false, fmt.Errorf("error scanning verify table exists result: %w", err)
	}

	return result, nil
}

func (pl *PostgresTransactionLogger) createTable(table string) error {
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		sequence SERIAL PRIMARY KEY,
		event_type INT NOT NULL,
		key TEXT NOT NULL,
		value TEXT
	)
	`, table)

	_, err := pl.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}
