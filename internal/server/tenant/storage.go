package tenant

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	corestorage "ttl-cli/internal/core/storage"
	storagesqlite "ttl-cli/internal/storage/sqlite"
)

type StorageManager struct {
	dataDir  string
	storages map[string]*storagesqlite.SQLiteStorage
	mu       sync.RWMutex
}

func NewStorageManager(dataDir string) *StorageManager {
	return &StorageManager{dataDir: dataDir, storages: make(map[string]*storagesqlite.SQLiteStorage)}
}

func (m *StorageManager) GetStorage(userID string) (corestorage.Storage, error) {
	m.mu.RLock()
	if storage, ok := m.storages[userID]; ok {
		m.mu.RUnlock()
		return storage, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()
	if storage, ok := m.storages[userID]; ok {
		return storage, nil
	}

	if err := os.MkdirAll(m.dataDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create tenant data directory: %w", err)
	}
	if err := os.Chmod(m.dataDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to set tenant data directory permissions: %w", err)
	}
	userDir := filepath.Join(m.dataDir, userID)
	if err := os.MkdirAll(userDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create user data directory: %w", err)
	}
	if err := os.Chmod(userDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to set user data directory permissions: %w", err)
	}
	for _, legacyName := range []string{"data.db", "data.bbolt"} {
		legacyPath := filepath.Join(userDir, legacyName)
		if _, err := os.Stat(legacyPath); err == nil {
			return nil, fmt.Errorf("user %s has incompatible legacy database %s", userID, legacyName)
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to inspect legacy database for user %s: %w", userID, err)
		}
	}
	storage := storagesqlite.NewSQLiteStorage()
	storage.SetDBPath(filepath.Join(userDir, "data.sqlite"))
	if err := storage.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize database for user %s: %w", userID, err)
	}
	m.storages[userID] = storage
	return storage, nil
}

func (m *StorageManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var lastErr error
	for id, storage := range m.storages {
		if err := storage.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close storage for user %s: %w", id, err)
		}
		delete(m.storages, id)
	}
	return lastErr
}

func (m *StorageManager) RemoveStorage(userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if storage, ok := m.storages[userID]; ok {
		if err := storage.Close(); err != nil {
			return fmt.Errorf("failed to close storage for user %s: %w", userID, err)
		}
		delete(m.storages, userID)
	}
	userDir := filepath.Join(m.dataDir, userID)
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to inspect user data directory: %w", err)
	}
	deletedDir := filepath.Join(m.dataDir, ".deleted")
	if err := os.MkdirAll(deletedDir, 0700); err != nil {
		return fmt.Errorf("failed to create deleted data directory: %w", err)
	}
	target := filepath.Join(deletedDir, fmt.Sprintf("%s-%d", userID, time.Now().UnixNano()))
	if err := os.Rename(userDir, target); err != nil {
		return fmt.Errorf("failed to isolate user data directory: %w", err)
	}
	return nil
}
