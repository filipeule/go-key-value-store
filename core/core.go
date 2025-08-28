package core

import (
	"errors"
	"sync"
)

var (
	ErrNoSuchKey = errors.New("no such key")
)

type KeyValueStore struct {
	m map[string]string
	transact TransactionLogger
	sync.RWMutex
}

func NewKeyValueStore(tl TransactionLogger) *KeyValueStore {
	return &KeyValueStore{
		m: make(map[string]string),
		transact: tl,
	}
}

func (store *KeyValueStore) Get(key string) (string, error) {
	store.RLock()
	val, ok := store.m[key]
	store.RUnlock()
	if !ok {
		return "", ErrNoSuchKey
	}

	return val, nil
}

func (store *KeyValueStore) Put(key, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()

	store.transact.WritePut(key, value)

	return nil
}

func (store *KeyValueStore) Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()

	store.transact.WriteDelete(key)

	return nil
}