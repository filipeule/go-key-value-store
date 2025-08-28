package transact

import (
	"errors"
	"fmt"
	"key-value-store/core"
	"os"
)

func NewTransactionLogger(logger string) (core.TransactionLogger, error) {
	switch logger {
	case "file":
		return NewFileTransactionLogger(os.Getenv("TLOG_FILENAME"))
	case "postgres":
		return NewPostgresTransactionLogger(PostgresDBParams{
			Host: os.Getenv("TLOG_DB_HOST"),
			DBName: os.Getenv("TLOG_DB_DATABASE"),
			User: os.Getenv("TLOG_DB_USERNAME"),
			Password: os.Getenv("TLOG_DB_PASSWORD"),
		})
	case "":
		return nil, errors.New("transaction logger type not defined")
	default:
		return nil, fmt.Errorf("no such transaction logger %s", logger)
	}
}