package db

import (
	corestorage "ttl-cli/internal/core/storage"
	storagebbolt "ttl-cli/internal/storage/bbolt"
)

type Storage = corestorage.Storage

// LocalStorage remains as a compatibility alias while callers migrate to internal/storage/bbolt.
type LocalStorage = storagebbolt.LocalStorage

func NewLocalStorage() *LocalStorage {
	return storagebbolt.NewLocalStorage()
}

func GetDBPath(confFile, storageType string) (string, error) {
	return storagebbolt.GetDBPath(confFile, storageType)
}
