package sync

import (
	"ttl-cli/internal/core/resource"
	corestorage "ttl-cli/internal/core/storage"
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
func (s *MirroredStorage) GetAllResources() (map[resource.ValJsonKey]resource.ValJson, error) {
	return s.local.GetAllResources()
}
func (s *MirroredStorage) SaveResource(k resource.ValJsonKey, v resource.ValJson) error {
	if err := s.local.SaveResource(k, v); err != nil {
		return err
	}
	return s.remote.SaveResource(k, v)
}
func (s *MirroredStorage) DeleteResource(k resource.ValJsonKey) error {
	if err := s.local.DeleteResource(k); err != nil {
		return err
	}
	return s.remote.DeleteResource(k)
}
func (s *MirroredStorage) UpdateResource(k resource.ValJsonKey, v resource.ValJson) error {
	if err := s.local.UpdateResource(k, v); err != nil {
		return err
	}
	return s.remote.UpdateResource(k, v)
}
func (s *MirroredStorage) GetTagStats() ([]resource.TagStat, error) { return s.local.GetTagStats() }
func (s *MirroredStorage) SaveAuditRecord(v resource.AuditRecord) error {
	if err := s.local.SaveAuditRecord(v); err != nil {
		return err
	}
	return s.remote.SaveAuditRecord(v)
}
func (s *MirroredStorage) GetAuditStats() (resource.AuditStats, error) {
	return s.local.GetAuditStats()
}
func (s *MirroredStorage) GetAllAuditRecords() ([]resource.AuditRecord, error) {
	return s.local.GetAllAuditRecords()
}
func (s *MirroredStorage) DeleteAuditRecords(k string) error {
	if err := s.local.DeleteAuditRecords(k); err != nil {
		return err
	}
	return s.remote.DeleteAuditRecords(k)
}
func (s *MirroredStorage) SaveHistoryRecord(v resource.HistoryRecord) error {
	if err := s.local.SaveHistoryRecord(v); err != nil {
		return err
	}
	return s.remote.SaveHistoryRecord(v)
}
func (s *MirroredStorage) GetAllHistoryRecords() ([]resource.HistoryRecord, error) {
	return s.local.GetAllHistoryRecords()
}
func (s *MirroredStorage) GetHistoryRecord(i int, o resource.SortOrder) (resource.HistoryRecord, error) {
	return s.local.GetHistoryRecord(i, o)
}
func (s *MirroredStorage) GetHistoryStats() (resource.HistoryStats, error) {
	return s.local.GetHistoryStats()
}
func (s *MirroredStorage) DeleteHistoryRecords(k string) error {
	if err := s.local.DeleteHistoryRecords(k); err != nil {
		return err
	}
	return s.remote.DeleteHistoryRecords(k)
}
func (s *MirroredStorage) SaveLogRecord(v resource.LogRecord) error {
	if err := s.local.SaveLogRecord(v); err != nil {
		return err
	}
	return s.remote.SaveLogRecord(v)
}
func (s *MirroredStorage) GetLogRecords(a, b string) ([]resource.LogRecord, error) {
	return s.local.GetLogRecords(a, b)
}
func (s *MirroredStorage) DeleteLogRecord(id int64) error {
	if err := s.local.DeleteLogRecord(id); err != nil {
		return err
	}
	return s.remote.DeleteLogRecord(id)
}
