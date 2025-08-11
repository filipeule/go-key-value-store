package main

import (
	"errors"
	"fmt"
	"key-value-store/internal/transaction"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	store = struct {
		sync.RWMutex
		m map[string]string
	}{
		m: make(map[string]string),
	}

	logger transaction.TransactionLogger

	ErrNoSuchKey = errors.New("no such key")
	ErrNoTransactionFound = errors.New("no transaction found")
)

func initializeTransactionLog(transactionType transaction.TransactionType) error {
	var err error

	switch transactionType {
	case transaction.PostgresTransaction:
		logger, err = transaction.NewPostgresTransactionLogger(transaction.PostgresDBParams{
			Host:     "localhost",
			DBName:   "transactions",
			User:     "postgres",
			Password: "postgres",
		})
		if err != nil {
			return fmt.Errorf("failed to create event logger: %w", err)
		}
	case transaction.FileTransaction:
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		loggerPath := filepath.Join(filepath.Dir(exePath), "transaction.log")

		logger, err = transaction.NewFileTransactionLogger(loggerPath)
		if err != nil {
			return fmt.Errorf("failed to create event logger: %w", err)
		}
	default:
		return ErrNoTransactionFound
	}

	events, errors := logger.ReadEvents()

	e, ok := transaction.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case transaction.EventDelete:
				err = Delete(e.Key)
			case transaction.EventPut:
				err = Put(e.Key, e.Value)
			}
		}
	}

	logger.Run()

	go func() {
		for err := range logger.Err() {
			log.Printf("error from transaction logger: %s\n", err)
		}
	}()

	return err
}

func Get(key string) (string, error) {
	store.RLock()
	val, ok := store.m[key]
	store.RUnlock()
	if !ok {
		return "", ErrNoSuchKey
	}

	return val, nil
}

func Put(key, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()

	return nil
}

func Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()

	return nil
}
