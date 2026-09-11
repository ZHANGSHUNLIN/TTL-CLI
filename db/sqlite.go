package db

import storagesqlite "ttl-cli/internal/storage/sqlite"

// SQLiteStorage remains as a compatibility alias while callers migrate to internal/storage/sqlite.
type SQLiteStorage = storagesqlite.SQLiteStorage

func NewSQLiteStorage() *SQLiteStorage {
	return storagesqlite.NewSQLiteStorage()
}
