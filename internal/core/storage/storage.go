package storage

import "ttl-cli/internal/core/resource"

// Storage is the shared persistence contract used by clients and servers.
type Storage interface {
	Init() error
	Close() error
	GetAllResources() (map[resource.Key]resource.Value, error)
	SaveResource(resource.Key, resource.Value) error
	DeleteResource(resource.Key) error
	UpdateResource(resource.Key, resource.Value) error
	GetTagStats() ([]resource.TagStat, error)
	SaveAuditRecord(resource.AuditRecord) error
	GetAuditStats() (resource.AuditStats, error)
	GetAllAuditRecords() ([]resource.AuditRecord, error)
	DeleteAuditRecords(string) error
	SaveHistoryRecord(resource.HistoryRecord) error
	GetAllHistoryRecords() ([]resource.HistoryRecord, error)
	GetHistoryRecord(int, resource.SortOrder) (resource.HistoryRecord, error)
	GetHistoryStats() (resource.HistoryStats, error)
	DeleteHistoryRecords(string) error
	SaveLogRecord(resource.LogRecord) error
	GetLogRecords(string, string) ([]resource.LogRecord, error)
	DeleteLogRecord(int64) error
}
