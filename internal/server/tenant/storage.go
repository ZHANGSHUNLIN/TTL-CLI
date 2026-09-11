package tenant

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	corestorage "ttl-cli/internal/core/storage"
	storagebbolt "ttl-cli/internal/storage/bbolt"
)

type StorageManager struct {
	dataDir  string
	storages map[string]*storagebbolt.LocalStorage
	mu       sync.RWMutex
}

func NewStorageManager(dataDir string) *StorageManager {
	return &StorageManager{dataDir: dataDir, storages: make(map[string]*storagebbolt.LocalStorage)}
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

	userDir := filepath.Join(m.dataDir, userID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create user data directory: %w", err)
	}
	storage := storagebbolt.NewLocalStorage()
	storage.SetDBPath(filepath.Join(userDir, "data.db"))
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
		_ = storage.Close()
		delete(m.storages, userID)
	}
	return os.RemoveAll(filepath.Join(m.dataDir, userID))
}
