package sync

import (
	corestorage "ttl-cli/internal/core/storage"
	"ttl-cli/models"
)

// Storage is the contract required by the legacy read-local/write-both mode.
type Storage = corestorage.Storage

// MirroredStorage reads from local storage and mirrors mutations to remote storage.
type MirroredStorage struct {
	local  Storage
	remote Storage
}

func NewMirroredStorage(local, remote Storage) *MirroredStorage {
	return &MirroredStorage{local: local, remote: remote}
}

func (s *MirroredStorage) Init() error {
	if err := s.local.Init(); err != nil {
		return err
	}
	return s.remote.Init()
}
func (s *MirroredStorage) Close() error {
	if err := s.local.Close(); err != nil {
		return err
	}
	return s.remote.Close()
}
func (s *MirroredStorage) GetAllResources() (map[models.ValJsonKey]models.ValJson, error) {
	return s.local.GetAllResources()
}
func (s *MirroredStorage) SaveResource(k models.ValJsonKey, v models.ValJson) error {
	if err := s.local.SaveResource(k, v); err != nil {
		return err
	}
	return s.remote.SaveResource(k, v)
}
func (s *MirroredStorage) DeleteResource(k models.ValJsonKey) error {
	if err := s.local.DeleteResource(k); err != nil {
		return err
	}
	return s.remote.DeleteResource(k)
}
func (s *MirroredStorage) UpdateResource(k models.ValJsonKey, v models.ValJson) error {
	if err := s.local.UpdateResource(k, v); err != nil {
		return err
	}
	return s.remote.UpdateResource(k, v)
}
func (s *MirroredStorage) GetTagStats() ([]models.TagStat, error) { return s.local.GetTagStats() }
func (s *MirroredStorage) SaveAuditRecord(v models.AuditRecord) error {
	if err := s.local.SaveAuditRecord(v); err != nil {
		return err
	}
	return s.remote.SaveAuditRecord(v)
}
func (s *MirroredStorage) GetAuditStats() (models.AuditStats, error) { return s.local.GetAuditStats() }
func (s *MirroredStorage) GetAllAuditRecords() ([]models.AuditRecord, error) {
	return s.local.GetAllAuditRecords()
}
func (s *MirroredStorage) DeleteAuditRecords(k string) error {
	if err := s.local.DeleteAuditRecords(k); err != nil {
		return err
	}
	return s.remote.DeleteAuditRecords(k)
}
func (s *MirroredStorage) SaveHistoryRecord(v models.HistoryRecord) error {
	if err := s.local.SaveHistoryRecord(v); err != nil {
		return err
	}
	return s.remote.SaveHistoryRecord(v)
}
func (s *MirroredStorage) GetAllHistoryRecords() ([]models.HistoryRecord, error) {
	return s.local.GetAllHistoryRecords()
}
func (s *MirroredStorage) GetHistoryRecord(i int, o models.SortOrder) (models.HistoryRecord, error) {
	return s.local.GetHistoryRecord(i, o)
}
func (s *MirroredStorage) GetHistoryStats() (models.HistoryStats, error) {
	return s.local.GetHistoryStats()
}
func (s *MirroredStorage) DeleteHistoryRecords(k string) error {
	if err := s.local.DeleteHistoryRecords(k); err != nil {
		return err
	}
	return s.remote.DeleteHistoryRecords(k)
}
func (s *MirroredStorage) SaveLogRecord(v models.LogRecord) error {
	if err := s.local.SaveLogRecord(v); err != nil {
		return err
	}
	return s.remote.SaveLogRecord(v)
}
func (s *MirroredStorage) GetLogRecords(a, b string) ([]models.LogRecord, error) {
	return s.local.GetLogRecords(a, b)
}
func (s *MirroredStorage) DeleteLogRecord(id int64) error {
	if err := s.local.DeleteLogRecord(id); err != nil {
		return err
	}
	return s.remote.DeleteLogRecord(id)
}
