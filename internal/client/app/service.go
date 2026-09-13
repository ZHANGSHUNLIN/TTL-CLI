package app

import (
	"context"
	"fmt"
	"time"

	"ttl-cli/internal/client/remote"
	clientsync "ttl-cli/internal/client/sync"
	"ttl-cli/internal/config"
	"ttl-cli/internal/core/resource"
	corestorage "ttl-cli/internal/core/storage"
	storagebbolt "ttl-cli/internal/storage/bbolt"
	storagesqlite "ttl-cli/internal/storage/sqlite"
)

type contextKey struct{}

// Service owns the storage dependency for one client command execution.
type Service struct {
	storage corestorage.Storage
}

func NewService(storage corestorage.Storage) *Service {
	return &Service{storage: storage}
}

func WithService(ctx context.Context, service *Service) context.Context {
	return context.WithValue(ctx, contextKey{}, service)
}

func ServiceFromContext(ctx context.Context) (*Service, error) {
	service, ok := ctx.Value(contextKey{}).(*Service)
	if !ok || service == nil {
		return nil, fmt.Errorf("storage not initialized")
	}
	return service, nil
}

func (s *Service) Storage() corestorage.Storage {
	if s == nil {
		return nil
	}
	return s.storage
}

func (s *Service) LocalStorage() (*storagebbolt.LocalStorage, bool) {
	localStorage, ok := s.Storage().(*storagebbolt.LocalStorage)
	return localStorage, ok
}

func (s *Service) Close() error {
	if s == nil || s.storage == nil {
		return nil
	}
	return s.storage.Close()
}

func (s *Service) storageOrError() (corestorage.Storage, error) {
	if s == nil || s.storage == nil {
		return nil, fmt.Errorf("storage not initialized")
	}
	return s.storage, nil
}

func (s *Service) GetAllResources() (map[resource.ValJsonKey]resource.ValJson, error) {
	storage, err := s.storageOrError()
	if err != nil {
		return nil, err
	}
	return storage.GetAllResources()
}

func (s *Service) SaveResource(key resource.ValJsonKey, value resource.ValJson) error {
	storage, err := s.storageOrError()
	if err != nil {
		return err
	}
	return storage.SaveResource(key, value)
}

func (s *Service) DeleteResource(key resource.ValJsonKey) error {
	storage, err := s.storageOrError()
	if err != nil {
		return err
	}
	return storage.DeleteResource(key)
}

func (s *Service) UpdateResource(key resource.ValJsonKey, value resource.ValJson) error {
	storage, err := s.storageOrError()
	if err != nil {
		return err
	}
	return storage.UpdateResource(key, value)
}

func (s *Service) GetTagStats() ([]resource.TagStat, error) {
	storage, err := s.storageOrError()
	if err != nil {
		return nil, err
	}
	return storage.GetTagStats()
}

func (s *Service) SaveAuditRecord(record resource.AuditRecord) error {
	storage, err := s.storageOrError()
	if err != nil {
		return err
	}
	return storage.SaveAuditRecord(record)
}

func (s *Service) GetAuditStats() (resource.AuditStats, error) {
	storage, err := s.storageOrError()
	if err != nil {
		return resource.AuditStats{}, err
	}
	return storage.GetAuditStats()
}

func (s *Service) GetAllAuditRecords() ([]resource.AuditRecord, error) {
	storage, err := s.storageOrError()
	if err != nil {
		return nil, err
	}
	return storage.GetAllAuditRecords()
}

func (s *Service) DeleteAuditRecords(resourceKey string) error {
	storage, err := s.storageOrError()
	if err != nil {
		return err
	}
	return storage.DeleteAuditRecords(resourceKey)
}

func (s *Service) GetAllHistoryRecords() ([]resource.HistoryRecord, error) {
	storage, err := s.storageOrError()
	if err != nil {
		return nil, err
	}
	return storage.GetAllHistoryRecords()
}

func (s *Service) DeleteHistoryRecords(resourceKey string) error {
	storage, err := s.storageOrError()
	if err != nil {
		return err
	}
	return storage.DeleteHistoryRecords(resourceKey)
}

func (s *Service) GetHistoryRecord(index int, order resource.SortOrder) (resource.HistoryRecord, error) {
	storage, err := s.storageOrError()
	if err != nil {
		return resource.HistoryRecord{}, err
	}
	return storage.GetHistoryRecord(index, order)
}

func (s *Service) SaveHistoryRecord(record resource.HistoryRecord) error {
	storage, err := s.storageOrError()
	if err != nil {
		return err
	}
	return storage.SaveHistoryRecord(record)
}

func (s *Service) SaveLogRecord(record resource.LogRecord) error {
	storage, err := s.storageOrError()
	if err != nil {
		return err
	}
	return storage.SaveLogRecord(record)
}

func (s *Service) GetLogRecords(startDate, endDate string) ([]resource.LogRecord, error) {
	storage, err := s.storageOrError()
	if err != nil {
		return nil, err
	}
	return storage.GetLogRecords(startDate, endDate)
}

func (s *Service) DeleteLogRecord(id int64) error {
	storage, err := s.storageOrError()
	if err != nil {
		return err
	}
	return storage.DeleteLogRecord(id)
}

func (s *Service) RecordAudit(resourceKey, operation string) error {
	return s.SaveAuditRecord(resource.AuditRecord{
		ResourceKey: resourceKey,
		Operation:   operation,
		Timestamp:   time.Now().Unix(),
		Count:       1,
	})
}

func (s *Service) RecordCommandHistory(operation, resourceKey string, _ bool) error {
	if operation == "completion" || operation == "help" || operation == "__complete" {
		return nil
	}
	record := resource.HistoryRecord{
		ID:          time.Now().UnixNano(),
		ResourceKey: resourceKey,
		Operation:   operation,
		Timestamp:   time.Now().Unix(),
		TimeStr:     time.Now().Format("2006-01-02 15:04:05"),
		Command:     operation,
	}
	return s.SaveHistoryRecord(record)
}

func (s *Service) CleanupResourceHistory(resourceKey string) (historyErr, auditErr error) {
	return s.DeleteHistoryRecords(resourceKey), s.DeleteAuditRecords(resourceKey)
}

// OpenStorage creates and initializes the selected client storage.
func OpenStorage(storageType, cloudAPIURL, cloudAPIKey string, cloudTimeout int, confFile string) (corestorage.Storage, error) {
	boltTimeout := 0
	if confFile == "" {
		if defaultConfPath, err := config.GetDefaultConfPath(); err == nil {
			if ttlConf, err := config.GetTtlConfFromFile(defaultConfPath); err == nil {
				boltTimeout = ttlConf.BoltDB.Timeout
			}
		}
	} else if ttlConf, err := config.GetTtlConfFromFile(confFile); err == nil {
		boltTimeout = ttlConf.BoltDB.Timeout
	}

	var storage corestorage.Storage
	switch storageType {
	case "sqlite":
		sqliteStorage := storagesqlite.NewSQLiteStorage()
		sqliteStorage.SetConfigFile(confFile)
		storage = sqliteStorage
	case "local", "bbolt":
		localStorage := storagebbolt.NewLocalStorage()
		localStorage.SetConfigFile(confFile)
		if boltTimeout > 0 {
			localStorage.SetTimeout(boltTimeout)
		}
		storage = localStorage
	case "cloud":
		if cloudAPIURL == "" || cloudAPIKey == "" {
			return nil, fmt.Errorf("cloud storage requires API URL and key")
		}
		storage = remote.NewStorage(cloudAPIURL, cloudAPIKey, cloudTimeout)
	case "sync":
		localStorage := storagebbolt.NewLocalStorage()
		localStorage.SetConfigFile(confFile)
		if boltTimeout > 0 {
			localStorage.SetTimeout(boltTimeout)
		}
		storage = clientsync.NewMirroredStorage(localStorage, remote.NewStorage(cloudAPIURL, cloudAPIKey, cloudTimeout))
	default:
		return nil, fmt.Errorf("unsupported storage type: %s (supported: sqlite, local/bbolt, cloud, sync)", storageType)
	}

	if err := storage.Init(); err != nil {
		_ = storage.Close()
		return nil, err
	}
	return storage, nil
}

func GetDBPath(confFile, storageType string) (string, error) {
	return storagebbolt.GetDBPath(confFile, storageType)
}
