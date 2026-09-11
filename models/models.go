package models

import "ttl-cli/internal/core/resource"

const (
	ORIGIN = resource.Origin
	TAG    = resource.Tag
)

const Version = resource.Version

type ValJsonKey = resource.Key
type ValJson = resource.Value
type AuditRecord = resource.AuditRecord
type AuditStats = resource.AuditStats
type HistoryRecord = resource.HistoryRecord
type HistoryStats = resource.HistoryStats
type LogRecord = resource.LogRecord
type TagStat = resource.TagStat
type SortOrder = resource.SortOrder

const (
	Ascending  = resource.Ascending
	Descending = resource.Descending
)

type TtlIni struct {
	StorageType string                     `ini:"storage_type"`
	DbPath      string                     `ini:"db_path"`
	Workspace   string                     `ini:"workspace"`
	BoltDB      BoltDBConfig               `ini:"bbolt"`
	Workspaces  map[string]WorkspaceConfig `ini:"-"`
}

type BoltDBConfig struct {
	Timeout int `ini:"timeout"`
}

type WorkspaceConfig struct {
	DbPath      string `ini:"db_path"`
	StorageType string `ini:"storage_type"`
}

type WorkspacesSection struct {
	Workspaces map[string]WorkspaceConfig `ini:"-"`
}
